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
	appID          = 1             // тестовое ID приложения, созданное миграцией
	appSecret      = "test-secret" // Секретный ключ приложения
	testAdminEmail = "test@email.test"
	testAdminPass  = "testtest"
)

// go test ./tests -count=1 -v
func TestCheckIsAdmin_HappyPath(t *testing.T) {
	ctx, st := suite.New(t) // Создаём Suite

	//login
	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    testAdminEmail,
		Password: testAdminPass,
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

	assert.Equal(t, testAdminEmail, claims["email"].(string))
	assert.Equal(t, testAdminPass, int(claims["app_id"].(float64)))

	// check is admin
	isAdmin, err := st.AuthClient.IsAdmin(ctx)

}
