package repositories

import (
	"context"
	"database/sql"

	"github.com/ariiiiph/ecommerce/internal/db"
	"github.com/ariiiiph/ecommerce/internal/models"
)

type AddressRepository struct {
	db db.DBTX
}

func NewAddressRepository(db db.DBTX) *AddressRepository {
	return &AddressRepository{
		db: db,
	}
}

func (r *AddressRepository) Create(ctx context.Context, address *models.Address) error {
	query := `
		INSERT INTO addresses (
			user_id,
			title,
			recipient_name,
			phone,
			country,
			city,
			address_line,
			postal_code,
			is_default
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING
			id,
			user_id,
			title,
			recipient_name,
			phone,
			country,
			city,
			address_line,
			postal_code,
			is_default,
			created_at,
			updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		address.UserID,
		address.Title,
		address.RecipientName,
		address.Phone,
		address.Country,
		address.City,
		address.AddressLine,
		address.PostalCode,
		address.IsDefault,
	).Scan(
		&address.ID,
		&address.UserID,
		&address.Title,
		&address.RecipientName,
		&address.Phone,
		&address.Country,
		&address.City,
		&address.AddressLine,
		&address.PostalCode,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

}

func (r *AddressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	query := `SELECT id, user_id, title, recipient_name, phone, country, city, address_line, postal_code, is_default, created_at, updated_at FROM addresses WHERE id = $1`

	address := &models.Address{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&address.ID,
		&address.UserID,
		&address.Title,
		&address.RecipientName,
		&address.Phone,
		&address.Country,
		&address.City,
		&address.AddressLine,
		&address.PostalCode,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *AddressRepository) GetAllByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	query := `SELECT id, user_id, title, recipient_name, phone, country, city, address_line, postal_code, is_default, created_at, updated_at FROM addresses WHERE user_id = $1 ORDER BY is_default DESC, id ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var addresses []*models.Address

	for rows.Next() {
		address := &models.Address{}

		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.Title,
			&address.RecipientName,
			&address.Phone,
			&address.Country,
			&address.City,
			&address.AddressLine,
			&address.PostalCode,
			&address.IsDefault,
			&address.CreatedAt,
			&address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return addresses, nil

}

func (r *AddressRepository) Update(ctx context.Context, address *models.Address) error {
	query := `UPDATE addresses SET title = $1, recipient_name = $2, phone = $3, country = $4, city = $5, address_line = $6, postal_code = $7, is_default = $8, updated_at = NOW() WHERE id = $9`

	result, err := r.db.ExecContext(
		ctx,
		query,
		address.Title,
		address.RecipientName,
		address.Phone,
		address.Country,
		address.City,
		address.AddressLine,
		address.PostalCode,
		address.IsDefault,
		address.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (r *AddressRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM addresses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AddressRepository) CreateAsDefault(ctx context.Context, address *models.Address) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE addresses
		 SET is_default = FALSE, updated_at = NOW()
		 WHERE user_id = $1 AND is_default = TRUE`,
		address.UserID,
	)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO addresses (
			user_id,
			title,
			recipient_name,
			phone,
			country,
			city,
			address_line,
			postal_code,
			is_default
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE)
		RETURNING
			id,
			user_id,
			title,
			recipient_name,
			phone,
			country,
			city,
			address_line,
			postal_code,
			is_default,
			created_at,
			updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		address.UserID,
		address.Title,
		address.RecipientName,
		address.Phone,
		address.Country,
		address.City,
		address.AddressLine,
		address.PostalCode,
	).Scan(
		&address.ID,
		&address.UserID,
		&address.Title,
		&address.RecipientName,
		&address.Phone,
		&address.Country,
		&address.City,
		&address.AddressLine,
		&address.PostalCode,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
}

func (r *AddressRepository) UpdateAsDefault(ctx context.Context, address *models.Address) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE addresses
		 SET is_default = FALSE, updated_at = NOW()
		 WHERE user_id = $1 AND is_default = TRUE AND id != $2`,
		address.UserID,
		address.ID,
	)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(
		ctx,
		`UPDATE addresses
		 SET title = $1,
		     recipient_name = $2,
		     phone = $3,
		     country = $4,
		     city = $5,
		     address_line = $6,
		     postal_code = $7,
		     is_default = TRUE,
		     updated_at = NOW()
		 WHERE id = $8`,
		address.Title,
		address.RecipientName,
		address.Phone,
		address.Country,
		address.City,
		address.AddressLine,
		address.PostalCode,
		address.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AddressRepository) GetDefaultByUserID(ctx context.Context, userID int64) (*models.Address, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			recipient_name,
			phone,
			country,
			city,
			address_line,
			postal_code,
			is_default,
			created_at,
			updated_at
		FROM addresses
		WHERE user_id = $1
		  AND is_default = TRUE
		LIMIT 1
	`

	address := &models.Address{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&address.ID,
		&address.UserID,
		&address.Title,
		&address.RecipientName,
		&address.Phone,
		&address.Country,
		&address.City,
		&address.AddressLine,
		&address.PostalCode,
		&address.IsDefault,
		&address.CreatedAt,
		&address.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return address, nil
}
