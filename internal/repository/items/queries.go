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

	updateItembyID = `
UPDATE items
SET name = $1,
    rating = $2,
    category = $3,
    image_url = $4,
    reputation = $5,
    price = $6,
    availability = $7,
    value = $8
WHERE id = $9;
`

	saveHistoryChanges = `
INSERT INTO table_history (table_name, table_key, data_before, data_after, user_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6);


`
)
