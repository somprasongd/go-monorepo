-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE todo_status AS ENUM ('open', 'done');

CREATE TABLE public.todos (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	"text" text NOT NULL,
	status todo_status NULL DEFAULT 'open'::todo_status,
	user_id uuid NOT NULL,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	CONSTRAINT todos_pkey PRIMARY KEY (id)
);

CREATE INDEX todos_status ON public.todos (status);
CREATE INDEX todos_user ON public.todos (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todos;
DROP TYPE todo_status;
-- +goose StatementEnd
