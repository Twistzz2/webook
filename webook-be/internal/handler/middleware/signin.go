package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SignInMiddlewareBuilder struct {
}

func NewSignInMiddlewareBuilder() *SignInMiddlewareBuilder {
	return &SignInMiddlewareBuilder{}
}

func (l *SignInMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signin" || ctx.Request.URL.Path == "/users/signup" {
			// 如果是登录或注册请求，直接放行
			return
		}

		session := sessions.Default(ctx)
		if session.Get("user") == nil {
			// 如果 session 中没有 user 信息，说明用户未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 如果有 user 信息，继续处理请求
		id := session.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
