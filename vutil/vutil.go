package vutil

import (
	"errors"
)

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

type CAPArray struct {
	N int
	A []interface{}
}

func (c *CAPArray) init() {
	if c.A == nil {
		c.A = make([]interface{}, 0, c.N+1)
	}
}

func (c *CAPArray) Geth(i int) (interface{}, error) {
	if len(c.A) != c.N {
		return nil, errors.New("CAPArray must be full")
	}
	if i <= 0 || i > c.N {
		return nil, errors.New("index out of range")
	}
	return c.A[c.N-i], nil
}

func (c *CAPArray) Seth(i int, e interface{}) error {
	if len(c.A) != c.N {
		return errors.New("CAPArray must be full")
	}
	c.A[c.N-i] = e
	return nil
}

func (c *CAPArray) Push(e interface{}) {
	c.init()
	c.A = append(c.A, e)
	if len(c.A) > c.N {
		c.A = c.A[1:len(c.A)]
	}
}

func (c *CAPArray) Pop() (e interface{}) {
	if c.A == nil || len(c.A) == 0 {
		return nil
	}
	e = c.A[len(c.A)-1]
	c.A = c.A[:len(c.A)-1]
	return e
}

func (c *CAPArray) Unshift(e interface{}) {
	c.init()
	c.A = append(c.A, e)
	copy(c.A[1:], c.A[0:])
	c.A[0] = e
	if len(c.A) > c.N {
		c.A = c.A[0 : len(c.A)-1]
	}
}

func (c *CAPArray) Shift() (e interface{}) {
	if c.A == nil || len(c.A) == 0 {
		return nil
	}
	e = c.A[0]
	c.A = c.A[1:]
	return e
}
