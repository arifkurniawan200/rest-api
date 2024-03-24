package user

import (
	"context"
	"database/sql"
	"template/internal/model"
	"template/internal/repository"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserHandler{db}
}

func (h UserHandler) BeginTx() (*sql.Tx, error) {
	return h.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h UserHandler) RegisterUser(c model.RequestRegisterUser) error {
	_, err := h.db.Exec(insertNewCostumer, c.FirstName, c.LastName, c.Email, c.Password)
	if err != nil {
		return err
	}
	return err
}

func (h UserHandler) GetUserByEmail(email string) (model.User, error) {
	var (
		data model.User
		err  error
	)
	rows, err := h.db.Query(getCostumerByEmail, email)
	if err != nil {
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&data.ID, &data.FirstName, &data.LastName, &data.Email, &data.Password, &data.Type,
			&data.CreatedAt, &data.UpdatedAt, &data.DeletedAt,
		); err != nil {
			return data, err
		}
	}

	if err = rows.Err(); err != nil {
		return data, err
	}
	return data, err
}
