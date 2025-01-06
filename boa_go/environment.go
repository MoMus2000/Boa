package main

type Env struct{
  mapper       map[string]interface{}
  enclosed_env *Env
}

func NewEnv(enclosing *Env) *Env{
  return &Env{
    mapper        : make(map[string]interface{}),
    enclosed_env  : enclosing,
  }
}

func (e *Env) define(name string, value interface{}) {
  e.mapper[name] = value
}

func (e *Env) get(name string) interface{}{
  val, exists := e.mapper[name]
  if !exists {
    panic("Error: undefined var")
  }
  if e.enclosed_env != nil {
    return e.enclosed_env.get(name)
  }
  return val
}

func (e *Env) assign(name string, value interface{}) {
  _ , exists := e.mapper[name]
  if exists {
    e.mapper[name] = value
    return
  }
  if e.enclosed_env != nil {
    e.enclosed_env.assign(name, value)
    return
  }
  if !exists {
    panic("Error: undefined var")
  }
}

