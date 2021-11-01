package utils

//分页
//row:  行数
//page: 页数
func Pagination(rows, page *int32) {
	//条案为空或者行数大于10 或乾小于０
	if rows == nil || *rows <= 0 {
		*rows = 10
	}

	//小于0或者没有页码
	if *page <= 0 || page == nil {
		*page = 0
	} else {
		*page = (*page - 1) * *rows
	}
}
