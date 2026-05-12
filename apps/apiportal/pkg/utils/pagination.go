package utils

type PaginationMetaData struct {
	TotalRecord int32 `json:"total_record"`
	TotalPages   int32 `json:"total_pages"`
	CurrentPage  int32 `json:"current_page"`
	PageSize     int32 `json:"page_size"`
}

type PaginationQueryParam struct {	
	Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
	Page   int32 `form:"page" binding:"omitempty,min=0"`
}


func (p *PaginationQueryParam) GetOffset() int32 {
	return (p.Page - 1) * p.Limit
}


func (p *PaginationQueryParam) GetMetadata(total_record  int32) PaginationMetaData {
	total_pages := total_record / p.Limit
	if total_record%p.Limit != 0 {
		total_pages += 1
	}

	return PaginationMetaData{
		TotalRecord: total_record,
		TotalPages:  total_pages,
		CurrentPage:  p.Page,
		PageSize:     p.Limit,
	}
}



