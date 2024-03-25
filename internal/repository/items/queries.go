package items

const (
	getListItem = `SELECT id, name, rating, category, image_url, reputation, price, availability, value,is_active
FROM items where is_active=true`

	getMyItems = `SELECT id, name, rating, category, image_url, reputation, price, availability,value,is_active
FROM items where created_by = $1`

	geItemByID = `SELECT id, name, rating, category, image_url, reputation, price, availability,value,is_active
FROM items where id = $1 and is_active = true`

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
    value = $8,
    is_active = $9
WHERE id = $10;
`

	saveHistoryChanges = `
INSERT INTO table_history (table_name, table_key, data_before, data_after, user_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6);
`

	getHistoryChanges = `
	SELECT id, table_name, table_key, data_before, data_after, user_id, created_at
FROM table_history WHERE table_name = $1 AND table_key = $2;
`
)
