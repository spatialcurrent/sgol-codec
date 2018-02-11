package codec

type Operation interface{
	GetTypeName() string
	Sgol() (string, error)
}

type AbstractOperation struct{
	Type string `json:"type" bson:"type" yaml:"type" hcl:"type"`
}

func (op *AbstractOperation) GetTypeName() string {
  return op.Type
}

type OperationPlaceholder struct{
	*AbstractOperation
}

func (op OperationPlaceholder) Sgol() (string, error) {
	return "PLACEHOLDER", nil
}

type AbstractOperationKey struct {
	*AbstractOperation
  Key string `json:"key" bson:"key" yaml:"key" hcl:"key"`
}

type OperationAdd struct {
	*AbstractOperationKey
}

func (op OperationAdd) Sgol() (string, error) {
	return "ADD "+op.Key, nil
}

type OperationDiscard struct {
	*AbstractOperationKey
}

func (op OperationDiscard) Sgol() (string, error) {
	return "DISCARD "+op.Key, nil
}
