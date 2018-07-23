package codec

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-graph/graph/elements"
)

type OperationAdd struct {
	*AbstractOperation
	Input bool `json:"input" bson:"input" yaml:"input" hcl:"input"`
	Entities []elements.Entity `json:"entities" bson:"entities" yaml:"entities" hcl:"entities"`
	Edges []elements.Edge `json:"edges" bson:"edges" yaml:"edges" hcl:"edges"`
}

func (op OperationAdd) HasInput() bool {
	return op.Input
}

func (op OperationAdd) Sgol() (string, error) {
	return "", errors.New("Add Operation cannot be represented as SGOL.")
}

func NewOperationAdd(input bool) OperationAdd {
	return OperationAdd{
		AbstractOperation: &AbstractOperation{Type: "ADD"},
		Input: input,
	}
}

func NewOperationAddWithElements(entities []elements.Entity, edges []elements.Edge) OperationAdd {
	return OperationAdd{
		AbstractOperation: &AbstractOperation{Type: "ADD"},
		Input: false,
		Entities: entities,
		Edges: edges,
	}
}
