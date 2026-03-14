package geom

import (
	"fmt"

	"github.com/Gleipnir-Technology/nidus-sync/platform/types"
)

func PostgisPointQuery(location types.Location) string {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), 4326)", location.Longitude, location.Latitude)
}
