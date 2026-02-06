package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

var (
	mockDistrictRootT           = buildTemplate("mock/district-root", "base")
	mockNuisanceT               = buildTemplate("mock/nuisance", "base")
	mockNuisanceSubmitCompleteT = buildTemplate("mock/nuisance-submit-complete", "base")
	mockRootT                   = buildTemplate("mock/root", "base")
	mockWaterT                  = buildTemplate("mock/water", "base")
)

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
	r.Get("/district/{slug}/water", renderMock(mockWaterT))
	r.Get("/nuisance", renderMock(mockNuisanceT))
	r.Get("/nuisance-submit-complete", renderMock(mockNuisanceSubmitCompleteT))
}

func makeContentURLMock(slug string) ContentURL {
	return ContentURL{
		Nuisance:       makeURLMock(slug, "nuisance"),
		SubmitComplete: makeURLMock(slug, "nuisance-submit-complete"),
		Tegola:         config.MakeURLTegola("/"),
		Water:          makeURLMock(slug, "water"),
	}
}
func makeURLMock(slug, p string) string {
	return config.MakeURLReport("/mock/district/%s/%s", slug, p)
}
func renderMock(t *html.BuiltTemplate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			slug = "delta-mvcd"
		}
		html.RenderOrError(
			w,
			t,
			ContentMock{
				District: ContentDistrict{
					Name:       "Delta MVCD",
					URLLogo:    config.MakeURLNidus("/api/district/%s/logo", slug),
					URLWebsite: "http://www.deltavcd.com/",
				},
				MapboxToken: config.MapboxToken,
				ReportID:    "abcd-1234-5678",
				URL:         makeContentURLMock(slug),
			},
		)
	}
}
