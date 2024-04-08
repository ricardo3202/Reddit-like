package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("加盐操作")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成 JWT token
func GenToken(userID int64, username string) (string, error) {
	// 创建一个自定义的数据
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				viper.GetDuration("auth.jwt_expire") * time.Hour).Unix(), // 过期时间
			Issuer: "bluebell", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c) // 这里的HS256这个字段是一个加密算法
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret) // 除了自己的数据还混淆了一段奇怪的东西
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	mc := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, err
	}
	return nil, errors.New("invalid token")
}
