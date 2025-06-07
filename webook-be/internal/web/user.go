package web

import "github.com/gin-gonic/gin"

// UserHandler 在这里定义和用户有关的路由
type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(router *gin.Engine) {
	ug := router.Group("/users")
	ug.POST("/signup", u.SignUp)  // 用户注册
	ug.POST("/signin", u.SignIn)  // 用户登录
	ug.POST("/edit", u.Edit)      // 用户信息编辑
	ug.GET("/profile", u.Profile) // 用户个人信息查看
}

func (u *UserHandler) SignUp(c *gin.Context) {
	// TODO: 处理用户注册逻辑
}

func (u *UserHandler) SignIn(c *gin.Context) {
	// TODO: 处理用户登录逻辑
}

func (u *UserHandler) Edit(c *gin.Context) {
	// TODO: 处理用户信息编辑逻辑
}

func (u *UserHandler) Profile(c *gin.Context) {
	// TODO: 处理用户个人信息查看逻辑
}
