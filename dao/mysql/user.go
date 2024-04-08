package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go.uber.org/zap"
)

// 把每一步对数据库的操作封装成函数
// 待logic层根据业务需求调用

const sercet = "asdasd"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) error {
	sqlStr := "select count(user_id) from user where username = ?"
	count := new(int)
	if err := db.Get(count, sqlStr, username); err != nil {
		return err
	}
	if *count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行sql语句入库
	sqlStr := "insert into user(user_id, username, password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

// encryptPassword 对密码进行加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(sercet))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 登录
func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.UserName)
	if err == sql.ErrNoRows {
		// 这里是没查到
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库的过程中出错了
		return err
	}
	//查询成功则需要判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserByID 根据用户id查询用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `SELECT user_id, username FROM user WHERE user_id = ?`
	if err = db.Get(user, sqlStr, uid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("wrong userID")
			err = ErrorInvalidID
		}
	}
	return
}
