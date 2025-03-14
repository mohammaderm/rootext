package authService

import (
	"strings"
	"time"

	"github.com/mohammaderm/rootext/entity"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	SignKey                  string        `koanf:"sign_key"`
	AccessExpirationTime     time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime    time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject            string        `koanf:"access_subject"`
	RefreshSubject           string        `koanf:"refresh_subject"`
	AuthMiddlewareContextKey string        `koanf:"authMiddlewareContextKey"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

type AuthService struct {
	config Config
}

func New(cfg Config) AuthService {
	return AuthService{
		config: cfg,
	}
}

func (s AuthService) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s AuthService) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s AuthService) createToken(userID uint, subject string, expireDuration time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.config.SignKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

func (s AuthService) ParseToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}

}
