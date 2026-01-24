package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
)

var (
	mockRootT         = buildTemplate("mock/root", "base")
	mockDistrictRootT = buildTemplate("mock/district-root", "base")
)

type ContentDistrict struct {
	LogoURL string
	Name    string
}
type ContentMock struct {
	District ContentDistrict
}

func addMockRoutes(r chi.Router) {
	r.Get("/", renderMock(mockRootT))
	r.Get("/district/{slug}", renderMock(mockDistrictRootT))
}

func renderMock(t *htmlpage.BuiltTemplate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		htmlpage.RenderOrError(
			w,
			t,
			ContentMock{
				District: ContentDistrict{
					LogoURL: config.MakeURLNidus("/api/district/%s/logo", slug),
					Name:    "Delta MCD",
				},
			},
		)
	}
}
