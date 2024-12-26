package main

type GloxCallable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
}

type GloxFunction struct {
	Declaration StmtFunction
}

func (f GloxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewEnvironment(interpreter.globals)
	for i, paramName := range f.Declaration.Params {
		environment.define(paramName.Lexeme, arguments[i])
	}
	err := interpreter.executeBlock(f.Declaration.Body, environment)
	if err == nil {
		return nil, nil
	} else if ret, ok := err.(Return); ok {
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
