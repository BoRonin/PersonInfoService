package models

type HttpPostRequest struct {
	Person
}

type HttpFilterRequest struct {
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	OrderBy string `json:"order_by"`
	PersonFilter
}

type HttpSearchResponse struct {
	Page     int          `json:"page"`
	PerPage  int          `json:"per_page"`
	Quantity int          `json:"quantity"`
	LastPage int          `json:"last_page"`
	Persons  []PersonFull `json:"persons"`
}

func NewHttpFilterRequest() HttpFilterRequest {
	return HttpFilterRequest{
		Page:    1,
		PerPage: 10,
		OrderBy: "name_asc",
	}
}
