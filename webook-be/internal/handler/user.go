package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Twistzz2/webook/webook-be/internal/domain"
	"github.com/Twistzz2/webook/webook-be/internal/repository/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"

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
	ug.POST("/login", u.Login)    // 用户登录
	ug.POST("/logout", u.Logout)  // 用户登出
	ug.POST("/edit", u.Edit)      // 用户信息编辑
	ug.GET("/profile", u.Profile) // 用户个人信息查看
}

func (u *UserHandler) SignUp(c *gin.Context) {
	// TODO: 处理用户注册逻辑

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
	// 最佳实践
	// errors.Is(err, service.ErrEmailAlreadyExists) // 判断错误类型
	if err == service.ErrEmailAlreadyExists {
		c.String(http.StatusBadRequest, "邮箱已被注册") // 400
		return
	}
	if err != nil {
		c.String(http.StatusInternalServerError, "系统错误") // 500
		return
	}

	c.String(http.StatusOK, "注册成功")
	fmt.Printf("%+v\n", req)

}

func (u *UserHandler) Login(c *gin.Context) {
	// TODO: 处理用户登录逻辑
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusBadRequest, "请求数据格式错误") // 400
		return
	}

	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		c.String(http.StatusUnauthorized, "账号或密码无效") // 401
		return
	}
	if err != nil {
		// 记录具体错误以便调试
		fmt.Printf("登录服务错误: %v\n", err)
		c.String(http.StatusInternalServerError, "系统错误") // 500
		return
	}
	// 这里登录成功后，需要将 session 拿到
	session := sessions.Default(c)
	// 先设置会话选项
	session.Options(sessions.Options{
		Path:     "/",   // 设置路径为根路径
		MaxAge:   86400, // 1天有效期
		HttpOnly: true,  // 防止JavaScript访问
		// Secure: true,      // 生产环境中设置
	}) // 使用userId作为会话键
	session.Set("userId", user.Id)
	// 设置更新时间，配合中间件使用 - 现在可以直接存储 time.Time
	session.Set("updateTime", time.Now())
	// 保存会话，使用新变量避免冲突
	if saveErr := session.Save(); saveErr != nil {
		fmt.Printf("保存会话错误: %v\n", saveErr)
		c.String(http.StatusInternalServerError, "系统错误: "+saveErr.Error())
		return
	}

	c.String(http.StatusOK, "登录成功") // 200
	return
}

func (u *UserHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{
		MaxAge: -1, // 设置 MaxAge 为 -1 删除 cookie
	})
	session.Clear() // 清除 session 中的所有数据

	session.Save()                  // 保存 session 的更改
	c.String(http.StatusOK, "登出成功") // 200
	return
}

func (u *UserHandler) Edit(c *gin.Context) {
	// TODO: 处理用户信息编辑逻辑
}

func (u *UserHandler) Profile(c *gin.Context) {
	// 定义返回结构
	type Profile struct {
		Email string `json:"email"`
	}

	// 获取会话
	session := sessions.Default(c)

	// 从会话中获取用户ID
	userId := session.Get("userId")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "请先登录",
		})
		return
	}

	// 类型断言
	id, ok := userId.(int64)
	if !ok {
		// 记录实际类型，便于调试
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "会话数据类型错误",
		})
		return
	}

	// 获取用户资料
	user, err := u.svc.Profile(c, id)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "系统错误",
			"error":   err.Error(),
		})
		return
	}

	// 返回用户资料
	c.JSON(http.StatusOK, Profile{
		Email: user.Email,
	})
}
