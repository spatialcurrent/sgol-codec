package codec

import (
	"github.com/spatialcurrent/go-graph/graph/exp"
)


type OperationRun struct {
	*AbstractOperation
	Function *exp.Function `json:"function" bson:"function" yaml:"function" hcl:"function"`
}

func NewOperationRun() *OperationRun {
	return &OperationRun{
		AbstractOperation: &AbstractOperation{Type: "RUN"},
	}
}

func (op OperationRun) Sgol() (string, error) {
	return "RUN "+op.Function.Sgol(), nil
}

func (op *OperationRun) SetFunction(f *exp.Function) {
	op.Function = f
}
