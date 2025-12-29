-- name: CreateFeed :one

INSERT INTO feeds (id, name, created_at, updated_at, url, userid)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- -- name: GetFeedByUserId :one

-- Select * from feeds where userid = $1;