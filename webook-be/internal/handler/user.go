package handler

import (
	"fmt"
	"net/http"

	"github.com/Twistzz2/webook/webook-be/internal/domain"
	"github.com/Twistzz2/webook/webook-be/internal/repository/service"
	regexp "github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
)

// UserHandler 在这里定义和用户有关的路由
type UserHandler struct {
	svc              *service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		svc:              svc,
		emailRegexExp:    emailExp,
		passwordRegexExp: passwordExp,
	}
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
	c.String(http.StatusOK, "你好，这里是注册接口")

	// 定义方法内部类 SignUpReq 来接收数据
	// 除了 SignUp 方法外，其他方法都用不了
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq

	// Bind 方法根据 Content-Type 自动选择绑定方式
	// 解析失败返回 400
	if err := c.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailRegexExp.MatchString(req.Email)
	if err != nil {
		// TODO: 记录日志
		c.String(http.StatusInternalServerError, "服务器错误") // 500
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "邮箱格式不正确") // 400
		return
	}

	if req.Password != req.ConfirmPassword {
		c.String(http.StatusBadRequest, "两次输入的密码不一致") // 400
		return
	}

	ok, err = u.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		// TODO: 记录日志
		c.String(http.StatusInternalServerError, "服务器错误") // 500
		return
	}
	if !ok {
		c.String(http.StatusBadRequest, "密码格式不正确") // 400
		return
	}

	// 调用 svc 的方法
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误") // 500
		return
	}

	c.String(http.StatusOK, "注册成功")
	fmt.Printf("%+v\n", req)

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
