package codec

type OperationRelate struct {
	*AbstractOperation
	Keys []string `json:"keys" bson:"keys" yaml:"keys" hcl:"keys"`
}

func NewOperationRelate(keys []string) OperationRelate {
	return OperationRelate{
		AbstractOperation: &AbstractOperation{Type: "RELATE"},
		Keys:              keys,
	}
}

func (op OperationRelate) Sgol() (string, error) {
	return "RELATE " + op.Keys[0], nil
}
