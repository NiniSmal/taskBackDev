package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"taskBackDev/entity"
	"time"
)

type Storage struct {
	conn *pgxpool.Pool
}

func NewStorage(conn *pgxpool.Pool) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) SaveToken(ctx context.Context, session entity.Session) error {
	query := "INSERT INTO sessions(user_id,user_name, refresh_token, expires_at, is_used) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.conn.Exec(ctx, query, session.UserID, session.UserName, session.RefreshToken, session.ExpiresAt, false)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateToken(ctx context.Context, session entity.Session) error {
	query := "UPDATE  sessions SET user_id=$1, user_name=$2, refresh_token=$3 ,expires_at=$4, is_used=$5 WHERE expires_at < $6"
	_, err := s.conn.Exec(ctx, query, session.UserID, session.UserName, session.RefreshToken, session.ExpiresAt, true, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) TokenByUserID(ctx context.Context, id int64) (entity.Session, error) {
	query := "SELECT user_id, user_name, refresh_token, expires_at, is_used FROM sessions WHERE user_id= $1"
	var session entity.Session

	err := s.conn.QueryRow(ctx, query, id).Scan(&session.UserID, &session.UserName, &session.RefreshToken, &session.ExpiresAt, &session.IsUsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, entity.ErrNotFound
		}

		return entity.Session{}, err
	}
	return session, nil
}

func (s *Storage) TokenByUserIP(ctx context.Context, name string) (entity.Session, error) {
	query := "SELECT user_id, user_name, refresh_token, expires_at, is_used FROM sessions WHERE user_name= $1"
	var session entity.Session

	err := s.conn.QueryRow(ctx, query, name).Scan(&session.UserID, &session.UserName, &session.RefreshToken, &session.ExpiresAt, &session.IsUsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Session{}, entity.ErrNotFound
		}

		return entity.Session{}, err
	}
	return session, nil
}
