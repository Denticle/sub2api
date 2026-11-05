package oauthprovider

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type authorizationCode struct {
	ID            int64
	ClientID      string
	UserID        int64
	RedirectURI   string
	CodeChallenge string
	ExpiresAt     time.Time
	ConsumedAt    *time.Time
}

type accessToken struct {
	ID        int64
	ClientID  string
	UserID    int64
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// repository 是本模块的全部持久化；只用标准库 database/sql，和 securityaudit 的做法一致。
type repository interface {
	createCode(ctx context.Context, codeHash string, code authorizationCode) error
	// consumeCode 原子地把授权码标记为已用，返回标记前的记录。
	// 已经被用过时返回记录且 ConsumedAt 非空，调用方据此判定重放。
	consumeCode(ctx context.Context, codeHash string, now time.Time) (*authorizationCode, error)
	createToken(ctx context.Context, tokenHash string, token accessToken, codeID int64) error
	revokeTokensByCode(ctx context.Context, codeID int64, now time.Time) error
	getToken(ctx context.Context, tokenHash string) (*accessToken, error)

	hasConsent(ctx context.Context, userID int64, clientID string) (bool, error)
	saveConsent(ctx context.Context, userID int64, clientID string) error

	getClientKeyID(ctx context.Context, userID int64, clientID string) (int64, error)
	setClientKeyID(ctx context.Context, userID int64, clientID string, apiKeyID int64) error

	deleteExpired(ctx context.Context, before time.Time) error
}

type sqlRepository struct{ db *sql.DB }

func newSQLRepository(db *sql.DB) *sqlRepository { return &sqlRepository{db: db} }

func (r *sqlRepository) createCode(ctx context.Context, codeHash string, c authorizationCode) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO oauth_authorization_codes (code_hash, client_id, user_id, redirect_uri, code_challenge, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		codeHash, c.ClientID, c.UserID, c.RedirectURI, c.CodeChallenge, c.ExpiresAt)
	return err
}

func (r *sqlRepository) consumeCode(ctx context.Context, codeHash string, now time.Time) (*authorizationCode, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	var c authorizationCode
	var consumed sql.NullTime
	err = tx.QueryRowContext(ctx, `
		SELECT id, client_id, user_id, redirect_uri, code_challenge, expires_at, consumed_at
		FROM oauth_authorization_codes WHERE code_hash = $1 FOR UPDATE`, codeHash).
		Scan(&c.ID, &c.ClientID, &c.UserID, &c.RedirectURI, &c.CodeChallenge, &c.ExpiresAt, &consumed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if consumed.Valid {
		t := consumed.Time
		c.ConsumedAt = &t
		return &c, tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, `UPDATE oauth_authorization_codes SET consumed_at = $2 WHERE id = $1`, c.ID, now); err != nil {
		return nil, err
	}
	return &c, tx.Commit()
}

func (r *sqlRepository) createToken(ctx context.Context, tokenHash string, t accessToken, codeID int64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO oauth_access_tokens (token_hash, client_id, user_id, authorization_code_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)`,
		tokenHash, t.ClientID, t.UserID, codeID, t.ExpiresAt)
	return err
}

func (r *sqlRepository) revokeTokensByCode(ctx context.Context, codeID int64, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE oauth_access_tokens SET revoked_at = $2
		WHERE authorization_code_id = $1 AND revoked_at IS NULL`, codeID, now)
	return err
}

func (r *sqlRepository) getToken(ctx context.Context, tokenHash string) (*accessToken, error) {
	var t accessToken
	var revoked sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, client_id, user_id, expires_at, revoked_at
		FROM oauth_access_tokens WHERE token_hash = $1`, tokenHash).
		Scan(&t.ID, &t.ClientID, &t.UserID, &t.ExpiresAt, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if revoked.Valid {
		v := revoked.Time
		t.RevokedAt = &v
	}
	return &t, nil
}

func (r *sqlRepository) hasConsent(ctx context.Context, userID int64, clientID string) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM oauth_consents WHERE user_id = $1 AND client_id = $2`, userID, clientID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *sqlRepository) saveConsent(ctx context.Context, userID int64, clientID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO oauth_consents (user_id, client_id) VALUES ($1, $2)
		ON CONFLICT (user_id, client_id) DO NOTHING`, userID, clientID)
	return err
}

func (r *sqlRepository) getClientKeyID(ctx context.Context, userID int64, clientID string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `SELECT api_key_id FROM oauth_client_api_keys WHERE user_id = $1 AND client_id = $2`, userID, clientID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	return id, err
}

func (r *sqlRepository) setClientKeyID(ctx context.Context, userID int64, clientID string, apiKeyID int64) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO oauth_client_api_keys (user_id, client_id, api_key_id) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, client_id) DO UPDATE SET api_key_id = EXCLUDED.api_key_id, created_at = NOW()`,
		userID, clientID, apiKeyID)
	return err
}

func (r *sqlRepository) deleteExpired(ctx context.Context, before time.Time) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM oauth_access_tokens WHERE expires_at < $1`, before); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM oauth_authorization_codes WHERE expires_at < $1`, before)
	return err
}
