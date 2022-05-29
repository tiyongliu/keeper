package backend

import (
	"context"
	"keeper/app/pkg/serializer"
)

type MMMM struct {
	ctx context.Context
}

func NewMMMM() *MMMM {
	return &MMMM{}
}

func (m *MMMM) GetVersion(params map[string]interface{}) *serializer.Response {
	return nil
}
