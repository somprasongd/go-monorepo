CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE DATABASE "users"
  WITH OWNER "postgres"
  ENCODING 'UTF8';

CREATE DATABASE "todos"
  WITH OWNER "postgres"
  ENCODING 'UTF8';