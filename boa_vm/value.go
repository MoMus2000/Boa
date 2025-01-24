package main

type ValueType int

const (
  VAL_NUMBER ValueType = iota
  VAL_BOOL
  VAL_NIL
)

type Value struct {
  valType ValueType
  boolean bool
  number  float32
}

type ValueArray struct {
  values []Value
}

func NewValue() ValueArray {
  return ValueArray{
    values : make([]Value, 0),
  }
}

func BoolVal(v bool) Value {
  return Value{boolean: v, valType: VAL_BOOL}
}

func NumberVal(v float32) Value {
  return Value{number: v, valType: VAL_NUMBER}
}

func NilVal() Value {
  return Value{valType: VAL_NIL}
}

func (v *Value) AsBoolean() bool {
  return v.boolean
}

func (v *Value) AsNumber() float32 {
  return v.number
}

func (v *Value) IsBool() bool {
  return v.valType == VAL_BOOL
}
func (v *Value) IsNumber() bool {
  return v.valType == VAL_NUMBER
}
func (v *Value) IsNil() bool {
  return v.valType == VAL_NIL
}

func (v *ValueArray) WriteValue(value Value) {
  v.values = append(v.values, value)
}

func (v *ValueArray) FreeValue() {
  v.values = make([]Value, 0)
}

