package tests

import (
	"SSO/pkg/utils"
	"SSO/tests/suite"
	"fmt"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	err := godotenv.Load("test.env") // Два уровня вверх от ./cmd/sso
	if err != nil {
		panic(err)
	}

	ctx, st := suite.New(t)

	firstName := gofakeit.Name()
	lastName := gofakeit.Name()
	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  pass,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, respReg.GetUserId())
	fmt.Println(respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Username: email,
		Password: pass,
		AppLogin: appID,
	})

	require.NoError(t, err)

	token := respLogin.GetAccessToken()
	assert.NotEmpty(t, token)

	fmt.Println(os.Getenv("JWT_SECRET_KEY"))
	tokenClaims, err := utils.ParseToken(token, os.Getenv("JWT_SECRET_KEY"))

	require.NoError(t, err)

	assert.Equal(t, appID, tokenClaims.AppID)
	assert.Equal(t, username, tokenClaims.Username)
	assert.Equal(t, respReg.GetUserId(), int64(tokenClaims.UserID))
}

func TestRegister_DuplicatedReg(t *testing.T) {
	ctx, st := suite.New(t)

	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	username := gofakeit.Username()
	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  pass,
	})
	assert.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		firstName   string
		lastName    string
		username    string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			username:    gofakeit.Username(),
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			username:    gofakeit.Username(),
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Both Empty",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			username:    gofakeit.Username(),
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
		{
			name:        "Register with empty first name",
			firstName:   "",
			lastName:    gofakeit.LastName(),
			username:    gofakeit.Username(),
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			expectedErr: "first name is required",
		},
		{
			name:        "Register with empty last name",
			firstName:   gofakeit.FirstName(),
			lastName:    "",
			username:    gofakeit.Username(),
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			expectedErr: "last name is required",
		},
		{
			name:        "Register with empty username",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			username:    "",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			expectedErr: "username is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
				Username:  tt.username,
				Email:     tt.email,
				Password:  tt.password,
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
		username    string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			username:    gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Username",
			username:    "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "username is required",
		},
		{
			name:        "Login with Both Empty Username and Password",
			username:    "",
			password:    "",
			appID:       appID,
			expectedErr: "username is required",
		},
		{
			name:        "Login with Non-Matching Password",
			username:    gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid username or password",
		},
		{
			name:        "Login without AppID",
			username:    gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Username:  gofakeit.Username(),
				Email:     gofakeit.Email(),
				Password:  randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Username: tt.username,
				Password: tt.password,
				AppLogin: tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
