package jwtauthservice_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	commonerr "waizlytest/common/errors"

	"waizlytest/repositories"
	mockrepositories "waizlytest/repositories/mock"

	authservice "waizlytest/services/auth"
	jwtauthservice "waizlytest/services/auth/jwt"
)

var publicCert = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCOqXTycDGJ50wv+hXPyj0dzDwQ
tB6jQAYC2gqRc6VN29kolhQG3xfbz/elrmfg7HPuBONZRmOGJETtSXQQ1Ld4IXAZ
Lt/gbkfA3bHPa6z2UHrqe/fOpHRNHTl56fLblJjDkH5Vipow3xlm4wFrSTAfp8WA
zK8po+6Ye19K4gXqLQIDAQAB
-----END PUBLIC KEY-----
`

var privateCert = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCOqXTycDGJ50wv+hXPyj0dzDwQtB6jQAYC2gqRc6VN29kolhQG
3xfbz/elrmfg7HPuBONZRmOGJETtSXQQ1Ld4IXAZLt/gbkfA3bHPa6z2UHrqe/fO
pHRNHTl56fLblJjDkH5Vipow3xlm4wFrSTAfp8WAzK8po+6Ye19K4gXqLQIDAQAB
AoGAJvlEabcc0X/O4IyByPKHH8zb2/RZKmAjREQs/u+JCWw2N1BIyFfKPNLj5O9w
kZIHWc8cDReduNfPvMOEYdS7Cj4BbnUEiOVNgX4wKH62ukPqBLzOBJ8LNc3LcGuz
EOm/QTCA8T5+NVQgvJfylUoOChwa/8j52WwkmKOuChkL5EECQQDKYkJKh+wQqU2o
V1AzQeFfQauZl+Topzkv3DRQz5zAwuqAXxb50BWvU8Ya5Uk9Jb6aMfXdAtzXlYQi
nQwV4+b7AkEAtHTWzLylTCP1Nu5rVUFQ9UfDegN3ndV5JWWIhwO5Cg2DRopJU3xE
j8VjkwuzXh+YgB6ydSPqVskcdtyWFqDK9wJBAJ+D2Oozvc3YE7x2rWDpMUDKWv2h
qivx+fIOJzH2oX+RYhGyYUKfTyg06HU0eYh3ooaYkEgVxIkrcM1zaR4r1lcCQCpG
/VP2FlviSL7X2LmeldMBPyDE0y9dJgbG5NeM3bsnM0xBdbBjesScIBoBMcKpcFcD
2hdrlKlEcLDAOUGP5j0CQD2GjEumeGtExzxX7HVAw/zl3RRPWg+YVDHLJrfg6iXl
c7C6DjFC0O5Kp32ysYKPkJQrObMfux5cyzfGib2OgX8=
-----END RSA PRIVATE KEY-----
`

