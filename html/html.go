package html

import (
	"net/http"
)

func RenderOrError(w http.ResponseWriter, template_name string, content interface{}) {
	templates.renderOrError(w, template_name, content)
}
