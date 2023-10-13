CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY NOT NULL,
    name character varying,
    phone character varying,
    zip character varying,
    city character varying,
    address character varying,
    region character varying,
    email character varying
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY NOT NULL,
    order_id integer,
    chrt_id integer,
    track_number character varying,
    price integer,
    rid character varying,
    name character varying,
    sale integer,
    size character varying,
    total_price integer,
    nm_id integer,
    brand character varying,
    status integer
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY NOT NULL,
    order_uid character varying,
    track_number character varying,
    entry character varying,
    locale character varying,
    internal_signature character varying,
    customer_id character varying,
    delivery_service character varying,
    shardkey character varying,
    sm_id integer,
    date_created character varying,
    oof_shard character varying
);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY NOT NULL,
    transaction character varying,
    request_id character varying,
    currency character varying,
    provider character varying,
    amount integer,
    payment_dt bigint,
    bank character varying,
    delivery_cost integer,
    goods_total integer,
    custom_fee integer
);
