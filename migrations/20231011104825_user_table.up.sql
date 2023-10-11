CREATE TABlE ref_user
(
    id         SERIAL PRIMARY KEY,
    login      TEXT UNIQUE,
    password   TEXT NOT NULL,
    created_on TIMESTAMP(0) WITH TIME ZONE DEFAULT now()
)