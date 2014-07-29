package vutil

import (
	"errors"
	"fmt"
)

func T(exp bool, a interface{}, b interface{}) interface{} {
	if exp {
		return a
	}
	return b
}

type CAPArray struct {
	N int
	a []interface{}
}

func (c *CAPArray) init() {
	if c.a == nil {
		c.a = make([]interface{}, 0, c.N+1)
	}
}

func (c *CAPArray) Fill(e interface{}) {
	for i := 0; i < c.N; i++ {
		c.Push(e)
	}
}

func (c *CAPArray) Geth(i int) (interface{}, error) {
	if len(c.a) != c.N {
		return nil, errors.New("CAPArray must be full")
	}
	if i <= 0 || i > c.N {
		return nil, errors.New("index out of range")
	}
	return c.a[c.N-i], nil
}

func (c *CAPArray) Seth(i int, e interface{}) error {
	if len(c.a) != c.N {
		return errors.New("CAPArray must be full")
	}
	c.a[c.N-i] = e
	return nil
}

func (c *CAPArray) Push(e interface{}) {
	c.init()
	c.a = append(c.a, e)
	if len(c.a) > c.N {
		c.a = c.a[1:len(c.a)]
	}
}

func (c *CAPArray) Pop() (e interface{}) {
	if c.a == nil {
		return
	}
	e = c.a[len(c.a)-1]
	c.a = c.a[:len(c.a)-1]
	return e
}

func (c *CAPArray) Unshift(e interface{}) {
	c.init()
	c.a = append(c.a, e)
	copy(c.a[1:], c.a[0:])
	c.a[0] = e
	if len(c.a) > c.N {
		c.a = c.a[0 : len(c.a)-1]
	}
}

func (c *CAPArray) Shift() (e interface{}) {
	if c.a == nil {
		return
	}
	e = c.a[0]
	c.a = c.a[1:]
	return e
}

func (c *CAPArray) Dump(f string) {
	if len(c.a) > 0 {
		fmt.Printf(fmt.Sprintf("n=%s first=%s last=%s { ", f, f, f), len(c.a), c.a[0], c.a[len(c.a)-1])
		for _, v := range c.a {
			fmt.Printf(f, v)
		}
		fmt.Println("}")
	} else {
		fmt.Println("{}")
	}
}
