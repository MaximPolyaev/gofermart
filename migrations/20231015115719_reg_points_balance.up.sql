CREATE TABlE reg_points_balance
(
    id         SERIAL PRIMARY KEY,
    order_id   INT REFERENCES doc_order (id) NOT NULL,
    points     NUMERIC(9, 2)                 NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT now()
)