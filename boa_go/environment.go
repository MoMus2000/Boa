package main

type Env struct{
  mapper       map[string]interface{}
  enclosed_env *Env
}

func NewEnv(enclosing Env) *Env{
  return &Env{
    mapper: make(map[string]interface{}),
    enclosed_env: &enclosing,
  }
}

func (e *Env) define() {

}

func (e *Env) get() {

}

func (e *Env) assign() {

}

