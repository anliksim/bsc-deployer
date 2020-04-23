package v1

import "github.com/nvellon/hal"

type Health struct {
	Status string
}

func (p Health) GetMap() hal.Entry {
	return hal.Entry{
		"status": p.Status,
	}
}
