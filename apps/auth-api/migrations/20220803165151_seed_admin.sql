-- +goose Up
-- +goose StatementBegin
-- user admin password admin
INSERT INTO "users" ("email","password","role") VALUES ('admin@mail.com','$argon2id$v=19$m=65536,t=3,p=4$bScd+0u6Msk7aVcYHuMj6w$ui/+Ylw4geyMrsbZfJWI+vuSmjLGDu9Onkjvnonzj6M','admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "users" WHERE "email" = 'admin@mail.com';
-- +goose StatementEnd
