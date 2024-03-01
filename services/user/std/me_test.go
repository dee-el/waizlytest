package stduserservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	commonerr "waizlytest/common/errors"

	"waizlytest/repositories"
	mockrepositories "waizlytest/repositories/mock"

	userservice "waizlytest/services/user"
	stduserservice "waizlytest/services/user/std"
)

func TestMe_GetProfile(t *testing.T) {
	type scenario struct {
		name     string
		req      func() (context.Context, int64)
		expected func() (userservice.Profile, error)

		service func() *stduserservice.StdService
	}

	scenarios := []scenario{
		{
			name: "[OK] Success",
			req: func() (context.Context, int64) {
				ctx := context.TODO()

				return ctx, 1
			},

			expected: func() (userservice.Profile, error) {
				prof := userservice.Profile{
					FullName: "Teest",
					Phone:    "123456789",
				}

				return prof, nil
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)
				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
					FullName: "Teest",
				}
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(usr, nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
		{
			name: "[FAILED] does not exist",
			req: func() (context.Context, int64) {
				ctx := context.TODO()

				return ctx, 1
			},

			expected: func() (userservice.Profile, error) {
				prof := userservice.Profile{}
				return prof, commonerr.BadRequest("user not found")
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				var usr *repositories.User
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(usr, nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			expectedResponse, expectedError := scn.expected()

			svc := scn.service()

			resp, err := svc.GetProfile(scn.req())
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

func TestMe_UpdateProfile(t *testing.T) {
	type scenario struct {
		name     string
		req      func() (context.Context, int64, userservice.UpdateRequest)
		expected func() error

		service func() *stduserservice.StdService
	}

	scenarios := []scenario{
		{
			name: "[OK] Success -- Update Phone",
			req: func() (context.Context, int64, userservice.UpdateRequest) {
				ctx := context.TODO()
				params := userservice.UpdateRequest{}

				phone := "909090"
				params.Phone = &phone

				return ctx, 1, params
			},

			expected: func() error {
				return nil
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)
				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
					FullName: "Teest",
				}
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(usr, nil)

				var emptyUsr *repositories.User
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(emptyUsr, nil)

				usrStorage.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
		{
			name: "[OK] Success -- Update Name",
			req: func() (context.Context, int64, userservice.UpdateRequest) {
				ctx := context.TODO()
				params := userservice.UpdateRequest{}

				name := "Adawd"
				params.FullName = &name

				return ctx, 1, params
			},

			expected: func() error {
				return nil
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)
				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
					FullName: "Teest",
				}
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(usr, nil)

				usrStorage.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
		{
			name: "[Failed] Conflicted Phone",
			req: func() (context.Context, int64, userservice.UpdateRequest) {
				ctx := context.TODO()
				params := userservice.UpdateRequest{}

				phone := "123123"
				params.Phone = &phone

				return ctx, 1, params
			},

			expected: func() error {
				return commonerr.BadRequest("user not found")
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				var emptyUsr *repositories.User
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(emptyUsr, nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
		{
			name: "[Failed] User does not exist",
			req: func() (context.Context, int64, userservice.UpdateRequest) {
				ctx := context.TODO()
				params := userservice.UpdateRequest{}

				phone := "123123"
				params.Phone = &phone

				return ctx, 1, params
			},

			expected: func() error {
				e := commonerr.Conflicted("")
				e.AddField("phone", fmt.Sprintf("[%s] already registered", "123123"))

				return e
			},

			service: func() *stduserservice.StdService {

				usrStorage := &mockrepositories.MockUserRepository{}

				pass, _ := bcrypt.GenerateFromPassword([]byte(`123123`), 6)
				usr := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
					FullName: "Teest",
				}
				usrStorage.On("FindUserByID", mock.Anything, mock.Anything).Return(usr, nil)

				usr2 := &repositories.User{
					ID:       1,
					Phone:    "123456789",
					Password: string(pass),
					FullName: "Teest",
				}
				usrStorage.On("FindUserByPhone", mock.Anything, mock.Anything).Return(usr2, nil)

				svc := stduserservice.NewService(usrStorage, usrStorage)

				return svc
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			expectedError := scn.expected()

			svc := scn.service()

			err := svc.UpdateProfile(scn.req())
			if err != nil {
				diff := cmp.Diff(expectedError.Error(), err.Error())
				if diff != "" {
					t.Errorf("err mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
