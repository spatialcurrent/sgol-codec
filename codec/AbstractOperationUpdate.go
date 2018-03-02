package codec

import (
	"github.com/spatialcurrent/sgol-codec/codec/update"
)

type AbstractOperationUpdate struct {
	Update update.Update `json:"update" bson:"update" yaml:"update" hcl:"update"`
}

func (op AbstractOperationUpdate) HasUpdate() bool {
	return len(op.Update.Key) > 0
}

func (op AbstractOperationUpdate) HasUpdateFilters() bool {
	return op.Update.HasFilters()
}

func (op *AbstractOperationUpdate) SetUpdate(u update.Update) {
	op.Update = u
}
