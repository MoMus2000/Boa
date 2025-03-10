package main

import "unsafe"

type ValueType int

const (
	VAL_NUMBER ValueType = iota
	VAL_BOOL
	VAL_NIL
	VAL_OBJ
)

type Value struct {
	valType ValueType
	boolean bool
	number  float32
	obj     *Object
}

type ValueArray struct {
	values []Value
}

func NewValue() ValueArray {
	return ValueArray{
		values: make([]Value, 0),
	}
}

func BoolVal(v bool) Value {
	return Value{boolean: v, valType: VAL_BOOL}
}

func NumberVal(v float32) Value {
	return Value{number: v, valType: VAL_NUMBER}
}

func (value *Value) isString() bool {
	return value.isObjectType(OBJ_STRING)
}

func (value *Value) isFunc() bool {
	return value.isObjectType(OBJ_FUNC)
}

func (value *Value) asFunc() *ObjectFunc {
	return (*ObjectFunc)(unsafe.Pointer(value.AsObj()))
}

func (value *Value) asString() *ObjectString {
	return (*ObjectString)(unsafe.Pointer(value.AsObj()))
}

func (value *Value) asCString() string {
	return (*ObjectString)(unsafe.Pointer(value.AsObj())).chars
}

func ObjVal(v *Object) Value {
	return Value{obj: v, valType: VAL_OBJ}
}

func (value *Value) isObjectType(objType ObjType) bool {
	return value.IsObj() && value.AsObj().objType == objType
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

func (v *Value) AsObj() *Object {
	return v.obj
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

func (v *Value) IsObj() bool {
	return v.valType == VAL_OBJ
}

func (v *ValueArray) WriteValue(value Value) {
	v.values = append(v.values, value)
}

func (v *ValueArray) FreeValue() {
	v.values = make([]Value, 0)
}
