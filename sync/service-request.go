package sync

import (
	"context"
	"time"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
)

type contentActiveServiceRequest struct {
	Created    time.Time
	LastAction time.Time
	NextStep   string
	Address    string
	PhotoCount uint
	Type       string
	URLDetail  string
}
type contentClosedServiceRequest struct {
	Employee         string
	Type             string
	Closed           time.Time
	Address          string
	TimeToResolution time.Duration
	URLDetail        string
}
type contentServiceRequestDetail struct{}
type contentServiceRequestList struct {
	ActiveRequests []contentActiveServiceRequest
	ClosedRequests []contentClosedServiceRequest
}

func getServiceRequestDetail(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	content := contentServiceRequestDetail{}
	return "sync/service-request-detail.html", content, nil
}
func getServiceRequestList(ctx context.Context, user *models.User) (string, interface{}, *errorWithStatus) {
	now := time.Now()
	content := contentServiceRequestList{
		ActiveRequests: []contentActiveServiceRequest{
			contentActiveServiceRequest{
				Created:    now.Add(-2 * time.Hour),
				LastAction: now.Add(-2 * time.Hour),
				NextStep:   "schedule-appointment",
				Address:    "123 Main St, Anytown",
				PhotoCount: 3,
				Type:       "biting-nuisance",
				URLDetail:  config.MakeURLNidus("/service-request/1"),
			},
			contentActiveServiceRequest{
				Created:    now.Add(-5 * time.Hour),
				LastAction: now.Add(-1 * time.Hour),
				NextStep:   "answer-question",
				Address:    "456 Elm St, Anytown",
				PhotoCount: 1,
				Type:       "standing-water",
				URLDetail:  config.MakeURLNidus("/service-request/1"),
			},
			contentActiveServiceRequest{
				Created:    now.Add(-1 * 24 * time.Hour),
				LastAction: now.Add(-3 * time.Hour),
				NextStep:   "add-to-route",
				Address:    "789 Oak St, Anytown",
				PhotoCount: 4,
				Type:       "active-breeding",
				URLDetail:  config.MakeURLNidus("/service-request/1"),
			},
			contentActiveServiceRequest{
				Created:    now.Add(-2 * 24 * time.Hour),
				LastAction: now.Add(-6 * time.Hour),
				NextStep:   "review",
				Address:    "101 Pine Ln, Anytown",
				PhotoCount: 0,
				Type:       "standing-water",
				URLDetail:  config.MakeURLNidus("/service-request/1"),
			},
		},
		ClosedRequests: []contentClosedServiceRequest{
			contentClosedServiceRequest{
				Employee:         "John Smith",
				Type:             "standing-water",
				Closed:           now.Add(-1 * 24 * time.Hour),
				Address:          "303 Ceder St, Anytown",
				TimeToResolution: 3 * 24 * time.Hour,
				URLDetail:        config.MakeURLNidus("/service-request/2"),
			},
			contentClosedServiceRequest{
				Employee:         "Maria Garcia",
				Type:             "biting-nuisance",
				Closed:           now.Add(-2 * 24 * time.Hour),
				Address:          "404 Birch St, Anytown",
				TimeToResolution: 1 * 24 * time.Hour,
				URLDetail:        config.MakeURLNidus("/service-request/2"),
			},
			contentClosedServiceRequest{
				Employee:         "Robert Johnson",
				Type:             "active-breeding",
				Closed:           now.Add(-4 * 24 * time.Hour),
				Address:          "404 Birch St, Anytown",
				TimeToResolution: 5 * 24 * time.Hour,
				URLDetail:        config.MakeURLNidus("/service-request/2"),
			},
			contentClosedServiceRequest{
				Employee:         "Sarah Lee",
				Type:             "standing-water",
				Closed:           now.Add(-7 * 24 * time.Hour),
				Address:          "606 Willow Way, Anytown",
				TimeToResolution: 2 * 24 * time.Hour,
				URLDetail:        config.MakeURLNidus("/service-request/2"),
			},
		},
	}
	return "sync/service-request-list.html", content, nil
}
