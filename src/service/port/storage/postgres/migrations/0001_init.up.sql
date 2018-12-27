CREATE TABLE port (
	id char(5) PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	city TEXT NOT NULL DEFAULT '',
	country TEXT NOT NULL DEFAULT '',
	alias TEXT[] NOT NULL,
	regions TEXT[] NOT NULL,
	coordinates TEXT[] NOT NULL,
	province TEXT NOT NULL DEFAULT '',
	timezone TEXT NOT NULL DEFAULT '',
	unlocs TEXT[] NOT NULL,
	code char(5) NOT NULL,

	created_at timestamp with time zone NOT NULL DEFAULT now(),
	updated_at timestamp with time zone NOT NULL DEFAULT now()
);
