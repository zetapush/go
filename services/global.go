package zpservice

import ()

type Pagination struct {
	PageNumber int    `json:"pageNumber,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
	Direction  string `json:"direction,omitempty"`
}

type ImpersonatedRequest struct {
	Owner string `json:"owner,omitempty"`
}
