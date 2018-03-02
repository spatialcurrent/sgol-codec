package codec

type OperationAdd struct {
	*AbstractOperationKey
}

func NewOperationAdd(key string) OperationAdd {
	return OperationAdd{
		&AbstractOperationKey{
			AbstractOperation: &AbstractOperation{Type: "ADD"},
			Key:               key,
		},
	}
}

func (op OperationAdd) Sgol() (string, error) {
	return "ADD " + op.Key, nil
}
