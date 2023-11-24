package errorModel

type Err int64

const (
	RequestErr Err = 0
	ParamErr   Err = 1
)

func (tr Err) String() string {
	switch tr {
	case 0:
		return "Network Error"
	case 1:
		return "Param Error"
	default:
		return ""
	}
}
