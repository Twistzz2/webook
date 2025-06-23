package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 定义不需要登录验证的路径
	publicPaths := map[string]bool{
		"/users/login":  true,
		"/users/signup": true,
	}

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 如果是登录或注册请求，直接放行
		if publicPaths[path] {
			ctx.Next()
			return
		}

		// 获取会话
		session := sessions.Default(ctx)

		// 如果有 user 信息，继续处理请求
		id := session.Get("userId")
		if id == nil {
			// 记录日志，便于调试
			fmt.Printf("认证失败: 会话中没有userId\n")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "请先登录",
			})
			return
		}

		// 设置会话选项
		session.Options(sessions.Options{
			Path:     "/",     // 设置路径为根路径
			MaxAge:   60 * 30, // 30分钟
			HttpOnly: true,    // 防止JavaScript访问
			// Secure: true,       // 生产环境中设置
		})
		// 刷新会话
		updateTime := session.Get("updateTime")
		now := time.Now()

		// 说明没有刷新过，刚登录
		if updateTime == nil {
			session.Set("updateTime", now.Unix())
			err := session.Save()
			if err != nil {
				fmt.Printf("保存会话失败: %v\n", err)
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			ctx.Next()
			return
		}

		// 类型断言 - 现在是int64的Unix时间戳
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			fmt.Printf("会话类型断言失败: updateTime不是int64类型\n")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 如果超过设定时间没有刷新过，刷新会话
		lastUpdateTime := time.Unix(updateTimeVal, 0)
		if now.Sub(lastUpdateTime) > time.Second*10 {
			session.Set("updateTime", now.Unix())
			err := session.Save()
			if err != nil {
				fmt.Printf("刷新会话失败: %v\n", err)
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		// 继续处理请求
		ctx.Next()
	}
}
