package geom

import (
	"fmt"
)

func PostgisPointQuery(longitude, latitude float64) string {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), 4326)", longitude, latitude)
}
