CREATE TABLE IF NOT EXISTS hotels (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    address text NOT NULL,
    location text NOT NULL
)