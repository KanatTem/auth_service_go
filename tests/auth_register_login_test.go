package tests

import (
	"auth_service/tests/suite"
	ssov1 "github.com/KanatTem/sso_proto/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	appID     = 1             // тестовое ID приложения, созданное миграцией
	appSecret = "test-secret" // Секретный ключ приложения
	passLen   = 10
)

// go test ./tests -count=1 -v
func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, passLen)

	//register
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	//register check
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	//login
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})

	//logic check
	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["userId"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}
