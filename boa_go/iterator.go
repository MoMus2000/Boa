package main

type Iterator interface {
  next() (any, bool)
  reset()
}

