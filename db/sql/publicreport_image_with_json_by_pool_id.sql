-- PublicreportImageWithJSONByPoolID
SELECT 
	"publicreport.image"."id" AS "id",
	"publicreport.image"."content_type" AS "content_type",
	"publicreport.image"."created" AS "created",
	"publicreport.image"."location" AS "location",
	COALESCE(ST_AsGeoJSON("publicreport.image"."location"), '{}') AS "location_json",
	"publicreport.image"."resolution_x" AS "resolution_x",
	"publicreport.image"."resolution_y" AS "resolution_y",
	"publicreport.image"."storage_uuid" AS "storage_uuid",
	"publicreport.image"."storage_size" AS "storage_size",
	"publicreport.image"."uploaded_filename" AS "uploaded_filename"
FROM "publicreport"."image" AS "publicreport.image"
INNER JOIN "publicreport"."pool_image" AS "publicreport.pool_image" ON ("publicreport.image"."id" = "publicreport.pool_image"."image_id")
WHERE ("publicreport.pool_image"."pool_id" = $1)
