package ratelimit

import "errors"

var ErrTooManyAttempts = errors.New("too many attempts")
