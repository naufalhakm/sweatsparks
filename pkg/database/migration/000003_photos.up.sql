CREATE TABLE photos (
    id PRIMARY KEY,
    user_id INT NOT NULL,
    url VARCHAR(200) NOT NULL,
    is_primary SMALLINT NOT NULL,
    uploaded_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);  
