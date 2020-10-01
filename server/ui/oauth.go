package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	Config     oauth2.Config
	ProfileURL string
	RedirectTo string
	LoginField string
}

func (cfg OAuth2) Attach(router gin.IRouter, storage SessionStorage) {
	router.GET("/login", func(gctx *gin.Context) {
		state := uuid.New().String()
		gctx.SetCookie("oauth", state, 3600, "", "", false, true)
		u := cfg.Config.AuthCodeURL(state)
		gctx.Redirect(http.StatusTemporaryRedirect, u)
	})
	router.GET("/callback", func(gctx *gin.Context) {
		savedState, _ := gctx.Cookie("oauth")
		if savedState != gctx.Query("state") {
			gctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		gctx.SetCookie("oauth", "", -1, "", "", false, true)

		token, err := cfg.Config.Exchange(gctx.Request.Context(), gctx.Query("code"))
		if err != nil {
			gctx.AbortWithError(http.StatusForbidden, err)
			return
		}

		if !token.Valid() {
			gctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		sessionID := uuid.New().String()
		session := newOAuthSession(token)
		err = session.fetchLogin(cfg.ProfileURL, cfg.LoginField)
		if err != nil {
			gctx.AbortWithError(http.StatusForbidden, err)
			return
		}

		gctx.SetCookie(sessionCookie, sessionID, 0, "", "", false, true)

		storage.Save(sessionID, session)
		log.Println("user", session.Login(), "authorized via oauth2")

		gctx.Redirect(http.StatusTemporaryRedirect, cfg.RedirectTo)
	})
}

func newOAuthSession(token *oauth2.Token) *oauthSession {
	return &oauthSession{
		token: oauth2.StaticTokenSource(token),
	}
}

type oauthSession struct {
	login string
	token oauth2.TokenSource
}

func (ss *oauthSession) Login() string {
	return ss.login
}

func (ss *oauthSession) Valid() bool {
	token, err := ss.token.Token()
	if err != nil {
		log.Println(err)
		return false
	}
	return token.Valid()
}

func (ss *oauthSession) GetJSON(url string, response interface{}) error {
	t, err := ss.token.Token()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	t.SetAuthHeader(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return json.NewDecoder(res.Body).Decode(response)
}

func (ss *oauthSession) fetchLogin(url string, field string) error {
	var profile = make(map[string]interface{})
	err := ss.GetJSON(url, &profile)
	if err != nil {
		return err
	}
	log.Printf("profile: %+v", profile)
	l, ok := profile[field]
	if !ok {
		return fmt.Errorf("not field %s in response", field)
	}
	s, ok := l.(string)
	if !ok {
		return fmt.Errorf("field %s is not string", field)
	}
	ss.login = s
	return nil
}
