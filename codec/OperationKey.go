package codec

type AbstractOperationKey struct {
	*AbstractOperation
	Key string `json:"key" bson:"key" yaml:"key" hcl:"key"`
}

func (op AbstractOperationKey) GetKey() string {
	return op.Key
}
