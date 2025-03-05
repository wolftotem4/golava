package middlewares

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wolftotem4/golava-core/instance"
	sessmid "github.com/wolftotem4/golava-core/session/middleware"
	"github.com/wolftotem4/golava/internal/app"
	"github.com/wolftotem4/golava/internal/bootstrap"
)

type startSession struct {
	mu sync.Mutex
	gin.HandlerFunc
}

func (s *startSession) StartSession(c *gin.Context) {
	if s.HandlerFunc == nil {
		var (
			i = instance.MustGetInstance(c)
			a = i.App.(*app.App)
		)

		err := s.makeMiddleware(a)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	s.HandlerFunc(c)
}

func (s *startSession) makeMiddleware(a *app.App) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.HandlerFunc == nil {
		factory, err := bootstrap.InitSession(a.DB, "sessions")
		if err != nil {
			return err
		}

		s.HandlerFunc = sessmid.StartSession(factory)
	}

	return nil
}

var StartSession gin.HandlerFunc = (&startSession{}).StartSession
