package main

import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if _, err := request.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, request)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler to show consent screen on client side
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Panic("Error when trying to get provider", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Panic("Error when trying to get login url", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		log.Println("Login")
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Panic("Error when trying to get provider")
		}
		credential, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Panic("Error when trying to complete auth", err)
		}
		user, err := provider.GetUser(credential)
		if err != nil {
			log.Panic("Error when trying to get user", err)
		}
		m := md5.New()
		_, _ = io.WriteString(m, user.Email())
		chatUser := &chatUser{User: user}
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarUrl, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Panic("Error when trying to get avatar url", err)
		}

		authCookievalue := objx.New(map[string]interface{}{
			"userId":     fmt.Sprintf("%x", m.Sum(nil)),
			"name":       user.Name(),
			"avatar_url": avatarUrl,
			"email":      user.Email(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookievalue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func InitAuthProvider() {
	gomniauth.SetSecurityKey("sjkdkasljdklsajklasjdklsajdklsajdklsajdkas")
	gomniauth.WithProviders(
		facebook.New("", "", "http://localhost:9090/auth/callback/facebook"),
		github.New("", "", "http://localhost:9090/auth/callback/github"),
		google.New(os.Getenv("GOOGLE_CLIENT_KEY"), os.Getenv("GOOGLE_SECRET_KEY"), "http://localhost:9090/auth/callback/google"),
	)
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
