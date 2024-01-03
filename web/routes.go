package web

import (
	"fmt"
	"html/template"
	"net/http"

	"goshrest/internal/common"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

const oAuthStateKey = "oauth_state"

type WebRoutes struct {
	oAuth2Conf *oauth2.Config
}

func NewWebRoutes(oAuth2Conf *oauth2.Config) *WebRoutes {
	return &WebRoutes{
		oAuth2Conf: oAuth2Conf,
	}
}

func (r *WebRoutes) Routes() []common.Router {
	return []common.Router{
		r.index(),
		r.gAuthInitRoute(),
		r.gAuthCallbackRoute(),
	}
}

func (r *WebRoutes) index() common.Router {
	return common.Router{
		Method:  "GET",
		Pattern: "/",
		Handler: func(res http.ResponseWriter, _ *http.Request) {
			t, _ := template.ParseFiles("web/templates/index.html")
			t.Execute(res, false)
		},
	}
}

func (r *WebRoutes) gAuthInitRoute() common.Router {
	return common.Router{
		Method:  "GET",
		Pattern: "/auth/google",
		Handler: func(res http.ResponseWriter, req *http.Request) {
			randState := uuid.New().String()

			redirectURL := r.oAuth2Conf.AuthCodeURL(randState, oauth2.AccessTypeOffline)

			http.SetCookie(res, &http.Cookie{
				Name:     oAuthStateKey,
				Value:    randState,
				HttpOnly: true,
				Path:     "/",
				MaxAge:   86400 * 30, // 30 days
			})
			http.Redirect(res, req, redirectURL, http.StatusFound)
		},
	}
}

func (r *WebRoutes) gAuthCallbackRoute() common.Router {
	return common.Router{
		Method:  "GET",
		Pattern: "/auth/google/callback",
		Handler: func(res http.ResponseWriter, req *http.Request) {
			queryParams := req.URL.Query()
			state := queryParams.Get("state")

			stateCookie, err := req.Cookie(oAuthStateKey)
			if err != nil {
				res.Write([]byte("Missing 'oauth_state' cookie"))
				return
			}

			if state != stateCookie.Value {
				res.Write([]byte("state mismatch"))
				return
			}

			code := queryParams.Get("code")
			token, err := r.oAuth2Conf.Exchange(req.Context(), code)
			if err != nil {
				res.Write([]byte("Code-Token Exchange Failed"))
				return
			}

			tokenSource := option.WithTokenSource(r.oAuth2Conf.TokenSource(req.Context(), token))

			userService, err := goauth2.NewService(req.Context(), tokenSource)
			if err != nil {
				panic(err)
			}

			ui, err := userService.Userinfo.V2.Me.Get().Do()
			if err != nil {
				panic(err)
			}

			fmt.Println(ui)

			t, _ := template.ParseFiles("web/templates/success.html")
			t.Execute(res, ui)
		},
	}
}
