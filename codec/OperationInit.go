package codec

type OperationInit struct {
	*AbstractOperationKey
}

func (op OperationInit) Sgol() (string, error) {
	return "INIT "+op.Key, nil
}
