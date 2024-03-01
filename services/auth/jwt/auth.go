package jwtauthservice

import (
	"context"
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	commonerr "waizlytest/common/errors"

	"waizlytest/repositories"
	authservice "waizlytest/services/auth"
)

var _ (authservice.Authn) = (*JWTAuth)(nil)
var _ (authservice.Authz) = (*JWTAuth)(nil)

type JWTAuth struct {
	userWriter repositories.UserWriter
	userReader repositories.UserReader

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTAuth(
	userWriter repositories.UserWriter,
	userReader repositories.UserReader,
	privateKey, publicKey string,
) (*JWTAuth, error) {
	pem, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	cert, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, err
	}

	instance := &JWTAuth{
		userWriter: userWriter,
		userReader: userReader,
		privateKey: pem,
		publicKey:  cert,
	}

	return instance, nil
}

func (s *JWTAuth) createToken(usr *repositories.User, now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  strconv.FormatInt(usr.ID, 10),
		"name": usr.FullName,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JWTAuth) ValidateToken(token string) (int64, error) {
	jot, err := s.getJWT(token)
	if err != nil {
		return 0, err
	}

	if jot == nil {
		return 0, fmt.Errorf("token is not valid")
	}

	if !jot.Valid {
		return 0, fmt.Errorf("token is not valid")
	}

	return s.getClaimedID(jot)
}

func (s *JWTAuth) Login(ctx context.Context, params authservice.LoginRequest) (authservice.LoginResponse, error) {
	lr := authservice.LoginResponse{}

	usr, err := s.userReader.FindUserByPhone(ctx, params.Phone)
	if err != nil {
		return lr, err
	}

	if usr == nil {
		e := commonerr.BadRequest("phone or password is wrong")
		return lr, e
	}

	// compare
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(params.Password))
	if err != nil {
		e := commonerr.BadRequest("phone or password is wrong")
		return lr, e
	}

	now := time.Now()
	err = s.userWriter.CreateUserAttendance(ctx, &repositories.UserAttendance{
		UserID:  usr.ID,
		LoginAt: now,
	})
	if err != nil {
		return lr, err
	}

	err = s.userWriter.SaveUserAttendanceSummary(ctx, usr.ID)
	if err != nil {
		return lr, err
	}

	token, err := s.createToken(usr, now)
	if err != nil {
		return lr, err
	}

	lr.ID = usr.ID
	lr.Token = token
	return lr, nil
}

func (s *JWTAuth) getJWT(token string) (*jwt.Token, error) {
	jot, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return jot, nil
}

func (s *JWTAuth) getClaimedID(jot *jwt.Token) (int64, error) {
	var id int64 = 0
	claims, ok := jot.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("token is not valid")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, fmt.Errorf("token is not valid")
	}

	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
