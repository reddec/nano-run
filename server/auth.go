package server

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (cfg Unit) enableAuthorization() func(gctx *gin.Context) {
	var handlers []AuthHandlerFunc
	if cfg.Authorization.JWT.Enable {
		handlers = append(handlers, cfg.Authorization.JWT.Create())
	}
	if cfg.Authorization.QueryToken.Enable {
		handlers = append(handlers, cfg.Authorization.QueryToken.Create())
	}
	if cfg.Authorization.HeaderToken.Enable {
		handlers = append(handlers, cfg.Authorization.HeaderToken.Create())
	}
	if cfg.Authorization.Basic.Enable {
		handlers = append(handlers, cfg.Authorization.Basic.Create())
	}
	return func(gctx *gin.Context) {
		var authorized = len(handlers) == 0
		for _, h := range handlers {
			if h(gctx.Request) {
				authorized = true
				break
			}
		}
		if !authorized {
			gctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		gctx.Next()
	}
}

type AuthHandlerFunc func(req *http.Request) bool

type JWT struct {
	Header string `yaml:"header"` // JWT header - by default Authorization
	Secret string `yaml:"secret"` // key to verify JWT
}

func (cfg JWT) Create() AuthHandlerFunc {
	header := cfg.Header
	if header == "" {
		header = "Authorization"
	}

	return func(req *http.Request) bool {
		rawToken := req.Header.Get(header)
		t, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unknown method")
			}
			return []byte(cfg.Secret), nil
		})
		return err == nil && t.Valid
	}
}

type QueryToken struct {
	Param  string   `yaml:"param"`  // query name - by default 'token'
	Tokens []string `yaml:"tokens"` // allowed tokens
}

func (cfg QueryToken) Create() AuthHandlerFunc {
	param := cfg.Param
	if param == "" {
		param = "token"
	}
	tokens := map[string]bool{}
	for _, k := range cfg.Tokens {
		tokens[k] = true
	}
	return func(req *http.Request) bool {
		token := req.URL.Query().Get(param)
		return tokens[token]
	}
}

type HeaderToken struct {
	Header string   `yaml:"header"` // header name - by default X-Api-Token
	Tokens []string `yaml:"tokens"` // allowed tokens
}

func (cfg HeaderToken) Create() AuthHandlerFunc {
	header := cfg.Header
	if header == "" {
		header = "X-Api-Token"
	}
	tokens := map[string]bool{}
	for _, k := range cfg.Tokens {
		tokens[k] = true
	}
	return func(req *http.Request) bool {
		token := req.URL.Query().Get(header)
		return tokens[token]
	}
}

type Basic struct {
	Users map[string]string `yaml:"users"` // users -> bcrypted password map
}

func (cfg Basic) Create() AuthHandlerFunc {
	return func(req *http.Request) bool {
		u, p, ok := req.BasicAuth()
		if !ok {
			return false
		}
		h, ok := cfg.Users[u]
		if !ok {
			return false
		}
		return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
	}
}
