-- name: CreateVerifyEmail :one
INSERT INTO "verify_emails" (
    email,
    secret_code
 ) VALUES ($1,$2,$3) 
RETURNING *;

-- name: UpdateVerifyEmail :one 
UPDATE "verify_emails"
SET
  is_used = true
WHERE
id = $1
AND secret_code = $2
AND is_used = FALSE
AND expired_at > NOW()  
RETURNING *;    
