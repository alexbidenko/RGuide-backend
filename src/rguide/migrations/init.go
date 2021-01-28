package migrations

var Schema = `
CREATE TABLE products (
    id int,
    title text,
    description text
);`
