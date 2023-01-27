package jwt

import (
	"bluebell/settings"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const bluebell = "bluebell"

var (
	mySecret        = []byte("cczj")
	ErrInvalidToken = errors.New("verify Token Failed")
)

type MyClaim struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 颁发token access token 和 refresh token
func GenToken(UserID int64, Username string) (atoken, rtoken string, err error) {
	rc := jwt.RegisteredClaims{
		ExpiresAt: getJWTTime(settings.Conf.AppConfig.AtokenAt),
		Issuer:    bluebell,
		IssuedAt:  getJWTTime(0),
		NotBefore: getJWTTime(-60),
	}
	at := MyClaim{
		UserID,
		Username,
		rc,
	}
	atoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, at).SignedString(mySecret)

	// refresh token 不需要保存任何用户信息
	rt := rc
	rt.ExpiresAt = getJWTTime(settings.Conf.AppConfig.RtokenAt)
	rtoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rt).SignedString(mySecret)
	return
}

// VerifyToken 验证Token
func VerifyToken(tokenID string) (*MyClaim, error) {
	var myc = new(MyClaim)
	token, err := jwt.ParseWithClaims(tokenID, myc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return myc, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return myc, nil
}

// RefreshToken 通过 refresh token 刷新 atoken
func RefreshToken(atoken, rtoken string) (newAtoken, newRtoken string, err error) {
	// rtoken 无效直接返回
	if _, err = jwt.Parse(rtoken, keyFunc); err != nil {
		return
	}
	// 从旧access token 中解析出claims数据
	var claim MyClaim
	_, err = jwt.ParseWithClaims(atoken, &claim, keyFunc)
	// 判断错误是不是因为access token 正常过期导致的
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claim.UserID, claim.Username)
	}
	return
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

func getJWTTime(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(t) * time.Second))
}
