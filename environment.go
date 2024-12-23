package main

type Environment struct {
  values map[string]interface{}
}

func NewEnvironment() Environment {
  return Environment{values: make(map[string]interface{})}
}

func (e *Environment) define(name string, value interface{}) {
  e.values[name] = value
}

func (e *Environment) assign(name Token, value interface{}) error {
  _, ok := e.values[name.Lexeme]
  if !ok {
    return &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
  }
  e.values[name.Lexeme] = value
  return nil
}

func (e *Environment) get(name Token) (interface{}, error) {
  value, ok := e.values[name.Lexeme]
  if !ok {
    return nil, &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
  }
  return value, nil
}
