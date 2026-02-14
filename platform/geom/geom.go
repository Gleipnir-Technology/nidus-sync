package geom

import (
	"fmt"
)

func PostgisPointQuery(longitude, latitude float64) string {
	return fmt.Sprintf("ST_GeometryFromText('Point(%f %f)')", longitude, latitude)
}
