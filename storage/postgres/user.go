package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/post/storage/repo"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{db: db}
}

func (ur *userRepo) Create(u *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users(
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			username,
			password,
			profile_image_url,
			type
		)values($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id,created_at
	`

	row := ur.db.QueryRow(
		query,
		u.FirstName,
		u.LastName,
		u.PhoneNumber,
		u.Email,
		u.Gender,
		u.UserName,
		u.Password,
		u.ProfileImageUrl,
		u.Type,
	)
	fmt.Println(u)

	if err := row.Scan(
		&u.Id,
		&u.CreatedAt,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *userRepo) Get(id int) (*repo.User, error) {
	var user repo.User

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			username,
			password,
			profile_image_url,
			type,
			created_at
		from users
		where id=$1
	`
	row := ur.db.QueryRow(query, id)
	if err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.Gender,
		&user.UserName,
		&user.Password,
		&user.ProfileImageUrl,
		&user.Type,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepo) GetAll(param repo.GetUserQuery) (*repo.GetAllUsersResult, error) {
	result := repo.GetAllUsersResult{
		Users: make([]*repo.User, 0),
	}

	offset := (param.Page - 1) * param.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", param.Limit, offset)
	filter := "WHERE true"
	if param.Search != "" {
		str := "%" + param.Search + "%"
		filter += fmt.Sprintf(` 
			and first_name ILIKE '%s' OR last_name ILIKE '%s' OR email ILIKE '%s' 
		OR username ILIKE '%s' OR phone_number ILIKE '%s'`, str, str, str, str, str)
	}

	if param.SortByDate == "" {
		param.SortByDate = "desc"
	}

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			username,
			password,
			profile_image_url,
			type,
			created_at
		FROM users
		` + filter + `
		ORDER BY created_at ` + param.SortByDate + ` ` + limit

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var usr repo.User
		if err := rows.Scan(
			&usr.Id,
			&usr.FirstName,
			&usr.LastName,
			&usr.PhoneNumber,
			&usr.Email,
			&usr.Gender,
			&usr.UserName,
			&usr.Password,
			&usr.ProfileImageUrl,
			&usr.Type,
			&usr.CreatedAt,
		); err != nil {
			return nil, err
		}
		result.Users = append(result.Users, &usr)
	}
	queryCount := `SELECT count(1) FROM users ` + filter
	err = ur.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

func (ur *userRepo) Update(usr *repo.User) (*repo.User, error) {
	query := `
		update users set 
			first_name=$1,
			last_name=$2,
			phone_number=$3,
			email=$4,
			gender=$5,
			username=$6,
			password=$7,
			profile_image_url=$8,
			type=$9
		where id=$10
		returning created_at
	`
	row := ur.db.QueryRow(
		query,
		usr.FirstName,
		usr.LastName,
		usr.PhoneNumber,
		usr.Email,
		usr.Gender,
		usr.UserName,
		usr.Password,
		usr.ProfileImageUrl,
		usr.Type,
		usr.Id,
	)

	if err := row.Scan(&usr.CreatedAt); err != nil {
		return nil, err
	}
	return usr, nil
}

func (ur *userRepo) Delete(id int) error {
	res, err := ur.db.Exec("delete from users where id=$1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (ur *userRepo) GetByEmail(email string) (*repo.User, error) {
	var result repo.User

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number,
			email,
			gender,
			password,
			username,
			profile_image_url,
			type,
			created_at
		FROM users
		WHERE email=$1
	`

	row := ur.db.QueryRow(query, email)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Email,
		&result.Gender,
		&result.Password,
		&result.UserName,
		&result.ProfileImageUrl,
		&result.Type,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (ur *userRepo) UpdatePassword(req *repo.UpdatePassword) error {
	query := `UPDATE users SET password=$1 WHERE id=$2`

	_, err := ur.db.Exec(query, req.Password, req.UserID)
	if err != nil {
		return err
	}

	return nil
}
