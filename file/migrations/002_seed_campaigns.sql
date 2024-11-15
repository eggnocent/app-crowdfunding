-- +migrate Up
INSERT INTO users (name, occupation, email, password_hash, avatar_file_name, role, token, created_at, updated_at)
VALUES
('User One', 'Engineer', 'userone@example.com', 'hashedpassword1', 'avatar1.jpg', 'admin', 'token1', NOW(), NOW()),
('User Two', 'Designer', 'usertwo@example.com', 'hashedpassword2', 'avatar2.jpg', 'user', 'token2', NOW(), NOW());

-- +migrate Down
DELETE FROM users WHERE email IN ('userone@example.com', 'usertwo@example.com');

