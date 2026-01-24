package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
)

var (
	mockDistrictRootT = buildTemplate("mock/district-root", "base")
	mockNuisanceT     = buildTemplate("mock/nuisance", "base")
	mockRootT         = buildTemplate("mock/root", "base")
)

type ContentDistrict struct {
	Name    string
	URLLogo string
}
type ContentURL struct {
	Nuisance string
}
type ContentMock struct {
	District    ContentDistrict
	MapboxToken string
	URL         ContentURL
}

func addMockRoutes(r chi.Router) {
	r.Get("/", renderMock(mockRootT))
	r.Get("/nuisance", renderMock(mockNuisanceT))
	r.Get("/district/{slug}", renderMock(mockDistrictRootT))
}

func makeContentURL() ContentURL {
	return ContentURL{
		Nuisance: makeURLMock("nuisance"),
	}
}

func makeURLMock(p string) string {
	return config.MakeURLReport("/mock/%s", p)
}
func renderMock(t *htmlpage.BuiltTemplate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		htmlpage.RenderOrError(
			w,
			t,
			ContentMock{
				District: ContentDistrict{
					Name:    "Delta MCD",
					URLLogo: config.MakeURLNidus("/api/district/%s/logo", slug),
				},
				MapboxToken: config.MapboxToken,
				URL:         makeContentURL(),
			},
		)
	}
}
