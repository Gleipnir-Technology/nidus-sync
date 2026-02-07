package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/go-chi/chi/v5"
)

type ContentMock struct {
	District    ContentDistrict
	MapboxToken string
	ReportID    string
	URL         ContentURL
}

func addMockRoutes(r chi.Router) {
	r.Get("/", renderMock("rmo/mock/root.html"))
	r.Get("/district/{slug}", renderMock("rmo/mock/district-root.html"))
	r.Get("/district/{slug}/nuisance", renderMock("rmo/mock/nuisance.html"))
	r.Get("/district/{slug}/nuisance-submit-complete", renderMock("rmo/mock/nuisance-submit-complete.html"))
	r.Get("/district/{slug}/water", renderMock("rmo/mock/water.html"))
	r.Get("/nuisance", renderMock("rmo/mock/nuisance.html"))
	r.Get("/nuisance-submit-complete", renderMock("rmo/mock/nuisance-submit-complete.html"))
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
func renderMock(t string) func(http.ResponseWriter, *http.Request) {
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
