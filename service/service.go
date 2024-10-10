package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
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

func (s *Service) NewJWT(name string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":         name,
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

func (s *Service) Tokens(ctx context.Context, name string, id int64) (entity.TokenResponse, error) {
	session, err := s.storage.TokenByUserID(ctx, id)
	if err != nil {
		return entity.TokenResponse{}, err
	}

	at, err := s.NewJWT(name, s.cfg.AccessTokenTTL)
	if err != nil {
		return entity.TokenResponse{}, fmt.Errorf("%s: %w", "Get Access Token", err)
	}
	accessToken, err := strconv.ParseInt(at, 10, 64)
	if err != nil {
		return entity.TokenResponse{}, err
	}
	tokenResponse := entity.TokenResponse{
		AccessToken:  session.RefreshToken,
		RefreshToken: accessToken,
	}
	return tokenResponse, nil
}

func (s *Service) CreateToken(ctx context.Context, userID int64, ipName string) (entity.TokenResponse, error) {
	refresh, err := s.NewRefreshToken()
	if err != nil {
		return entity.TokenResponse{}, err
	}
	session := entity.Session{
		UserID:       userID,
		UserName:     ipName,
		RefreshToken: refresh,
		ExpiresAt:    time.Now().Add(s.cfg.RefreshTokenTTL),
		IsUsed:       false,
	}
	err = s.storage.SaveToken(ctx, session)
	if err != nil {
		return entity.TokenResponse{}, err
	}

	at, err := s.NewJWT(ipName, s.cfg.AccessTokenTTL)
	if err != nil {
		return entity.TokenResponse{}, fmt.Errorf("%s: %w", "Get Access Token", err)
	}

	accessToken, err := strconv.ParseInt(at, 10, 64)
	if err != nil {
		return entity.TokenResponse{}, err
	}
	tokenResponse := entity.TokenResponse{
		AccessToken:  session.RefreshToken,
		RefreshToken: accessToken,
	}
	return tokenResponse, nil
}

func (s *Service) SwitchToken(ctx context.Context, tokenCookie string, ipName string) (int64, error) {
	oldToken, err := s.storage.TokenByUserIP(ctx, ipName)
	if err != nil {
		return 0, err
	}
	err = checkRefreshToken(tokenCookie, oldToken.RefreshToken)
	if err != nil {
		return 0, err
	}

	err = s.storage.UpdateToken(ctx, oldToken)
	if err != nil {
		return 0, err
	}
	return oldToken.UserID, nil
}
