package main

import (
	"fmt"
)

type GloxClass struct {
	Name    string
	Methods map[string]GloxFunction
}

type GloxInstance struct {
	Klass  *GloxClass
	Fields map[string]interface{}
}

func NewGloxClass(name string, methods map[string]GloxFunction) GloxClass {
	return GloxClass{
		Name:    name,
		Methods: methods,
	}
}

func NewGloxInstance(klass GloxClass) GloxInstance {
	return GloxInstance{
		Klass:  &klass,
		Fields: make(map[string]interface{}),
	}
}

func (c GloxClass) String() string {
	return fmt.Sprintf("%s", c.Name)
}

func (c *GloxClass) FindMethod(name string) *GloxFunction {
	if value, ok := c.Methods[name]; ok {
		return &value
	} else {
    return nil
	}
}

func (i GloxInstance) String() string {
	return fmt.Sprintf("%s Instance", i.Klass.Name)
}

func (f GloxClass) Arity() int {
  if initializer, ok := f.Methods["init"]; ok {
    return initializer.Arity()
  }
	return 0
}

func (f GloxClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewGloxInstance(f)
  if initializer, ok := f.Methods["init"]; ok {
    _, err := initializer.Bind(&instance).Call(interpreter, arguments)
    if err != nil {
      return nil, err
    }
  }
	return instance, nil
}

func (i *GloxInstance) Get(name Token) (interface{}, error) {
	if value, ok := i.Fields[name.Lexeme]; ok {
		return value, nil
	} else {
    method := i.Klass.FindMethod(name.Lexeme)
    if method != nil {
      return method.Bind(i), nil
    }
		return nil, &RuntimeError{
			token:   name,
			message: fmt.Sprintf("Undefined propety '%v'", name.Lexeme),
		}
	}
}

func (i *GloxInstance) Set(name Token, value interface{}) {
	i.Fields[name.Lexeme] = value
}
