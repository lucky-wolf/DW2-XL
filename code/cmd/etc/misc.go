package etc

// WHERE THE FUCK IS THE STANDARD DEFINITION OF THIS !@##@!q$
type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func MulDiv[T Int](value, numerator, denominator T) T {
	return value * numerator / denominator
}

func MulDivRoundUp[T Int](value, numerator, denominator T) T {
	return MulDiv(value+denominator-1, numerator, denominator)
}
