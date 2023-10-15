CREATE TABlE doc_order
(
    id         SERIAL PRIMARY KEY,
    number     BIGINT UNIQUE,
    user_id    INT REFERENCES ref_user (id) NOT NULL,
    status     VARCHAR(100)                 NOT NULL DEFAULT 'NEW',
    changed_at TIMESTAMP(0) WITH TIME ZONE           DEFAULT now(),
    created_at TIMESTAMP(0) WITH TIME ZONE           DEFAULT now()
)