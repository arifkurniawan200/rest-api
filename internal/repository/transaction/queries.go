package transaction

const insertPurchase = `INSERT INTO item_transactions (user_id, item_id, transaction_type, transaction_date, quantity, total_price, availability)
VALUES ($1, $2, $3, $4, $5, $6, $7);`
