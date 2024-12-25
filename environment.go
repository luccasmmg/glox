package main

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
	globals   *Environment
}

func GlobalEnvironment() *Environment {
  return &Environment{values: make(map[string]interface{}), enclosing: nil, globals: nil}
}

func NewEnvironment(enclosing ...*Environment) Environment {
	var env *Environment
	if len(enclosing) > 0 {
		env = enclosing[0]
	}
	global := GlobalEnvironment()
	global.define("clock", Time{})
	return Environment{values: make(map[string]interface{}), enclosing: env, globals: global}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) assign(name Token, value interface{}) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			e.enclosing.assign(name, value)
			return nil
		}
		return &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
	}
	e.values[name.Lexeme] = value
	return nil
}

func (e *Environment) get(name Token) (interface{}, error) {
	value, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
	}
	return value, nil
}
