package customer

type CustomersQueryOptions struct {
	Page       int
	Limit      int
	OrderByCol string
	Order      string
	FirstName  string
	LastName   string

	ResultPages int
	ResultPage  int
	ResultRecs  int
}

func (qo *CustomersQueryOptions) SetPaginationAttrs(cnt int) {
	pages := 1
	if cnt > 0 && qo.Limit > 0 {
		pages = int(cnt) / qo.Limit
		if int(cnt)%qo.Limit > 0 {
			pages += 1
		}
	}

	if qo.Page > pages {
		qo.Page = pages
	}

	qo.ResultPages = pages
	qo.ResultPage = qo.Page
	qo.ResultRecs = int(cnt)
}
