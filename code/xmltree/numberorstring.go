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

func (ns SimpleValue) Int() int {
	return int(ns.Float())
}
func (ns SimpleValue) Float() float64 {
	switch t := ns.value.(type) {
	case float64:
		return t
	default:
		return 0
	}
}
func (ns SimpleValue) String() string {
	switch t := ns.value.(type) {
	case float64:
		return fmt.Sprint(t)
	case string:
		return t
	default:
		return ""
	}
}

func (ns *SimpleValue) SetInt(v int) {
	ns.value = float64(v)
}
func (ns *SimpleValue) SetFloat(v float64) {
	ns.value = v
}
func (ns *SimpleValue) SetString(v string) {
	ns.value = v
}
