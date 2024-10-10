package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"taskBackDev/db_test"
	"taskBackDev/entity"
	"testing"
	"time"
)

func TestRepositorySaveToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conn := db_test.DBConnection(t)
	st := NewStorage(conn)

	expiresAt := time.Now().UTC().Round(time.Millisecond)
	session := entity.Session{
		UserID:       1,
		UserName:     uuid.NewString(),
		RefreshToken: uuid.NewString(),
		ExpiresAt:    expiresAt,
		IsUsed:       false,
	}
	err := st.SaveToken(ctx, session)
	require.NoError(t, err)

	sessionDB, err := st.TokenByUserID(ctx, session.UserID)
	require.Error(t, err)
	require.Equal(t, session, sessionDB)
}

func TestRepositoryUpdateToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	conn := db_test.DBConnection(t)
	st := NewStorage(conn)

	expiresAt := time.Now().UTC().Round(time.Millisecond)
	session := entity.Session{
		UserID:       1,
		UserName:     uuid.NewString(),
		RefreshToken: uuid.NewString(),
		ExpiresAt:    expiresAt,
		IsUsed:       false,
	}
	err := st.SaveToken(ctx, session)
	require.NoError(t, err)
	session2 := entity.Session{
		UserID:       session.UserID,
		UserName:     session.UserName,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		IsUsed:       true,
	}
	err = st.UpdateToken(ctx, session2)
	require.NoError(t, err)
	sessionDB, err := st.TokenByUserIP(ctx, session2.UserName)
	require.NoError(t, err)
	require.Equal(t, session, sessionDB)

}
