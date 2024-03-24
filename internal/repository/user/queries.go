package user

const (
	insertNewCostumer  = `INSERT INTO users(first_name,last_name,email, password) VALUES ($1, $2, $3, $4)`
	getCostumerByEmail = `SELECT id, first_name, last_name, email, password, type, created_at, updated_at, deleted_at
			  FROM users
			  WHERE email = $1`
)
