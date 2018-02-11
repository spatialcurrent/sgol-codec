package codec

type OperationRelate struct {
	*AbstractOperation
	Keys []string `json:"keys" bson:"keys" yaml:"keys" hcl:"keys"`
}

func (op OperationRelate) Sgol() (string, error) {
	return "RELATE "+op.Keys[0], nil
}
