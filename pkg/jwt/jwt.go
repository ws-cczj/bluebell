package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = 6 * time.Hour
const TokenIssuedAtDuration = 0
const TokenNotBeforeDuration = -1 * time.Second
const TokenIssuer = "bluebell"

var mySecret = []byte("cczj")
var ErrorInvalidToken = errors.New("verify Token Failed")

type MyClaim struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 颁发token
func GenToken(UserID int64, Username string) (string, error) {
	t := MyClaim{
		UserID,
		Username,
		jwt.RegisteredClaims{
			ExpiresAt: getJWTTime(TokenExpireDuration),    // 过期时间
			Issuer:    TokenIssuer,                        // 签发人
			IssuedAt:  getJWTTime(TokenIssuedAtDuration),  // 签发时间
			NotBefore: getJWTTime(TokenNotBeforeDuration), // 生效时间
		},
	}
	// 生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	// 通过mySecret写上签名
	return token.SignedString(mySecret)
}

// VerifyToken 验证Token
func VerifyToken(tokenID string) (*MyClaim, error) {
	var myc = new(MyClaim)
	token, err := jwt.ParseWithClaims(tokenID, myc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrorInvalidToken
	}
	return myc, nil
}

func getJWTTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}
