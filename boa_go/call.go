package main

type Callable interface{
  call(interpreter *Interpreter, args []interface{})
}

type BuiltInFunc struct {
  ident string
}

func (call *BuiltInFunc) call(_ *Interpreter, args []interface{}){

}

type CallableFunc struct {
  declaration *FunctionStatement
}

func (call *CallableFunc) call(interpreter *Interpreter, args []interface{}){
  env := NewEnv(interpreter.env)
  if len(call.declaration.args) != len(args) {
    panic("Wrong number of args to function: "+ call.declaration.ident.Lexeme.(string))
  }
  for i, arg := range args{
    param_name := call.declaration.args[i]
    env.define(param_name, arg)
  }
  interpreter.execute_block(call.declaration.body, env)
}
