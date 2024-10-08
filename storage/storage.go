package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"taskBackDev/entity"
)

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) SetSession(ctx context.Context, id uuid.UUID, session entity.Session) error {
	query := "INSERT INTO  sessions (user_id, refresh_token ,expires_at) VALUES ($1, $2, $3)"
	_, err := s.conn.Exec(ctx, query, id, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UserIDByRefreshToken(ctx context.Context, refreshToken string) (entity.Session, error) {
	query := "SELECT user_id, refresh_token ,expires_at FROM sessions WHERE refresh_token = $1 "

	var session entity.Session
	err := s.conn.QueryRow(ctx, query, refreshToken).Scan(&session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (s *Storage) TokenByUserID(ctx context.Context, id string) (entity.Session, error) {
	query := "SELECT user_id,refresh_token, expires_at FROM sessions WHERE user_id= $1"
	var session entity.Session

	err := s.conn.QueryRow(ctx, query, id).Scan(&session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return entity.Session{}, err
	}
	return session, nil
}
