package main

import "fmt"

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (env *Environment) get(token Token) (interface{}, error) {
	if value, ok := env.values[token.Lexeme]; ok {
		return value, nil
	} else if env.enclosing != nil {
		return env.enclosing.get(token)
	} else {
		return nil, &RuntimeError{
			token:   token,
			message: fmt.Sprintf("Undefined variable '%v'", token.Lexeme),
		}
	}
}

func (env *Environment) assign(token Token, value interface{}) error {
	if _, ok := env.values[token.Lexeme]; ok {
		env.values[token.Lexeme] = value
		return nil
	} else if env.enclosing != nil {
		return env.enclosing.assign(token, value)
	} else {
		return &RuntimeError{
			token:   token,
			message: fmt.Sprintf("Undefined variable '%v'", token.Lexeme),
		}
	}
}

func (env *Environment) getAt(distance int, name string) interface{} {
  return env.ancestor(distance).values[name]
}

func (env *Environment) assignAt(distance int, name Token, value interface{}) {
  env.ancestor(distance).values[name.Lexeme] = value
}

func (env *Environment) ancestor(distance int) *Environment {
  var _env = env
	for i := 0; i < distance; i-- {
    _env = _env.enclosing
  }
  return _env
}
