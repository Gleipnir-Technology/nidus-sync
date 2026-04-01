package rmo

import (
	"net/http"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/html"
	"github.com/gorilla/mux"
)

type ContentMock struct {
	District ContentDistrict
	ReportID string
	URL      ContentURL
}

func addMockRoutes(r *mux.Router) {
	r.HandleFunc("/", renderMock("rmo/mock/root.html"))
	r.HandleFunc("/district/{slug}", renderMock("rmo/mock/district-root.html"))
	r.HandleFunc("/district/{slug}/nuisance-submit-complete", renderMock("rmo/mock/nuisance-submit-complete.html"))
	r.HandleFunc("/nuisance", renderMock("rmo/mock/nuisance.html"))
	r.HandleFunc("/nuisance-submit-complete", renderMock("rmo/mock/nuisance-submit-complete.html"))
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
		vars := mux.Vars(r)
		slug := vars["slug"]
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
				ReportID: "abcd-1234-5678",
				URL:      makeContentURLMock(slug),
			},
		)
	}
}
