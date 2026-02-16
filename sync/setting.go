package sync

import (
	"context"
	"net/http"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/sm"
	"github.com/Gleipnir-Technology/nidus-sync/arcgis"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	//"github.com/rs/zerolog/log"
	"github.com/stephenafamo/scan"
)

type contentDistrict struct {
	Address                 string  `db:"address"`
	Agency                  string  `db:"agency"`
	Centroid                string  `db:"st_asgeojson"`
	Contact                 string  `db:"contact"`
	City1                   string  `db:"city1"`
	City2                   string  `db:"city2"`
	Fax                     string  `db:"fax1"`
	GID                     int32   `db:"gid"`
	Phone1                  string  `db:"phone1"`
	Phone2                  string  `db:"phone2"`
	GeneralManager          string  `db:"general_mg"`
	PostalCode              string  `db:"postal_c_1"`
	ShapeArea               string  `db:"shape_area"`
	SurfaceAreaMetersSquare string  `db:"area_4326_sqm"`
	Website                 string  `db:"website"`
	XMin                    float32 `db:"st_xmin"`
	XMax                    float32 `db:"st_xmax"`
	YMin                    float32 `db:"st_ymin"`
	YMax                    float32 `db:"st_ymax"`
}
type contentSettingDistrict struct {
	District contentDistrict
	URL      ContentURL
	User     User
}

type contentSettingIntegration struct {
	ArcGISOAuth *models.OauthToken
	URL         ContentURL
	User        User
}

func getSetting(w http.ResponseWriter, r *http.Request, u *models.User) {
	userContent, err := contentForUser(r.Context(), u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	data := ContentAuthenticatedPlaceholder{
		URL:  newContentURL(),
		User: userContent,
	}
	html.RenderOrError(w, "sync/settings.html", data)
}
func getSettingDistrict(w http.ResponseWriter, r *http.Request, u *models.User) {
	ctx := r.Context()
	userContent, err := contentForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	org, err := u.Organization().One(ctx, db.PGInstance.BobDB)
	var district contentDistrict
	gid := int32(0)
	if org.ImportDistrictGid.IsValue() {
		gid = org.ImportDistrictGid.MustGet()
		district, err = bob.One[contentDistrict](ctx, db.PGInstance.BobDB, psql.Select(
			sm.From("import.district"),
			sm.Columns(
				"address",
				"agency",
				"area_4326_sqm",
				"city1",
				"city2",
				"contact",
				"fax1",
				"general_mg",
				"gid",
				"phone1",
				"phone2",
				"postal_c_1",
				"website",
				psql.F("ST_AsGeoJSON", "centroid_4326"),
				psql.F("ST_XMin", "extent_4326"),
				psql.F("ST_YMin", "extent_4326"),
				psql.F("ST_XMax", "extent_4326"),
				psql.F("ST_YMax", "extent_4326"),
			),
			sm.Where(psql.Quote("gid").EQ(psql.Arg(gid))),
		), scan.StructMapper[contentDistrict]())
		if err != nil {
			respondError(w, "Failed to get extents", err, http.StatusInternalServerError)
			return
		}
	}
	data := contentSettingDistrict{
		District: district,
		URL:      newContentURL(),
		User:     userContent,
	}
	html.RenderOrError(w, "sync/setting-district.html", data)
}
func getSettingIntegration(w http.ResponseWriter, r *http.Request, u *models.User) {
	ctx := r.Context()
	userContent, err := contentForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get user content", err, http.StatusInternalServerError)
		return
	}
	oauth, err := arcgis.GetOAuthForUser(ctx, u)
	if err != nil {
		respondError(w, "Failed to get oauth", err, http.StatusInternalServerError)
		return
	}
	data := contentSettingIntegration{
		ArcGISOAuth: oauth,
		URL:         newContentURL(),
		User:        userContent,
	}
	html.RenderOrError(w, "sync/setting-integration.html", data)
}

type contentSettingPesticide struct{}

func getSettingPesticide(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentSettingPesticide{}
	return "sync/setting-pesticide.html", content, nil
}
func getSettingPesticideAdd(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentSettingPesticide{}
	return "sync/setting-pesticide-add.html", content, nil
}
