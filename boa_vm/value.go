package main

type Value float32;

type ValueArray struct {
  values []Value
}

func NewValue() ValueArray {
  return ValueArray{
    values : make([]Value, 0),
  }
}

func (v *ValueArray) WriteValue(value Value) {
  v.values = append(v.values, value)
}

func (v *ValueArray) FreeValue() {
  v.values = make([]Value, 0)
}

