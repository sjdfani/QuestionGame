package authservice

import (
	"QuestionGame/entity"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func New(signKey, accessSubject, refreshSubject string,
	accessExpirationTime, refreshExpirationTime time.Duration) Service {
	return Service{
		signKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.accessSubject, s.accessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.refreshSubject, s.refreshExpirationTime)
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.signKey), nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else if claims, ok := token.Claims.(*Claims); ok {
		fmt.Println(claims.UserID, claims.RegisteredClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) createToken(userID uint, subject string, expirationTime time.Duration) (string, error) {
	// create a signer for rsa 256
	// TODO: replace with RS 256 later

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.signKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	// Create token string
	return tokenString, nil
}
