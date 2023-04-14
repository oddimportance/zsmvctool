package persistence

import ()

type PaginationDetails struct {
	RowsInTotal             int
	PagesTotal              int
	PagePresent             int
	PageNext                int
	PagePrevious            int
	ItemsPerPage            int
	ItemFrom                int
	ItemUntil               int
	SearchParamsAsUrlString string
	SearchParams            map[string]string
}
