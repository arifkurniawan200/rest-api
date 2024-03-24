package items

const (
	getListItem = `SELECT id, name, rating, category, image_url, reputation, price, availability
FROM items `

	getMyItems = `SELECT id, name, rating, category, image_url, reputation, price, availability
FROM items where created_by = $1`

	geItemByID = `SELECT id, name, rating, category, image_url, reputation, price, availability
FROM items where id = $1`
)
