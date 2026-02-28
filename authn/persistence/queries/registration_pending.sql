-- name: GetRegistrationPending :one
SELECT id,
       email,
       code,
       code_expires_at,
       created_at
FROM registration_pending
WHERE id = @id::TEXT
LIMIT 1;


-- name: GetRegistrationPendingByEmail :one
SELECT *
FROM registration_pending
WHERE email = @email::TEXT
LIMIT 1;

-- name: CreateRegistrationPending :one
INSERT INTO registration_pending(email, code, code_expires_at)
VALUES (@email::TEXT, @code::TEXT, @code_expires_at::TIMESTAMP)
RETURNING *;

-- name: DeleteRegistrationPending :exec
DELETE FROM registration_pending
WHERE id = @id::TEXT;

-- name: DeleteRegistrationPendingByEmail :exec
DELETE FROM registration_pending
WHERE email = @email::TEXT;
