package dtos

type Pagination struct {
	PageNo   int64 `json:"pageNo"`
	PageSize int64 `json:"pageSize"`
}

func (p *Pagination) GetSkipAndLimit() (int64, int64) {
	skip := p.PageNo*p.PageSize - p.PageSize
	limit := p.PageSize
	return skip, limit
}

type QueryFilter struct {
	Filter map[string]interface{} `json:"filter,omitempty"`
}

type QuerySort struct {
	Sort map[string]*bool `json:"sort,omitempty"`
}

type IdRequestDTO struct {
	Id string `json:"id"`
}

type IdArrayRequestDTO struct {
	Ids []string `json:"ids"`
}

type ResponseBaseDTO struct {
	Count int16 `json:"_count"`
}

type QueryRequestDTO struct {
	QueryFilter
	QuerySort
}

type ResponseWithCount[T any] struct {
	Count int64 `bson:"totalCount" json:"totalCount"`
	Data  []T   `bson:"data" json:"data"`
}

type DeleteQueryParams struct {
	HardDelete  bool `json:"hard"`
	ForceDelete bool `json:"force"`
	Cascade     bool `json:"cascade"`
}
