package ui

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"nano-run/server"
)

func Expose(units []server.Unit, uiDir string) http.Handler {
	router := gin.New()
	router.LoadHTMLGlob(filepath.Join(uiDir, "*.html"))
	Attach(router, units)
	return router
}

func Attach(router gin.IRouter, units []server.Unit) {
	ui := &uiRouter{
		units: units,
	}
	router.GET("", func(gctx *gin.Context) {
		gctx.Redirect(http.StatusTemporaryRedirect, "units")
	})
	router.GET("/units", ui.listUnits)
	router.GET("/unit/:name", ui.unitInfo)
}

type uiRouter struct {
	units []server.Unit
}

func (ui *uiRouter) unitInfo(gctx *gin.Context) {
	name := gctx.Param("name")
	var unit *server.Unit
	for _, u := range ui.units {
		if u.Name() == name {
			unit = &u
			break
		}
	}
	if unit == nil {
		gctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	var reply struct {
		Unit *server.Unit
	}
	reply.Unit = unit
	gctx.HTML(http.StatusOK, "unit-info.html", reply)
}

func (ui *uiRouter) listUnits(gctx *gin.Context) {
	var reply struct {
		Units []server.Unit
	}
	reply.Units = ui.units
	gctx.HTML(http.StatusOK, "units-list.html", reply)
}
