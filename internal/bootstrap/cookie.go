package bootstrap

import (
	"net/http"
	"os"

	"github.com/wolftotem4/golava-core/cookie"
)

func initCookie() *cookie.CookieManager {
	path := os.Getenv("SESSION_PATH")
	if path == "" {
		path = "/"
	}

	var secure bool
	secureStr := os.Getenv("SESSION_SECURE_COOKIE")
	if secureStr == "true" || secureStr == "1" {
		secure = true
	}

	var sameSite http.SameSite
	switch os.Getenv("SESSION_SAME_SITE") {
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}

	return &cookie.CookieManager{
		Path:     path,
		Domain:   os.Getenv("SESSION_DOMAIN"),
		Secure:   secure,
		SameSite: sameSite,
	}
}
