package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

var jwtSecret = []byte("Jin-Vue-Admin")

type Claims struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Authority int    `json:"authority"`
	IsAdmin   int    `json:"is_admin"`
	jwt.StandardClaims
}

// tokenManager 用于管理 Token
type tokenManager struct {
	M      sync.Map      // 使用同步映射存储有效 Token 和到期时间戳
	Expire time.Duration // Token 过期时长
}

var TkMan = tokenManager{
	Expire: 24 * time.Hour,
}

// GenerateToken 签发用户Token
func GenerateToken(id uint, username string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(TkMan.Expire)
	claims := Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	// 保存 token 到合并后的 map 中
	TkMan.M.Store(token, expireTime.Unix())
	return token, nil
}

// InvalidToken 将指定的 Token 标记为失效，使其无法通过 ParseToken 验证
func InvalidToken(token string) {
	// 从 map 中删除指定的 token，使其失效
	TkMan.M.Delete(token)
}

// ParseToken 验证用户 Token 是否有效
// 如果 Token 失效，返回错误信息。
func ParseToken(token string) (*Claims, error) {
	expireTime, ok := TkMan.M.Load(token)
	if !ok {
		return nil, errors.New("invalid token")
	}
	if time.Now().Unix() > expireTime.(int64) { // Token 已过期
		return nil, errors.New("token has expired")
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || tokenClaims == nil {
		return nil, errors.New("fail to parse token: " + err.Error())
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok || !tokenClaims.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
