package logic

import (
	"fmt"
	"strconv"
	"strings"
)

type Value []bool

func (v Value) String() string {
	n := 0
	for i, b := range v {
		if b {
			n += 1 << i
		}
	}
	return strconv.Itoa(n)
}

type Constants struct {
	values []Value
}

func NewConstants(value []Value) *Constants {
	return &Constants{value}
}

func (c *Constants) Values() []Value {
	return c.values
}

func (c *Constants) Name() string {
	return "constants"
}

func (c *Constants) InputNames() []string {
	return nil
}

func (c *Constants) OutputNames() []string {
	names := make([]string, len(c.values))
	for i := range c.values {
		names[i] = fmt.Sprintf("c%d", i)
	}
	return names
}

func (c *Constants) String() string {
	var sb strings.Builder
	for i, value := range c.values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(value.String())
	}
	return sb.String()
}
