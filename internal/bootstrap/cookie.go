package bootstrap

import (
	"net/http"
	"os"

	"github.com/wolftotem4/golava-core/cookie"
	"github.com/wolftotem4/golava-core/encryption"
	"github.com/wolftotem4/golava/internal/env"
)

func InitCookie(encrypter encryption.IEncrypter) *cookie.CookieFactory {
	path := os.Getenv("SESSION_PATH")
	if path == "" {
		path = "/"
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

	domain := os.Getenv("SESSION_DOMAIN")
	secure := env.Bool(os.Getenv("SESSION_SECURE_COOKIE"))

	return &cookie.CookieFactory{
		Manager: func() cookie.IEncryptableCookieManager {
			manager := &cookie.CookieManager{
				Path:     path,
				Domain:   domain,
				Secure:   secure,
				SameSite: sameSite,
			}

			return cookie.NewEncryptableCookieManager(manager, encrypter)
		},
	}
}
