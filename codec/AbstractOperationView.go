package codec

import (
	"github.com/spatialcurrent/sgol-codec/codec/view"
)

type AbstractOperationView struct {
	View view.View `json:"view" bson:"view" yaml:"view" hcl:"view"`
}

func (op AbstractOperationView) HasFilters() bool {
	return op.View.HasFilters()
}

func (op *AbstractOperationView) SetView(v view.View) {
	op.View = v
}
