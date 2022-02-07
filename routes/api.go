package routes

import (
	"gohub/app/http/controllers/api/v1/auth"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	v1 := r.Group("v1")
	{
		// 注册认证路由
		authGroup := v1.Group("/auth")
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
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			// 发送手机验证码
			authGroup.POST("/verify-codes/phone", vcc.SendUsingPhone)
			// 发送邮件验证码
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)
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
	}
}
