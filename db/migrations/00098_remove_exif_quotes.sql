-- +goose Up
UPDATE publicreport.image_exif
SET value = SUBSTRING(value, 2, LENGTH(value) - 2)
WHERE value LIKE '"%"' 
  AND LENGTH(value) >= 2;
-- +goose Down
