package model

import "github.com/nvellon/hal"

type None struct {
}

func (p None) GetMap() hal.Entry {
	return hal.Entry{}
}
