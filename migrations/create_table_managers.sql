CREATE TABLE customers(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customer_tokens(
    token TEXT NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL REFERENCES customers,
    expire TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP+INTERVAL '1 hour',
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

create table managers (
    id bigserial primary key,
    name text not null,
    login text not null unique,
    password text not null,
    salary integer not null check(salary >0 ),
    plan integer not null check(plan >0 ),
    boss_id bigint not null,
    department text,
    active boolean not null default true,
    created timestamp not null default current_timestamp
);