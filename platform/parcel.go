package platform

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/paulmach/orb/geojson"
	"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

func ParcelEnvelope(ctx context.Context, parcel_id int32) (*geojson.Polygon, error) {
	type _Row struct {
		Envelope string `db:"st_asgeojson"`
	}
	row, err := bob.One(ctx, db.PGInstance.BobDB, psql.Select(
		sm.Columns(
			psql.F("ST_AsGeoJSON", psql.F("ST_Envelope", psql.Raw("geometry")))(),
		),
		sm.From("parcel"),
		sm.Where(psql.Quote("id").EQ(psql.Arg(parcel_id))),
	), scan.StructMapper[_Row]())
	if err != nil {
		return nil, fmt.Errorf("query parcel: %w", err)
	}
	var polygon geojson.Polygon
	log.Info().Str("envelope", row.Envelope).Msg("about to unmarshal")
	err = json.Unmarshal([]byte(row.Envelope), &polygon)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return &polygon, nil
}
