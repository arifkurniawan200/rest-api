package items

const (
	getListItem = `SELECT id, name, rating, category, image_url, reputation, price, availability, value
FROM items `

	getMyItems = `SELECT id, name, rating, category, image_url, reputation, price, availability,value
FROM items where created_by = $1`

	geItemByID = `SELECT id, name, rating, category, image_url, reputation, price, availability,value
FROM items where id = $1`

	insertNewItem = `INSERT INTO items (name, rating, category, image_url, reputation, price, availability, value)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
)
