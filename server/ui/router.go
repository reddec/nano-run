package ui

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/gin-gonic/gin"

	"nano-run/server"
)

func Expose(units []server.Unit, uiDir string) http.Handler {
	router := gin.New()
	Attach(router, units, uiDir)
	return router
}

func Attach(router gin.IRouter, units []server.Unit, uiDir string) {
	ui := &uiRouter{
		dir:   uiDir,
		units: units,
	}
	router.GET("/units", ui.listUnits)
	router.GET("/unit/:name", ui.unitInfo)
}

type uiRouter struct {
	dir   string
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

func (ui *uiRouter) getTemplate(name string) *template.Template {
	t, err := template.New("").Funcs(sprig.HtmlFuncMap()).ParseFiles(filepath.Join(ui.dir, name))
	if err == nil {
		return t
	}
	t, err = template.New("").Parse("<html><body>Ooops... Page not found</body></html>")
	if err != nil {
		panic(err)
	}
	return t
}
