-- name: GetUserConfirmationByUserID :one
SELECT *
  FROM user_confirmations
 WHERE user_id = @userID::TEXT
 LIMIT 1;

-- name: GetUserConfirmationByToken :one
SELECT *
  FROM user_confirmations
 WHERE token = @token::TEXT
   AND expires_at > now()
 LIMIT 1;

-- name: CreateUserConfirmation :one
INSERT INTO user_confirmations(user_id, token, expires_at)
     VALUES (@userID::TEXT, @token::TEXT, @expiresAt::TIMESTAMP)
  RETURNING *;

-- name: ConfirmUserRegistration :exec
UPDATE user_confirmations
   SET confirmed_at = now(),
       updated_at = now()
 WHERE user_id = @userID::TEXT
   AND confirmed_at IS NULL;

-- name: DeleteUserConfirmation :exec
DELETE FROM user_confirmations
 WHERE user_id = @userID::TEXT;

-- name: DeleteExpiredUserConfirmations :exec
DELETE FROM user_confirmations
 WHERE expires_at < now();
