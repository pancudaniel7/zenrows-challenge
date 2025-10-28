INSERT INTO zenrows."user" (username, password_hash) VALUES
('alice', crypt('alicepass', gen_salt('bf'))),
('bob', crypt('bobpass', gen_salt('bf')));
