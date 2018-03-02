package codec

type OperationInit struct {
	*AbstractOperationKey
}

func NewOperationInit(key string) OperationInit {
	return OperationInit{
		&AbstractOperationKey{
			AbstractOperation: &AbstractOperation{Type: "DISCARD"},
			Key:               key,
		},
	}
}

func (op OperationInit) Sgol() (string, error) {
	return "INIT " + op.Key, nil
}
