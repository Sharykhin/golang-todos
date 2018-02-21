CREATE TABLE todos (
  id integer PRIMARY KEY,
  title text NOT NULL,
  description text NOT NULL,
  completed integer,
  created DATETIME DEFAULT CURRENT_TIMESTAMP
)