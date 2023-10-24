CREATE TABlE user_lock
(
    user_id INT REFERENCES ref_user (id)  NOT NULL PRIMARY KEY
)