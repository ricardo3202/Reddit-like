package controller

// controller这一层主要就是做路由的处理、请求参数的处理、以及可能会用到的重定向

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	// 1. 获取参数，校验参数
	p := new(models.ParamSignUp)
	// 接收HTTP请求，将请求体中的 JSON 数据绑定到 p 变量上
	if err := c.ShouldBindJSON(p); err != nil { // ShouldBindJSON只能校验字段和格式
		// 请求参数有误，直接返回响应
		zap.L().Error("Signup with invalid param", zap.Error(err)) // zap.Any() 空接口类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWihMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//	//"msg": errs.Translate(trans),
		//})
		return
	}
	//手动对请求参数进行详细的业务规则校验,但是手动进行参数校验会显得啰嗦
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	//	// 请求参数有误，直接返回响应
	//	zap.L().Error("Signup with invalid param") // zap.Any() 空接口类型
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	// 2. 处理业务逻辑
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SingUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserExist)
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
}

func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及校验参数
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Signup with invalid param", zap.Error(err)) // zap.Any() 空接口类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWihMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2. 分发业务给logic，p就是已经获取到的数据
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID, // 前端能识别的数字范围是小于后端的数字范围的，于是用字符串去传值，tag里写上,sring
		"user_name": user.UserName,
		"token":     user.Token,
	})
}
