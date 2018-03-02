package codec

import (
	"strconv"
)

type OperationLimit struct {
	*AbstractOperation
	Limit int `json:"limit" bson:"limit" yaml:"limit" hcl:"limit"`
}

func NewOperationLimit(limit int) OperationLimit {
	return OperationLimit{
		AbstractOperation: &AbstractOperation{Type: "LIMIT"},
		Limit:             limit,
	}
}

func (op OperationLimit) GetLimit() int {
	return op.Limit
}

func (op OperationLimit) Sgol() (string, error) {
	return "LIMIT " + strconv.Itoa(op.Limit), nil
}
