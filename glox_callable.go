package main

type GloxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
}

type GloxFunction struct {
	IsInitializer bool
	Declaration   StmtFunction
	Closure       *Environment
}

func (f GloxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(f.Closure)
	for i, paramName := range f.Declaration.Params {
		environment.define(paramName.Lexeme, arguments[i])
	}
	err := interpreter.executeBlock(f.Declaration.Body, &environment)
	if err == nil {
		return nil, nil
	} else if ret, ok := err.(Return); ok {
    if f.IsInitializer {
      return f.Closure.getAt(0, "this"), nil
    }
		return ret.Value, nil
	} else {
		return nil, err
	}
}

func (f GloxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f GloxFunction) String() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}

func (f GloxFunction) Bind(instance *GloxInstance) GloxFunction {
	environment := NewEnvironment(f.Closure)
	environment.define("this", instance)
	return GloxFunction{
		Declaration: f.Declaration, Closure: &environment, IsInitializer: f.IsInitializer}
}
