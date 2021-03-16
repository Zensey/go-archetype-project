package customer

type CustomersQueryOptions struct {
	OrderByCol string
	Order      string
	FirstName  string
	LastName   string
}

type PaginationOptions struct {
	Page  int
	Limit int

	ResultPages int
	ResultPage  int
	ResultRecs  int
}

func (o *PaginationOptions) SetPaginationAttrs(cnt int) {
	pages := 1
	if cnt > 0 && o.Limit > 0 {
		pages = cnt / o.Limit
		if cnt%o.Limit > 0 {
			pages += 1
		}
	}

	if o.Page > pages {
		o.Page = pages
	}

	o.ResultPages = pages
	o.ResultPage = o.Page
	o.ResultRecs = cnt
}
