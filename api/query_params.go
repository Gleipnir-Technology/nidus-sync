package api

type queryParams struct {
	Limit *int    `schema:"limit"`
	Query *string `schema:"query"`
	Sort  *string `schema:"sort"`
	Type  *string `schema:"type"`
}

func (qp queryParams) SortOrDefault(default_name string, ascending bool) (string, bool) {
	if qp.Sort == nil {
		return default_name, ascending
	}
	s := *qp.Sort
	if s == "" {
		return default_name, ascending
	}
	a := true
	if s[0] == '-' {
		a = false
	}
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	return s, a
}
