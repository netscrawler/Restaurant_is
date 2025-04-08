package models

import "fmt"

type CodeMsg struct {
	code int
}

func NewCodeMsg(code int) *CodeMsg {
	return &CodeMsg{code: code}
}

func (c *CodeMsg) String() string {
	return fmt.Sprintf("Your verification code: %d", c.code)
}
