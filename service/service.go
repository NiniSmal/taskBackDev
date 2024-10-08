package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"taskBackDev/config"
	"taskBackDev/entity"
	"taskBackDev/storage"
	"time"
)

type Service struct {
	storage *storage.Storage
	cfg     *config.Config
}

func NewService(storage *storage.Storage, cfg *config.Config) *Service {
	return &Service{storage: storage,
		cfg: cfg,
	}
}

func (s *Service) SingIn(ctx context.Context, id uuid.UUID) (entity.TokenResponse, error) {
	return s.createSession(ctx, id)
}

func (s *Service) createSession(ctx context.Context, id uuid.UUID) (entity.TokenResponse, error) {
	var (
		resp entity.TokenResponse
		err  error
	)

	resp.AccessToken, err = s.NewJWT(id.String(), AccessTokenTTL)
	if err != nil {
		return entity.TokenResponse{}, err
	}
	resp.RefreshToken, err = s.NewRefreshToken()
	if err != nil {
		return entity.TokenResponse{}, err
	}

	session := entity.Session{
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 48),
	}
	err = s.storage.SetSession(ctx, id, session)
	if err != nil {
		return entity.TokenResponse{}, err
	}
	return resp, nil
}

var AccessTokenTTL time.Duration

func (s *Service) GetAccessToken(userID string) (string, error) {
	accessToken, err := s.NewJWT(userID, AccessTokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "Get Access Token", err)
	}

	return accessToken, nil
}

func (s *Service) NewJWT(id string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"id":         id,
		"expires_at": jwt.NewNumericDate(time.Now().Add(ttl)),
	})

	return token.SignedString([]byte(entity.JwtSecretKey))
}

func (s *Service) NewRefreshToken() (string, error) {
	refreshToken := uuid.NewString()
	refreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func hashRefreshToken(refreshToken string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkRefreshToken(refreshToken, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(refreshToken))
	return err
}

func (s *Service) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(entity.JwtSecretKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (s *Service) RefreshTokens(ctx context.Context, id string) (entity.TokenResponse, error) {
	session, err := s.storage.TokenByUserID(ctx, id)
	if err != nil {
		return entity.TokenResponse{}, err
	}
	accessToken, err := s.NewJWT(id, AccessTokenTTL)
	if err != nil {
		return entity.TokenResponse{}, fmt.Errorf("%s: %w", "Get Access Token", err)
	}
	tokenResponse := entity.TokenResponse{
		AccessToken:  session.RefreshToken,
		RefreshToken: accessToken,
	}
	return tokenResponse, nil
}
func (s *Service) SwitchToken(ctx context.Context, newToken string, userID string) error {
	oldToken, err := s.storage.TokenByUserID(ctx, userID)
	if err != nil {
		return err
	}
	err = checkRefreshToken(newToken, oldToken.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}
