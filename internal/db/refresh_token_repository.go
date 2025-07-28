package db

import (
	"auth-service/internal/model"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	SaveRefreshToken(record model.RefreshTokenRecord) error
	UpdateRefreshToken(userID uuid.UUID, newRecord model.RefreshTokenRecord) error
	GetRefreshTokensByUserID(userID uuid.UUID) ([]model.RefreshTokenRecord, error)
	RevokeTokensByUserID(userID uuid.UUID) error
	RevokeUser(userID uuid.UUID) error
	IsUserRevoked(userID uuid.UUID) (bool, error)
}

type refreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

// Реализация методов:
// Сохранение новой записи в БД
func (r *refreshTokenRepo) SaveRefreshToken(record model.RefreshTokenRecord) error {
	query := `INSERT INTO refresh_tokens
		(user_id, token_hash, issued_at, expires_at, revoked, used, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(context.Background(),
		query,
		record.UserID,
		record.TokenHash,
		record.IssuedAt,
		record.ExpiresAt,
		record.Revoked,
		record.Used,
		record.UserAgent,
		record.IP,
	)

	return err
}

// Получить запись из БД по UserID
func (r *refreshTokenRepo) GetRefreshTokensByUserID(userID uuid.UUID) ([]model.RefreshTokenRecord, error) {
	query := `SELECT user_id, token_hash, issued_at, expires_at, revoked, used, user_agent, ip 
              FROM refresh_tokens WHERE user_id = $1`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []model.RefreshTokenRecord
	for rows.Next() {
		var record model.RefreshTokenRecord
		err := rows.Scan(
			&record.UserID,
			&record.TokenHash,
			&record.IssuedAt,
			&record.ExpiresAt,
			&record.Revoked,
			&record.Used,
			&record.UserAgent,
			&record.IP,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, record)
	}

	return tokens, nil
}

// Обновить запись
func (r *refreshTokenRepo) UpdateRefreshToken(userID uuid.UUID, newRecord model.RefreshTokenRecord) error {
	query := `UPDATE refresh_tokens SET
		token_hash = $1,
		issued_at = $2,
		expires_at = $3,
		revoked = $4,
		used = $5,
		user_agent = $6,
		ip = $7
		WHERE user_id = $8`

	_, err := r.db.Exec(context.Background(),
		query,
		newRecord.TokenHash,
		newRecord.IssuedAt,
		newRecord.ExpiresAt,
		newRecord.Revoked,
		newRecord.Used,
		newRecord.UserAgent,
		newRecord.IP,
		userID,
	)

	return err
}

// Отозвать токен
func (r *refreshTokenRepo) RevokeTokensByUserID(userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	_, err := r.db.Exec(context.Background(), query, userID)

	return err
}

func (r *refreshTokenRepo) RevokeUser(userID uuid.UUID) error {
	query := `INSERT INTO revoked_users (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING`
	_, err := r.db.Exec(context.Background(), query, userID)

	return err
}

func (r *refreshTokenRepo) IsUserRevoked(userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM revoked_users WHERE user_id = $1)`
	var exists bool
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&exists)

	return exists, err
}
