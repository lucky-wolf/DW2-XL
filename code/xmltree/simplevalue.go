package xmltree

import "fmt"

// we encapsulate a number or a string (useful for conveying simple xml values)
type SimpleValue struct {
	value any
}

func CreateInt(value int) SimpleValue {
	return CreateNumber(float64(value))
}
func CreateNumber(value float64) SimpleValue {
	return SimpleValue{value: value}
}
func CreateString(value string) SimpleValue {
	return SimpleValue{value: value}
}

func (sv SimpleValue) ApplyTo(e *XMLValue) {
	if e != nil {
		e.SetString(sv.String())
	}
}

func (sv SimpleValue) Int() int {
	switch t := sv.value.(type) {
	case float64:
		return int(t)
	case int:
		return t
	default:
		return 0
	}
}

func (sv SimpleValue) Float() float64 {
	switch t := sv.value.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	default:
		return 0
	}
}
func (sv SimpleValue) String() string {
	switch t := sv.value.(type) {
	case float64:
		return fmt.Sprint(t)
	case int:
		return fmt.Sprint(t)
	case string:
		return t
	default:
		return ""
	}
}

func (sv *SimpleValue) SetInt(v int) {
	sv.value = float64(v)
}
func (sv *SimpleValue) SetFloat(v float64) {
	sv.value = v
}
func (sv *SimpleValue) SetString(v string) {
	sv.value = v
}
