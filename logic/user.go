package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"fmt"
)

// 存放业务逻辑代码

// SignUp 注册
func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		// 数据库查询出错
		return err
	}
	// 2. 生成UID
	userID := snowflake.GenID()
	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(user)
}

// Login 登录
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.Username,
		Password: p.Password,
	}

	// 传递的是指针,于是可以拿到user.UserID
	if err := mysql.Login(user); err != nil {
		fmt.Println("这里出错了")
		return nil, err
	}
	// 生成JWT
	//return "", nil
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.Token = token
	return
}
