package main

import "time"

type Time struct{}

func (t *Time) arity() int {
  return 0
}

func (t *Time) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
  return float64(time.Now().Unix()), nil
} 

func (t *Time) String() string {
  return "<native fn>"
}
