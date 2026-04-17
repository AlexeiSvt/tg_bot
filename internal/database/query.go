package database

const (
DB_Create = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT UNIQUE,
    first_name TEXT,
    last_name TEXT,
    class TEXT,
    disciplines JSONB
);`

Insert_Query = `INSERT INTO users (tg_id, first_name, last_name, class, disciplines)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tg_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			class = EXCLUDED.class,
			disciplines = EXCLUDED.disciplines
		RETURNING id`
)