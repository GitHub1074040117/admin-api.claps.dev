package util

// 自定义错误类型
type Err struct {
	Code    int    // 错误类型
	Message string // 错误信息提示
	Errord  error  // 内部错误信息
}

const (
	OK    = 200
	ERROR = 400
)

// 返回成功信息
func Success() *Err {
	var err Err
	err.Code = OK
	err.Message = ""
	err.Errord = nil
	return &err
}

// 返回失败信息
func Fail(msg string) *Err {
	var err Err
	err.Code = ERROR
	err.Message = msg
	return &err
}

func IsOk(err *Err) bool {
	if err.Code == OK {
		return true
	}
	return false
}
