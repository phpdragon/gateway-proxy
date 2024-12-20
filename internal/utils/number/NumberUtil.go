package number

import "strconv"

func Str2int(str string) (int, error) {
	return strconv.Atoi(str)
}

func Str2int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Int2Str(n int) string {
	return strconv.Itoa(n)
}

func Int64ToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
