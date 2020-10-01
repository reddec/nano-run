package ui

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"

	"nano-run/server"
	"nano-run/services/meta"
	"nano-run/worker"
)

func Expose(units []server.Unit, workers []*worker.Worker, uiDir string, auth Authorization) http.Handler {
	router := gin.New()
	router.SetFuncMap(sprig.HtmlFuncMap())
	router.LoadHTMLGlob(filepath.Join(uiDir, "*.html"))
	Attach(router, units, workers, auth)
	return router
}

func Attach(router gin.IRouter, units []server.Unit, workers []*worker.Worker, auth Authorization) {
	ui := &uiRouter{
		units: make(map[string]unitInfo),
	}

	for i := range units {
		u := units[i]
		w := workers[i]
		ui.units[u.Name()] = unitInfo{
			Unit:   u,
			Worker: w,
		}
	}
	sessions := &memorySessions{}
	if r, ok := router.(*gin.RouterGroup); ok {
		router.Use(func(gctx *gin.Context) {
			gctx.Set("base", r.BasePath())
			gctx.Next()
		})
	}
	router.GET("", func(gctx *gin.Context) {
		gctx.Redirect(http.StatusTemporaryRedirect, "unit/")
	})
	auth.attach(router.Group("/auth"), "login.html", sessions)

	restricted := router.Group("/unit/").Use(auth.restrict(func(gctx *gin.Context) string {
		b := base(gctx)
		return b.Rel("/auth/")
	}, sessions))

	restricted.GET("/", ui.listUnits)
	restricted.GET("/:name/", ui.unitInfo)
	restricted.POST("/:name/", ui.unitInvoke)
	restricted.GET("/:name/history", ui.unitHistory)
	restricted.GET("/:name/request/:request/", ui.unitRequestInfo)
	restricted.POST("/:name/request/:request/", ui.unitRequestRetry)
	restricted.GET("/:name/request/:request/payload", ui.unitRequestPayload)
	restricted.GET("/:name/request/:request/attempt/:attempt/", ui.unitRequestAttemptInfo)
	restricted.GET("/:name/request/:request/attempt/:attempt/result", ui.unitRequestAttemptResult)
}

type uiRouter struct {
	units map[string]unitInfo
}

func (ui *uiRouter) unitRequestAttemptResult(gctx *gin.Context) {
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	attemptID := gctx.Param("attempt")

	f, err := info.Worker.Blobs().Get(attemptID)
	if err != nil {
		gctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer f.Close()
	gctx.Status(http.StatusOK)
	_, _ = io.Copy(gctx.Writer, f)
}

type attemptInfo struct {
	requestInfo
	AttemptID string
	Attempt   meta.Attempt
}

func (ui *uiRouter) unitRequestAttemptInfo(gctx *gin.Context) {
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	requestId := gctx.Param("request")
	request, err := info.Worker.Meta().Get(requestId)
	if err != nil {
		gctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	attemptID := gctx.Param("attempt")

	var found bool
	var attempt meta.Attempt
	for _, atp := range request.Attempts {
		if atp.ID == attemptID {
			attempt = atp
			found = true
			break
		}
	}
	if !found {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	gctx.HTML(http.StatusOK, "unit-request-attempt-info.html", attemptInfo{
		requestInfo: requestInfo{
			baseResponse: base(gctx),
			unitInfo:     info,
			Request:      request,
			RequestID:    requestId,
		},
		AttemptID: attemptID,
		Attempt:   attempt,
	})
}

type requestInfo struct {
	unitInfo
	baseResponse
	Request   *meta.Request
	RequestID string
}

func (ui *uiRouter) unitRequestInfo(gctx *gin.Context) {
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	requestId := gctx.Param("request")
	request, err := info.Worker.Meta().Get(requestId)
	if err != nil {
		gctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	gctx.HTML(http.StatusOK, "unit-request-info.html", requestInfo{
		baseResponse: base(gctx),
		unitInfo:     info,
		Request:      request,
		RequestID:    requestId,
	})
}

func (ui *uiRouter) unitRequestRetry(gctx *gin.Context) {
	name := gctx.Param("name")
	item, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	requestID := gctx.Param("request")

	id, err := item.Worker.Retry(gctx.Request.Context(), requestID)
	if err != nil {
		log.Println("failed to retry:", err)
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	gctx.Redirect(http.StatusSeeOther, base(gctx).Rel("/unit", name, "request", id))
}

func (ui *uiRouter) unitRequestPayload(gctx *gin.Context) {
	name := gctx.Param("name")
	item, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	requestID := gctx.Param("request")
	info, err := item.Worker.Meta().Get(requestID)
	if !ok {
		gctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	gctx.Header("Last-Modified", info.CreatedAt.Format(time.RFC850))
	f, err := item.Worker.Blobs().Get(requestID)
	if err != nil {
		log.Println("failed to get data:", err)
		gctx.AbortWithError(http.StatusInternalServerError, err)
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

func (ui *uiRouter) unitInfo(gctx *gin.Context) {
	type viewUnit struct {
		unitInfo
		baseResponse
	}
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	gctx.HTML(http.StatusOK, "unit-info.html", viewUnit{
		unitInfo:     info,
		baseResponse: base(gctx),
	})
}

func (ui *uiRouter) unitInvoke(gctx *gin.Context) {
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	data := gctx.PostForm("body")
	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(data))
	if err != nil {
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	id, err := info.Worker.Enqueue(req)
	if err != nil {
		gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	gctx.Redirect(http.StatusSeeOther, base(gctx).Rel("/unit", name, "request", id))
}

func (ui *uiRouter) listUnits(gctx *gin.Context) {
	var reply struct {
		baseResponse
		Units []server.Unit
	}

	var units = make([]server.Unit, 0, len(ui.units))
	for _, info := range ui.units {
		units = append(units, info.Unit)
	}
	reply.baseResponse = base(gctx)
	reply.Units = units
	gctx.HTML(http.StatusOK, "units-list.html", reply)
}

func (ui *uiRouter) unitHistory(gctx *gin.Context) {
	type viewUnit struct {
		unitInfo
		baseResponse
	}
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	gctx.HTML(http.StatusOK, "unit-history.html", viewUnit{
		unitInfo:     info,
		baseResponse: base(gctx),
	})
}
