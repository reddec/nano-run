package ui

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	redirectToCookie = "redirect-to"
	ctxLogin         = "login"
	ctxAuthorized    = "authorized"
)

type Strategy struct {
	Icon        string   `yaml:"icon"`
	Title       string   `yaml:"title"`
	Key         string   `yaml:"key"`
	Secret      string   `yaml:"secret"`
	AuthURL     string   `yaml:"auth_url"`
	TokenURL    string   `yaml:"token_url"`
	ProfileURL  string   `yaml:"profile_url"`
	CallbackURL string   `yaml:"callback_url"`
	LoginField  string   `yaml:"login_field"`
	Scopes      []string `yaml:"scopes"`
}

func (st *Strategy) Enable(router gin.IRouter, sessionStorage SessionStorage) {
	if st == nil {
		return
	}
	ep := &OAuth2{
		Config: oauth2.Config{
			ClientID:     st.Key,
			ClientSecret: st.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  st.AuthURL,
				TokenURL: st.TokenURL,
			},
			RedirectURL: st.CallbackURL,
			Scopes:      st.Scopes,
		},
		ProfileURL: st.ProfileURL,
		RedirectTo: "../success",
		LoginField: st.LoginField,
	}

	ep.Attach(router, sessionStorage)
}

type Authorization struct {
	OAuth2       *Strategy `yaml:"oauth2"`
	AllowedUsers []string  `yaml:"users"`
}

func (auth Authorization) Enabled() bool {
	return auth.OAuth2 != nil
}

func (auth Authorization) attach(router gin.IRouter, loginTemplate string, sessionStorage SessionStorage) {
	auth.OAuth2.Enable(router.Group("/oauth2/"), sessionStorage)

	router.GET("/logout", func(gctx *gin.Context) {
		id, _ := gctx.Cookie(sessionCookie)
		sessionStorage.Delete(id)
		gctx.Redirect(http.StatusTemporaryRedirect, "../")
	})
	router.GET("/success", func(gctx *gin.Context) {
		redirectTo, err := gctx.Cookie(redirectToCookie)
		if err == nil && redirectTo != "" {
			gctx.SetCookie(redirectToCookie, "", -1, "", "", false, true)
			gctx.Redirect(http.StatusTemporaryRedirect, redirectTo)
		} else {
			gctx.Redirect(http.StatusTemporaryRedirect, "../../")
		}
	})
	router.GET("/", func(gctx *gin.Context) {
		var reply struct {
			Auth Authorization
		}
		reply.Auth = auth
		gctx.HTML(http.StatusOK, loginTemplate, reply)
	})
}

func (auth Authorization) restrict(redirectTo func(gctx *gin.Context) string, sessionStorage SessionStorage) gin.HandlerFunc {
	if !auth.Enabled() {
		return func(gctx *gin.Context) {
			gctx.Set(ctxAuthorized, false)
			gctx.Set(ctxLogin, "anonymous")
			gctx.Next()
		}
	}
	return func(gctx *gin.Context) {
		sessionID, err := gctx.Cookie(sessionCookie)
		session, ok := sessionStorage.Get(sessionID)
		if err != nil || !ok || session == nil || !session.Valid() {
			gctx.SetCookie(redirectToCookie, gctx.Request.RequestURI, 3600, "", "", false, true)
			gctx.Redirect(http.StatusTemporaryRedirect, redirectTo(gctx))
			gctx.Abort()
			return
		}
		login := session.Login()

		found := len(auth.AllowedUsers) == 0
		for _, u := range auth.AllowedUsers {
			if u == login {
				found = true
				break
			}
		}
		if !found {
			log.Println("user", login, "not allowed")
			gctx.Redirect(http.StatusTemporaryRedirect, redirectTo(gctx))
			gctx.Abort()
			return
		}
		gctx.Set(ctxAuthorized, true)
		gctx.Set(ctxLogin, login)
		gctx.Next()
	}
}
