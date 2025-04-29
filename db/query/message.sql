
-- name: CreateThread :one
INSERT INTO thread (title)
VALUES ($1)
RETURNING *;

-- name: CreateMessage :one
INSERT INTO message (content,thread_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC;

-- name: DeleteMessageById :one
DELETE FROM message WHERE id = $1
RETURNING id;

-- name: DeleteMessageByThreadId :one
DELETE FROM message WHERE thread_id = $1
RETURNING thread_id;

-- name: UpdateMessage :exec
UPDATE message 
SET content = $2
WHERE id = $1
RETURNING *;

-- name: GetThreadById :one
SELECT * FROM thread WHERE id = $1;

-- name: GetMessageByThreadPaginated :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateProductCategory :one
INSERT INTO product_categories (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: CreateSingleProduct :one
INSERT INTO products (
    category_id,
    name,
    slug,
    description,
    type,
    regular_price,
    sale_price,
    sku,
    stock_id,
    main_image_url
)
VALUES ($1, $2, $3, $4, 'single', $5, $6, $7, $8, $9)
RETURNING id;

-- name: CreateVariableProduct :one
INSERT INTO products (
    category_id,
    name,
    slug,
    description,
    type,
    main_image_url
)
VALUES ($1, $2, $3, $4, 'variable', $5)
RETURNING id;

-- name: CreateStock :one
INSERT INTO stock (
    quantity,
    low_stock_threshold
)
VALUES ($1, $2)
RETURNING id;

-- name: CreateProductGallery :exec
INSERT INTO product_gallery (
    product_id,
    image_url
)
VALUES ($1, $2);

-- name: CreateProductVariation :one
INSERT INTO product_variations (
    product_id,
    variation_name,
    variation_value,
    regular_price,
    sale_price,
    sku,
    stock_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;
