

CREATE TABLE IF NOT EXISTS campaigns(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);



CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id),
    name VARCHAR(255) NOT NUll,
    description VARCHAR(255) NOT NULL DEFAULT '',
    priority SERIAL,/* max +1*/
    removed BOOLEAN DEFAULT false,
    created_at TIMESTAMP    NOT NULL default now()
);
create index ON items using btree(name)  ;
insert into campaigns(name) values ('Первая запись');

