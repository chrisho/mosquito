package paginate

import (
	"github.com/chrisho/mosquito/helper"
)

var (
	DefaultPageSize int32 = 10 // 默认显示条数
	PaginateAction        = map[int32]string{0: "PreviousPage", 1: "NextPage", }
	PaginateMode          = map[int32]string{0: "PageNumber", 1: "MaxId"}
)

// 设置 默认每页 10 条，页码 第 1 页
func SetDefaultPaginateOptions(in *PageOptions) *PageOptions {

	// set default pageSize ：
	if in.PageSize < 1 {
		in.PageSize = DefaultPageSize
	}

	// set default first page : 1
	if in.PageNumber < 1 {
		in.PageNumber = 1
	}

	return in
}

// 获取 select 参数
func GetSelectOptions(
	in *PageOptions, whereStr string, whereParams []interface{},
) (
	where string, params []interface{}, orderBy, offset interface{},
) {

	switch PaginateMode[in.PaginateMode] {
	case "MaxId":
		whereStr, whereParams, orderBy, offset = GetMaxIdSelectOptions(in, whereStr, whereParams)
	case "PageNumber":
		orderBy, offset = GetPageNumberSelectOptions(in)
	}

	return whereStr, whereParams, orderBy, offset
}

// PaginateMode ： PageNumber
func GetPageNumberSelectOptions(in *PageOptions) (orderBy, offset interface{}) {

	offset = in.PageSize * (in.PageNumber - 1)

	if in.OrderBy == "" {
		orderBy = "id desc"
	}

	return
}

// Action : 1 == NextPage , 0 == PreviousPage
// PaginateMode ： MaxId
// Action default NextPage
// default sort : id desc
func GetMaxIdSelectOptions(
	in *PageOptions, whereStr string, whereParams []interface{},
) (
	where string, params []interface{}, orderBy, offset interface{},
) {

	switch PaginateAction[in.Action] {
	case "NextPage":
		if in.MaxId > 0 {
			whereStr += " and id <= ? "
			whereParams = append(whereParams, in.MaxId)
			offset = in.PageSize
		} else {
			offset = 0
		}
		orderBy = "id desc"
	case "PreviousPage":
		whereStr += " and id > ? "
		whereParams = append(whereParams, in.MaxId)
		orderBy = "id asc"
		offset = 0
	default:
		orderBy = "id desc"
		offset = 0
	}

	return whereStr, whereParams, orderBy, offset
}

// set paginate data
// total TotalRecords entries, total TotalPages pages, MaxId , MinId and PageSize
func SetPaginateData(in *PageOptions, records, maxId, minId int32) (paginate PageResult) {

	paginate.TotalRecords = records
	paginate.PageSize = in.PageSize
	paginate.PaginateMode = PaginateMode[in.PaginateMode]
	paginate.Action = PaginateAction[in.Action]

	if paginate.TotalRecords%in.PageSize == 0 {
		paginate.TotalPages = paginate.TotalRecords / in.PageSize
	} else {
		paginate.TotalPages = paginate.TotalRecords/in.PageSize + 1
	}

	switch PaginateMode[in.PaginateMode] {
	case "PageNumber":
		paginate.PageNumber = in.PageNumber
	case "MaxId":
		paginate.MaxId = maxId
		paginate.MinId = minId
	}

	return
}

// sort data : default id desc
// in.Action : 0-PreviousPage-asc , 1-NextPage-desc
// sort usersMap desc
func SetReturnDataSort(in *PageOptions, slice interface{}) {

	switch PaginateMode[in.PaginateMode] {
	case "PageNumber":
		// 不翻转
	case "MaxId":
		switch PaginateAction[in.Action] {
		case "NextPage":
			// 不翻转
		case "PreviousPage":
			helper.ReverseSlice(slice) // 翻转数组
		}
	}
}
