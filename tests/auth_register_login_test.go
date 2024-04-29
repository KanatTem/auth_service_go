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
	appID      = 1             // тестовое ID приложения, созданное миграцией
	appSecret  = "test-secret" // Секретный ключ приложения
	emptyAppID = 0
	passLen    = 10
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

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()

	pass := gofakeit.Password(true, true, true, true, false, passLen)

	//first register
	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	//register check
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())
	//second register (error)
	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "missing password",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passLen),
			expectedErr: "missing email",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "missing email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "missing password",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passLen),
			appID:       appID,
			expectedErr: "missing email",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "missing email",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passLen),
			appID:       appID,
			expectedErr: "invalid credentials",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passLen),
			appID:       emptyAppID,
			expectedErr: "missing appID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, passLen),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
