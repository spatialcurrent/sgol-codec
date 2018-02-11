package codec

import (
	"strconv"
)

type OperationLimit struct {
	*AbstractOperation
	Limit int `json:"limit" bson:"limit" yaml:"limit" hcl:"limit"`
}

func (op OperationLimit) Sgol() (string, error) {
	return "LIMIT "+strconv.Itoa(op.Limit), nil
}
