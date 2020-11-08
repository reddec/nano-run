package ui

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"

	"nano-run/server"
	"nano-run/services/meta"
	"nano-run/worker"
)

func Expose(units []server.Unit, workers []*worker.Worker, cronEntries []*server.CronEntry, uiDir string, auth Authorization) http.Handler {
	router := gin.New()
	router.SetFuncMap(sprig.HtmlFuncMap())
	router.LoadHTMLGlob(filepath.Join(uiDir, "*.html"))
	Attach(router, units, workers, cronEntries, auth)
	return router
}

func Attach(router gin.IRouter, units []server.Unit, workers []*worker.Worker, cronEntries []*server.CronEntry, auth Authorization) {
	ui := &uiRouter{
		units: make(map[string]unitInfo),
	}

	var offset int
	for i := range units {
		u := units[i]
		w := workers[i]
		var ce []*server.CronEntry

		var last int
		for last = offset; last < len(cronEntries); last++ {
			if cronEntries[last].Worker != w {
				break
			}
		}
		ce = cronEntries[offset:last]
		offset = last

		ui.units[u.Name()] = unitInfo{
			Unit:        u,
			Worker:      w,
			CronEntries: ce,
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

	guard := auth.restrict(func(gctx *gin.Context) string {
		b := base(gctx)
		return b.Rel("/auth/")
	}, sessions)

	unitsRoutes := router.Group("/unit/").Use(guard)

	unitsRoutes.GET("/", ui.listUnits)
	unitsRoutes.GET("/:name/", ui.unitInfo)
	unitsRoutes.POST("/:name/", ui.unitInvoke)
	unitsRoutes.GET("/:name/history", ui.unitHistory)
	unitsRoutes.GET("/:name/request/:request/", ui.unitRequestInfo)
	unitsRoutes.POST("/:name/request/:request/", ui.unitRequestRetry)
	unitsRoutes.GET("/:name/request/:request/payload", ui.unitRequestPayload)
	unitsRoutes.GET("/:name/request/:request/attempt/:attempt/", ui.unitRequestAttemptInfo)
	unitsRoutes.GET("/:name/request/:request/attempt/:attempt/result", ui.unitRequestAttemptResult)
	unitsRoutes.GET("/:name/cron/:index", ui.unitCronInfo)

	cronRoutes := router.Group("/cron/").Use(guard)
	cronRoutes.GET("/", ui.listCron)
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
		_ = gctx.AbortWithError(http.StatusNotFound, err)
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
	requestID := gctx.Param("request")
	request, err := info.Worker.Meta().Get(requestID)
	if err != nil {
		_ = gctx.AbortWithError(http.StatusNotFound, err)
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
			RequestID:    requestID,
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
	requestID := gctx.Param("request")
	request, err := info.Worker.Meta().Get(requestID)
	if err != nil {
		_ = gctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	gctx.HTML(http.StatusOK, "unit-request-info.html", requestInfo{
		baseResponse: base(gctx),
		unitInfo:     info,
		Request:      request,
		RequestID:    requestID,
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
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
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
		_ = gctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	gctx.Header("Last-Modified", info.CreatedAt.Format(time.RFC850))
	f, err := item.Worker.Blobs().Get(requestID)
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
	req, err := http.NewRequestWithContext(gctx.Request.Context(), http.MethodPost, "/", bytes.NewBufferString(data))
	if err != nil {
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	id, err := info.Worker.Enqueue(req)
	if err != nil {
		_ = gctx.AbortWithError(http.StatusInternalServerError, err)
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
	sort.Slice(units, func(i, j int) bool {
		return units[i].Name() < units[j].Name()
	})
	reply.baseResponse = base(gctx)
	reply.Units = units
	gctx.HTML(http.StatusOK, "units-list.html", reply)
}

func (ui *uiRouter) listCron(gctx *gin.Context) {
	type uiEntry struct {
		Index int
		Entry *server.CronEntry
	}

	var reply struct {
		baseResponse
		Entries []uiEntry
	}

	var specs = make([]uiEntry, 0, len(ui.units))
	for _, info := range ui.units {
		for i, spec := range info.CronEntries {
			specs = append(specs, uiEntry{
				Index: i,
				Entry: spec,
			})
		}
	}
	sort.Slice(specs, func(i, j int) bool {
		if specs[i].Entry.Config.Name() < specs[j].Entry.Config.Name() {
			return true
		} else if specs[i].Entry.Config.Name() == specs[j].Entry.Config.Name() && specs[i].Index < specs[j].Index {
			return true
		}
		return false
	})
	reply.baseResponse = base(gctx)
	reply.Entries = specs
	gctx.HTML(http.StatusOK, "cron-list.html", reply)
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

func (ui *uiRouter) unitCronInfo(gctx *gin.Context) {
	type viewUnit struct {
		unitInfo
		baseResponse
		Cron  *server.CronEntry
		Label string
	}
	name := gctx.Param("name")
	info, ok := ui.units[name]
	if !ok {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	strIndex := gctx.Param("index")
	index, err := strconv.Atoi(strIndex)

	if err != nil || index < 0 || index >= len(info.CronEntries) {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	entry := info.CronEntries[index]
	gctx.HTML(http.StatusOK, "unit-cron-info.html", viewUnit{
		unitInfo:     info,
		baseResponse: base(gctx),
		Cron:         entry,
		Label:        entry.Spec.Label(strconv.Itoa(index)),
	})
}
