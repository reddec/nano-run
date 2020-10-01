package ui

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Context    *gin.Context
	Login      string
	Authorized bool
}

func (br baseResponse) Rel(path string, tail ...string) string {
	var chunks = append([]string{path}, tail...)
	path = strings.Join(chunks, "/")
	if len(path) == 0 || path[0] != '/' {
		return path
	}
	toRoot := strings.Repeat("../", strings.Count(br.Context.Request.RequestURI, "/"))
	if len(toRoot) > 0 {
		toRoot = toRoot[:len(toRoot)-1]
	}
	return toRoot + br.Context.GetString("base") + path[1:]
}

func base(gctx *gin.Context) baseResponse {
	return baseResponse{
		Authorized: gctx.GetBool(ctxAuthorized),
		Context:    gctx,
		Login:      gctx.GetString(ctxLogin),
	}
}
