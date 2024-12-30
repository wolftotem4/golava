# Golava

this document is still work in progress

## Login Attempt Protection (Redis)

1. Requires redis package.

```shell
$ go get github.com/redis/go-redis/v9
```

2. Edits `internal/app/app.go`, add line `Redis *redis.Client` to **App** struct.

```go
import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wolftotem4/golava-core/golava"
)

type App struct {
	golava.App
	DB    *sqlx.DB
    
    // add this line
	Redis *redis.Client
}
```

3. Creates `internal/bootstrap/redis.go`

```go
package bootstrap

import "github.com/redis/go-redis/v9"

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
```

4. Edits `internal/bootstrap/app.go`

```go
return &app.App{
		DB:    db,

    	// add this line
		Redis: initRedis(),

		App: golava.App{
```

5. Rename `internal/ratelimit/login/redis.go.example`, `internal/ratelimit/middleware/ratelimit.go.example`, strip `.example` from their file names.
6. Congratulations!  Now, you can use `LoginRateLimit` middleware for protecting your login requests

```go
import (
    // ...
    ratemid "github.com/wolftotem4/golava/internal/ratelimit/middleware"
    // ...
)

// ...

r.POST("/login", ratemid.LoginRateLimit("username", 3, 1*time.Minute), home.SubmitLogin)
```

