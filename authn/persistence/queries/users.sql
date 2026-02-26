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
       deleted_at
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
         deleted_at
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
