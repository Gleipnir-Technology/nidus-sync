package publicreport

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/htmlpage"
	"github.com/go-chi/chi/v5"
)

var (
	mockDistrictRootT           = buildTemplate("mock/district-root", "base")
	mockNuisanceT               = buildTemplate("mock/nuisance", "base")
	mockNuisanceSubmitCompleteT = buildTemplate("mock/nuisance-submit-complete", "base")
	mockRootT                   = buildTemplate("mock/root", "base")
	mockStatusT                 = buildTemplate("mock/status", "base")
)

type ContentDistrict struct {
	Name    string
	URLLogo string
}
type ContentURL struct {
	Nuisance               string
	NuisanceSubmitComplete string
	Status                 string
	Tegola                 string
}
type ContentMock struct {
	District    ContentDistrict
	MapboxToken string
	ReportID    string
	URL         ContentURL
}

func addMockRoutes(r chi.Router) {
	r.Get("/", renderMock(mockRootT))
	r.Get("/district/{slug}", renderMock(mockDistrictRootT))
	r.Get("/district/{slug}/nuisance", renderMock(mockNuisanceT))
	r.Get("/district/{slug}/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT))
	r.Get("/district/{slug}/status", renderMock(mockStatusT))
	r.Get("/nuisance", renderMock(mockNuisanceT))
	r.Get("/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT))
	r.Get("/status", renderMock(mockStatusT))
}

func makeContentURL(slug string) ContentURL {
	return ContentURL{
		Nuisance:               makeURLMock(slug, "nuisance"),
		NuisanceSubmitComplete: makeURLMock(slug, "nuisance-submit-complete"),
		Status:                 makeURLMock(slug, "status"),
		Tegola:                 config.MakeURLTegola("/"),
	}
}

func makeURLMock(slug, p string) string {
	return config.MakeURLReport("/mock/district/%s/%s", slug, p)
}
func renderMock(t *htmlpage.BuiltTemplate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			slug = "delta-mvcd"
		}
		htmlpage.RenderOrError(
			w,
			t,
			ContentMock{
				District: ContentDistrict{
					Name:    "Delta MVCD",
					URLLogo: config.MakeURLNidus("/api/district/%s/logo", slug),
				},
				MapboxToken: config.MapboxToken,
				ReportID:    "abcd-1234-5678",
				URL:         makeContentURL(slug),
			},
		)
	}
}
