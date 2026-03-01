-- name: GetUser :one
SELECT id,
       email,
       password_hash,
       first_name,
       last_name,
       locale,
       timezone,
       created_at,
       updated_at,
       deleted_at,
       password_reset_token,
       password_reset_token_expires_at
FROM users
WHERE id = @id::TEXT
LIMIT 1;


-- name: GetUserByEmail :one
SELECT *
  FROM users
 WHERE email = @email::TEXT
   AND deleted_at is NULL
 LIMIT 1;

-- name: ListUsers :many
  SELECT id,
         email,
         password_hash,
         first_name,
         last_name,
         locale,
         timezone,
         created_at,
         updated_at,
         deleted_at,
         password_reset_token,
         password_reset_token_expires_at
    FROM users
   WHERE (sqlc.narg('email')::TEXT IS NULL OR email = sqlc.narg('email'))
     AND (sqlc.narg('created_start_range')::TIMESTAMP IS NULL OR created_at >= sqlc.narg('created_start_range')::TIMESTAMP)
     AND (sqlc.narg('created_end_range')::TIMESTAMP IS NULL OR created_at <= sqlc.narg('created_end_range')::TIMESTAMP)
     AND (deleted_at is null) = @active::BOOLEAN;

-- name: CreateUser :one
INSERT INTO users(email, password_hash)
     VALUES (@email::TEXT, @password_hash::BYTEA)
  RETURNING *;

-- name: UpdateUser :one
    UPDATE users
       SET email         = @email::TEXT,
           password_hash = @password_hash::BYTEA,
           first_name    = @first_name::TEXT,
           last_name     = @last_name::TEXT,
           locale        = @locale::TEXT,
           timezone      = @timezone::TEXT,
           updated_at    = now()
     WHERE id = @id::TEXT
 RETURNING *;


-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE id = @id::TEXT;

-- name: SetPasswordResetToken :exec
UPDATE users
SET password_reset_token = @token::TEXT,
    password_reset_token_expires_at = @expiresAt::TIMESTAMP,
    updated_at = now()
WHERE id = @userID::TEXT;

-- name: GetUserByPasswordResetToken :one
SELECT *
  FROM users
 WHERE password_reset_token = @token::TEXT
   AND password_reset_token_expires_at > now()
 LIMIT 1;

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = @passwordHash::BYTEA,
    password_reset_token = NULL,
    password_reset_token_expires_at = NULL,
    updated_at = now()
WHERE id = @userID::TEXT;

-- name: CreateUserConfirmation :one
INSERT INTO user_confirmations(user_id, token, expires_at)
     VALUES (@userID::TEXT, @token::TEXT, @expiresAt::TIMESTAMP)
  RETURNING *;

-- name: GetUserConfirmationByToken :one
SELECT id, user_id, token, expires_at, confirmed_at, created_at, updated_at
  FROM user_confirmations
 WHERE token = @token::TEXT
   AND expires_at > now()
   AND confirmed_at IS NULL
 LIMIT 1;

-- name: ConfirmUserRegistration :exec
UPDATE user_confirmations
SET confirmed_at = now(),
    updated_at = now()
WHERE user_id = @userID::TEXT;

-- name: DeleteUserConfirmation :exec
DELETE FROM user_confirmations WHERE user_id = @userID::TEXT;
