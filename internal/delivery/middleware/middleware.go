package middleware

import (
	"template/pkg/auth"
	"template/pkg/database/cached"
	"template/pkg/log"
)

type Middleware struct {
	logger  *log.Logs
	jwtUtil auth.JWTUtil
	session cached.Session
}

func InitMiddleware(logger *log.Logs, util auth.JWTUtil, session cached.Session) Middleware {
	return Middleware{
		logger:  logger,
		jwtUtil: util,
		session: session,
	}
}
