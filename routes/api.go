package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("v1")

	// 全局限流中间件：每小时限流。这里是所有 API （根据 IP）请求加起来。
	// 作为参考 Github API 每小时最多 60 个请求（根据 IP）。
	// 测试时，可以调高一点。
	v1.Use(middlewares.LimitIP("200-H"))
	{
		// 注册认证路由
		authGroup := v1.Group("/auth")
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断 Email 是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// 手机号注册
			authGroup.POST("/signup/using-phone", suc.SignupUsingPhone)
			// 邮箱注册
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)
			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("20-H"), vcc.ShowCaptcha)
			// 发送手机验证码
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			// 发送邮件验证码
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)
			// 登录
			lc := new(auth.LoginController)
			// 使用手机号，短信验证码进行登录
			authGroup.POST("/login/using-phone", lc.LoginByPhone)
			// 支持手机号，Email 和 用户名
			authGroup.POST("/login/using-password", lc.LoginByPassword)
			// 刷新 token
			authGroup.POST("/login/refresh-token", lc.RefreshToken)
			// 重置密码
			pc := new(auth.PasswordController)
			// 使用手机重置密码
			authGroup.POST("/password-reset/using-phone", pc.ResetByPhone)
			// 使用邮件重置密码
			authGroup.POST("/password-reset/using-email", pc.ResetByEmail)
		}

		uc := new(controllers.UsersController)
		// 获取当前用户
		v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("", uc.Index)
		}
	}
}
