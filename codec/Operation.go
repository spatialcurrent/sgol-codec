package codec

import (
	"github.com/spatialcurrent/go-graph/graph"
)

type Operation interface {
	GetTypeName() string
	Validate(schema graph.Schema) error
	GetInputType() bool
	GetOutputType() bool
	Sgol() (string, error)
}

type AbstractOperation struct {
	Type string `json:"type" bson:"type" yaml:"type" hcl:"type"`
}

func (op AbstractOperation) GetTypeName() string {
	return op.Type
}

func (op AbstractOperation) Validate(schema graph.Schema) error {
	return nil
}

func (op AbstractOperation) GetInputType() string {
	return "void"
}

func (op AbstractOperation) GetOutputType() string {
	return "void"
}
