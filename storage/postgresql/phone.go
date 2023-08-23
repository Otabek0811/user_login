package postgresql

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type phoneRepo struct {
	db *pgxpool.Pool
}

func NewPhoneRepo(db *pgxpool.Pool) *phoneRepo {
	return &phoneRepo{
		db: db,
	}
}

func (r *phoneRepo) Create(ctx context.Context, req *models.CreatePhone) (string, error) {
	var (
		query string
		id    string
	)
	id = uuid.NewString()

	query = `
		INSERT INTO phones(
			id, 
			user_id,
			phone,
			description,
			is_fax,
			updated_at 
		)
		VALUES ( $1, $2, $3, $4, $5, now())
	`
	_, err := r.db.Exec(ctx, query,
		id,
		req.UserID,
		req.Phone,
		req.Description,
		req.IsFax,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *phoneRepo) GetByID(ctx context.Context, req *models.PhonePrimaryKey) (*models.Phone, error) {

	var (
		query string
		phone  models.Phone
	)

	query = `
		SELECT
			id, 
			user_id,
			phone,
			description,
			is_fax,
			CAST(created_at::timestamp AS VARCHAR),
			CAST(updated_at::timestamp AS VARCHAR)
		FROM phones
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&phone.Id,
		&phone.UserID,
		&phone.Phone,
		&phone.Description,
		&phone.IsFax,
		&phone.CreatedAt,
		&phone.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &phone, nil
}

func (r *phoneRepo) GetList(ctx context.Context, req *models.GetListPhoneRequest) (resp *models.GetListPhoneResponse, err error) {

	resp = &models.GetListPhoneResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id, 
			user_id,
			phone,
			description,
			is_fax,
			CAST(created_at::timestamp AS VARCHAR),
			CAST(updated_at::timestamp AS VARCHAR)
		FROM phones
	`

	if len(req.Search) > 0 {
		filter += " AND phone ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	if len(req.UserID)>0{
		filter+= " AND user_id =  '" + req.UserID + "'"
	}

	query += filter + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var phone models.Phone
		err = rows.Scan(
			&resp.Count,
			&phone.Id,
			&phone.UserID,
			&phone.Phone,
			&phone.Description,
			&phone.IsFax,
			&phone.CreatedAt,
			&phone.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Phones = append(resp.Phones, &phone)
	}

	return resp, nil
}

func (r *phoneRepo) Update(ctx context.Context, req *models.UpdatePhone) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
		phones
		SET
			id = :id,
			user_id = :user_id,
			phone = :phone,
			description = :description,
			is_fax = :is_fax,
			updated_at = now()
		WHERE id = :id and user_id = :user_id
	`

	params = map[string]interface{}{
		"id":       req.Id,
		"user_id":     req.UserID,
		"phone":    req.Phone,
		"description": req.Description,
		"is_fax":      req.IsFax,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *phoneRepo) Delete(ctx context.Context, req *models.PhonePrimaryKey) (int64, error) {
	query := `
		DELETE 
		FROM phones
		WHERE id = $1 
	`

	result, err := r.db.Exec(ctx, query, req.Id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
