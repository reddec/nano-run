package api

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"nano-run/services/meta"
	"nano-run/worker"
)

// Expose worker as HTTP handler:
//     POST /                  - post task async, returns 303 See Other and location
//     PUT  /                  - post task synchronously, supports ?wait=<duration> parameter for custom wait time
//     GET  /:id               - get task info.
//     POST /:id               - retry task, redirects to /:id
//     GET  /:id/completed     - redirect to completed attempt (or 404 if task not yet complete)
//     GET  /:id/attempt/:atid - get attempt result (as-is).
//     GET  /:id/request       - replay request (as-is).
func Expose(router gin.IRouter, wrk *worker.Worker, defaultWaitTime time.Duration) {
	handler := &workerHandler{wrk: wrk, defaultWait: defaultWaitTime}
	router.POST("/", func(gctx *gin.Context) {
		id, err := wrk.Enqueue(gctx.Request)
		if err != nil {
			log.Println("failed to enqueue:", err)
			_ = gctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		gctx.Header("X-Correlation-Id", id)
		gctx.Redirect(http.StatusSeeOther, id+"?"+gctx.Request.URL.RawQuery)
	})
	router.PUT("/", handler.createSyncTask)
	taskRoutes := router.Group("/:id")
	taskRoutes.GET("", handler.getTask)
	taskRoutes.POST("", handler.retry)
	taskRoutes.DELETE("", handler.completeRequest)
	taskRoutes.GET("/completed", handler.getComplete)
	// get attempt result as-is.
	taskRoutes.GET("/attempt/:attemptId", handler.getAttempt)
	// get recorded request.
	taskRoutes.GET("/request", handler.getRequest)
}

type workerHandler struct {
	wrk         *worker.Worker
	defaultWait time.Duration
}

func (wh *workerHandler) createSyncTask(gctx *gin.Context) {
	var queryParams struct {
		Wait time.Duration `query:"wait" form:"wait"`
	}
	queryParams.Wait = wh.defaultWait

	if err := gctx.BindQuery(&queryParams); err != nil {
		return
	}
	if queryParams.Wait <= 0 {
		gctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tracker, err := wh.wrk.EnqueueWithTracker(gctx.Request)
	if err != nil {
		log.Println("failed to enqueue:", err)
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.Header("X-Correlation-Id", tracker.ID())
	select {
	case <-tracker.Done():
		gctx.Redirect(http.StatusSeeOther, tracker.ID()+"/completed?"+gctx.Request.URL.RawQuery)
	case <-time.After(queryParams.Wait):
		gctx.AbortWithStatus(http.StatusGatewayTimeout)
	}
}

// get request meta information.
func (wh *workerHandler) getTask(gctx *gin.Context) {
	requestID := gctx.Param("id")
	info, err := wh.wrk.Meta().Get(requestID)
	if err != nil {
		log.Println("failed access request", requestID, ":", err)
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	gctx.Header("X-Correlation-Id", requestID)
	gctx.Header("Content-Version", strconv.Itoa(len(info.Attempts)))
	// modification time
	setLastModify(gctx, info)

	gctx.Header("Age", strconv.FormatInt(int64(time.Since(info.CreatedAt)/time.Second), 10))
	if info.Complete {
		gctx.Header("X-Status", "complete")
	} else {
		gctx.Header("X-Status", "processing")
	}

	if len(info.Attempts) > 0 {
		gctx.Header("X-Last-Attempt", info.Attempts[len(info.Attempts)-1].ID)
		gctx.Header("X-Last-Attempt-At", info.Attempts[len(info.Attempts)-1].CreatedAt.Format(time.RFC850))
	}

	if info.Complete {
		lastAttempt := info.Attempts[len(info.Attempts)-1]
		gctx.Request.URL.Path += "/attempt/" + lastAttempt.ID
		gctx.Header("Location", gctx.Request.URL.String())
	}
	gctx.IndentedJSON(http.StatusOK, info)
}

func (wh *workerHandler) retry(gctx *gin.Context) {
	requestID := gctx.Param("id")
	id, err := wh.wrk.Retry(gctx.Request.Context(), requestID)
	if err != nil {
		log.Println("failed to retry:", err)
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.Header("X-Correlation-Id", id)
	gctx.Redirect(http.StatusSeeOther, id+"?"+gctx.Request.URL.RawQuery)
}

func (wh *workerHandler) completeRequest(gctx *gin.Context) {
	requestID := gctx.Param("id")
	info, err := wh.wrk.Meta().Get(requestID)
	if err != nil {
		log.Println("failed access request", requestID, ":", err)
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !info.Complete {
		err = wh.wrk.Meta().Complete(requestID)
		if err != nil {
			log.Println("failed to mark request as complete:", err)
			_ = gctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	gctx.AbortWithStatus(http.StatusNoContent)
}

func (wh *workerHandler) getComplete(gctx *gin.Context) {
	requestID := gctx.Param("id")
	info, err := wh.wrk.Meta().Get(requestID)
	if err != nil {
		log.Println("failed access request", requestID, ":", err)
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !info.Complete {
		gctx.AbortWithStatus(http.StatusTooEarly)
		return
	}
	lastAttempt := info.Attempts[len(info.Attempts)-1]
	gctx.Redirect(http.StatusMovedPermanently, "attempt/"+lastAttempt.ID)
}

func (wh *workerHandler) getAttempt(gctx *gin.Context) {
	requestID := gctx.Param("id")
	attemptID := gctx.Param("attemptId")
	info, err := wh.wrk.Meta().Get(requestID)
	if err != nil {
		log.Println("failed access request", requestID, ":", err)
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	var attempt meta.Attempt
	var found bool
	for _, atp := range info.Attempts {
		if atp.ID == attemptID {
			found = true
			attempt = atp
			break
		}
	}
	if !found {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	body, err := wh.wrk.Blobs().Get(attempt.ID)
	if err != nil {
		log.Println("failed to get body:", err)
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer body.Close()
	gctx.Header("Last-Modified", attempt.CreatedAt.Format(time.RFC850))
	if info.Complete {
		gctx.Header("X-Status", "complete")
	} else {
		gctx.Header("X-Status", "processing")
	}
	gctx.Header("X-Processed", "true")
	for k, v := range attempt.Headers {
		gctx.Request.Header[k] = v
	}
	_, _ = io.Copy(gctx.Writer, body)
}

func (wh *workerHandler) getRequest(gctx *gin.Context) {
	requestID := gctx.Param("id")
	info, err := wh.wrk.Meta().Get(requestID)
	if err != nil {
		log.Println("failed access request", requestID, ":", err)
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	gctx.Header("Last-Modified", info.CreatedAt.Format(time.RFC850))
	f, err := wh.wrk.Blobs().Get(requestID)
	if err != nil {
		log.Println("failed to get data:", err)
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	gctx.Header("X-Method", info.Method)
	gctx.Header("X-Request-Uri", info.URI)

	for k, v := range info.Headers {
		gctx.Request.Header[k] = v
	}
	gctx.Status(http.StatusOK)
	_, _ = io.Copy(gctx.Writer, f)
}

func setLastModify(gctx *gin.Context, info *meta.Request) {
	if info.Complete {
		gctx.Header("Last-Modified", info.CompleteAt.Format(time.RFC850))
	} else if len(info.Attempts) > 0 {
		gctx.Header("Last-Modified", info.Attempts[len(info.Attempts)-1].CreatedAt.Format(time.RFC850))
	} else {
		gctx.Header("Last-Modified", info.CreatedAt.Format(time.RFC850))
	}
}
