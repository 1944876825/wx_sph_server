package common

import (
	"errors"
	"strings"
	"time"
	"wx_video_help/db"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("w4nbeavdgWCyqwjtEN")

type UserClaims struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
	jwt.RegisteredClaims
}

func ParseToken(tokenString string) (*UserClaims, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("令牌过期")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

func GenerateToken(user *db.SphUser) (tokenString string, err error) {
	expirationTime := time.Now().Add(24 * 30 * time.Hour)
	claim := UserClaims{
		UserID:   user.ID,
		NickName: user.NickName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(jwtKey)
	return tokenString, err
}
