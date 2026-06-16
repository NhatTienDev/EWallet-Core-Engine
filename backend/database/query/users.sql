-- name: CreateUser :one
INSERT INTO users (full_name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreatePasswordReset :one
INSERT INTO password_resets (user_id, hashed_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetValidPasswordReset :one
-- Get token that is not used and not expired
SELECT * FROM password_resets
WHERE hashed_token = $1 AND is_used = FALSE AND expires_at > NOW()
LIMIT 1;

-- name: MarkPasswordResetUsed :exec
UPDATE password_resets
SET is_used = TRUE
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET hashed_password = $2
WHERE id = $1;