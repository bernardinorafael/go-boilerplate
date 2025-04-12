package dto

type SearchParams struct {
	Term  string `json:"term,omitempty"`
	Sort  string `json:"sort,omitempty"`
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
}
