package update

import (
	"github.com/spatialcurrent/sgol-codec/codec/view"
)

type Update struct {
	Key  string `json:"key" bson:"key" yaml:"key" hcl:"key"`
	View view.View   `json:"view" bson:"view" yaml:"view" hcl:"view"`
}

func New(key string, v view.View) Update {
	return Update{
		Key:  key,
		View: v,
	}
}

func (u Update) HasFilters() bool {
	return u.View.HasFilters()
}

func (u Update) Sgol() string {

	s := "UPDATE " + u.Key

	if u.HasFilters() {
		s += u.View.Sgol()
	}

	return s
}
