package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	ATokenExpiredDuration  = 600
	RTokenExpiredDuration  = 30 * 24 * 3600
	TokenIssuedAtDuration  = 0
	TokenNotBeforeDuration = -60
	TokenIssuer            = "bluebell"
	TokenSecret            = "cczj"
	ErrToken               = "verify Token Failed"
)

var (
	Now             time.Time
	mySecret        = []byte(TokenSecret)
	ErrInvalidToken = errors.New(ErrToken)
)

type MyClaim struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 颁发token access token 和 refresh token
func GenToken(UserID int64, Username string) (atoken, rtoken string, err error) {
	Now = time.Now()
	rc := jwt.RegisteredClaims{
		ExpiresAt: getJWTTime(ATokenExpiredDuration),
		Issuer:    TokenIssuer,
		IssuedAt:  getJWTTime(TokenIssuedAtDuration),
		NotBefore: getJWTTime(TokenNotBeforeDuration),
	}
	at := MyClaim{
		UserID,
		Username,
		rc,
	}
	atoken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, at).SignedString(mySecret)

	// refresh token 不需要保存任何用户信息
	rt := rc
	rt.ExpiresAt = getJWTTime(RTokenExpiredDuration)
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

func getJWTTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(Now.Add(t * time.Second))
}
