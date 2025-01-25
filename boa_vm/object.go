package main

type ObjType int

const (
  OBJ_STRING ObjType = iota
)

type Object struct {
  objType ObjType
}

func (o *Object) ObjType() ObjType {
  return o.objType
}

type ObjectString struct {
  obj    Object
  length int
  chars  string
}

