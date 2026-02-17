package html

import (
	"net/http"
	"time"
)

func RenderOrError(w http.ResponseWriter, template_name string, content interface{}) {
	templates.renderOrError(w, template_name, content)
}

var startedTime time.Time

func SetStartedTime() {
	startedTime = time.Now()
}
