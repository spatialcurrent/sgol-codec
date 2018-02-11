package codec

type OperationFetch struct {
	*AbstractOperationKey
}

func (op OperationFetch) Sgol() (string, error) {
	return "FETCH "+op.Key, nil
}
