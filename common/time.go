package common

//按时间粒度取时间的格式
func GetTimeFormat(period int) string {
	switch period {
	case 1, 5:
		return "2006-01-02 15:04"
	case 10, 15, 30:
		return "2006-01-02 15:0"
	case 60:
		return "2006-01-02 15"
	case 1440:
		return "2006-01-02"
	default:
		return "2006-01-02 15:04"
	}
}