func TestJWTAuth_Login(t *testing.T) {
	type scenario struct {
		name     string
		req      func() (context.Context, authservice.LoginRequest)
		expected func() (authservice.LoginResponse, error)

		service func() *jwtauthservice.JWTAuth
	}

	scenarios := []scenario{
		{
			name: "[OK] Success",
			req: func() (context.Context, authservice.LoginRequest) {
				ctx := context.TODO()

				req := authservice.LoginRequest{
					Password: "123123",
					Phone:    "321321",
				}

				return ctx, req
			},

			expected: func() (authservice.LoginResponse, error) {
				resp := authservice.LoginResponse{
					ID:    1,
					Token: "",
				}

				return resp, nil
			},

			service: func() *jwtauthservice.JWTAuth {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)

				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
				}
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(usr, nil)
				usrStorage.On("CreateUserAttendance", mock.Anything, mock.AnythingOfType("*repositories.UserAttendance")).Return(nil)
				usrStorage.On("SaveUserAttendanceSummary", mock.Anything, mock.Anything).Return(nil)

				svc, err := jwtauthservice.NewJWTAuth(usrStorage, usrStorage, privateCert, publicCert)
				if err != nil {
					t.Errorf("failed instatiate service: %v\n", err)
					return nil
				}

				return svc
			},
		},
		{
			name: "[OK] User does not exist",
			req: func() (context.Context, authservice.LoginRequest) {
				ctx := context.TODO()

				req := authservice.LoginRequest{
					Password: "123123",
					Phone:    "321321",
				}

				return ctx, req
			},

			expected: func() (authservice.LoginResponse, error) {
				return authservice.LoginResponse{}, commonerr.BadRequest("phone or password is wrong")
			},

			service: func() *jwtauthservice.JWTAuth {
				usrStorage := &mockrepositories.MockUserRepository{}

				var usr *repositories.User
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(usr, nil)

				svc, err := jwtauthservice.NewJWTAuth(usrStorage, usrStorage, privateCert, publicCert)
				if err != nil {
					t.Errorf("failed instatiate service: %v\n", err)
					return nil
				}

				return svc
			},
		},
		{
			name: "[Failed] Wrong Password",
			req: func() (context.Context, authservice.LoginRequest) {
				ctx := context.TODO()

				req := authservice.LoginRequest{
					Password: "123123",
					Phone:    "321321",
				}

				return ctx, req
			},

			expected: func() (authservice.LoginResponse, error) {
				return authservice.LoginResponse{}, commonerr.BadRequest("phone or password is wrong")
			},

			service: func() *jwtauthservice.JWTAuth {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`12312311`), 6)

				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
				}
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(usr, nil)

				usrStorage.On("CreateUserAttendance", mock.Anything, mock.AnythingOfType("*repositories.UserAttendance")).Return(nil)
				usrStorage.On("SaveUserAttendanceSummary", mock.Anything, mock.Anything).Return(nil)

				svc, err := jwtauthservice.NewJWTAuth(usrStorage, usrStorage, privateCert, publicCert)
				if err != nil {
					t.Errorf("failed instatiate service: %v\n", err)
					return nil
				}

				return svc
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			expectedResponse, expectedError := scn.expected()

			svc := scn.service()

			resp, err := svc.Login(scn.req())
			if err != nil {
				diff := cmp.Diff(expectedError.Error(), err.Error())
				if diff != "" {
					t.Errorf("err mismatch (-want +got):\n%s", diff)
				}
			}

			diff := cmp.Diff(expectedResponse, resp, cmpopts.IgnoreFields(authservice.LoginResponse{}, "Token"))
			if diff != "" {
				t.Errorf("resp mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestJWTAuth_ValidateToken(t *testing.T) {
	type scenario struct {
		name     string
		req      func() string
		expected func() (int64, error)

		service func() *jwtauthservice.JWTAuth
	}

	scenarios := []scenario{
		{
			name: "[OK] Success",
			req: func() string {
				return "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkyMDQwNzMsImlhdCI6MTcwOTIwMDQ3MywibmFtZSI6IiIsInN1YiI6IjEifQ.AOqYJqUsFKKvb_r9OUufkdZFwfkDNanWjORjP5ptFLibic4IyUuXwA9Fnffnr8205PNC8r_ebADfbsPxw1y_jPI556_cU9xxDxC3wQLPNdP80Oi-6I--aMIIsML0xXwwSqrRiSepmNpuuQM9ZMkZwtxJOXGvnYlPiP37hx9K8-k"
			},

			expected: func() (int64, error) {
				return 1, nil
			},

			service: func() *jwtauthservice.JWTAuth {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)

				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
				}
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(usr, nil)

				svc, err := jwtauthservice.NewJWTAuth(usrStorage, usrStorage, privateCert, publicCert)
				if err != nil {
					t.Errorf("failed instatiate service: %v\n", err)
					return nil
				}

				return svc
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			expectedResponse, expectedError := scn.expected()

			svc := scn.service()

			resp, err := svc.ValidateToken(scn.req())
			if err != nil {
				diff := cmp.Diff(expectedError.Error(), err.Error())
				if diff != "" {
					t.Errorf("err mismatch (-want +got):\n%s", diff)
				}
			}

			diff := cmp.Diff(expectedResponse, resp)
			if diff != "" {
				t.Errorf("resp mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
