package pkg

import (
	"auth/config"
	"fmt"
	"net/http"
	"time"
)

type cookieService struct {
	domain   string
	path     string
	sameSite http.SameSite
	secure   bool
	httpOnly bool
}

type CookieService interface {
	Set(w http.ResponseWriter, name, value string, expires time.Time)
	Clear(w http.ResponseWriter, name string)
}

func NewCookieService(cfg config.Cookie) CookieService {
	var sameSite http.SameSite
	switch cfg.SameSite {
	case "Strict":
		sameSite = http.SameSiteStrictMode
	case "Lax":
		sameSite = http.SameSiteLaxMode
	case "None":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	return &cookieService{
		domain:   cfg.Domain,
		path:     cfg.Path,
		sameSite: sameSite,
		secure:   cfg.Secure,
		httpOnly: cfg.HttpOnly,
	}
}

func (s *cookieService) Set(w http.ResponseWriter, name, value string, expires time.Time) {
	fmt.Println("expires: ", expires)
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     s.path,
		Domain:   s.domain,
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
		SameSite: s.sameSite,
	}
	http.SetCookie(w, cookie)
}

func (s *cookieService) Clear(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   s.path,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
