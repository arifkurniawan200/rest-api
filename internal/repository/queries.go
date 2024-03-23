package repository

const (
	insertNewCostumer  = `INSERT INTO users(NIK, full_name, born_place, born_date,is_admin, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)`
	getCostumerByEmail = `
        SELECT
            id,
            NIK,
            full_name,
            born_place,
            born_date,
            is_admin,
            email,
            password,
            created_at,
            updated_at,
            deleted_at
        FROM
            users
        WHERE
            email = ?
    `
)
