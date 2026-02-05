-- PublicreportImageWithJSONByNuisanceID
SELECT 
	"publicreport.image"."id" AS "id",
	"publicreport.image"."content_type" AS "content_type",
	"publicreport.image"."created" AS "created",
	"publicreport.image"."location" AS "location",
	ST_AsGeoJSON("publicreport.image"."location") AS "location_json",
	"publicreport.image"."resolution_x" AS "resolution_x",
	"publicreport.image"."resolution_y" AS "resolution_y",
	"publicreport.image"."storage_uuid" AS "storage_uuid",
	"publicreport.image"."storage_size" AS "storage_size",
	"publicreport.image"."uploaded_filename" AS "uploaded_filename"
FROM "publicreport"."image" AS "publicreport.image"
INNER JOIN "publicreport"."nuisance_image" AS "publicreport.nuisance_image" ON ("publicreport.image"."id" = "publicreport.nuisance_image"."image_id")
WHERE ("publicreport.nuisance_image"."nuisance_id" = $1)
