// 框架函数库

package helper

// 判断数组下标是否存在
func IssetArrayIndex(arr []interface{}, index int) bool {
	for i := range arr {
		if i == index {
			return true
		}
	}

	return false
}
