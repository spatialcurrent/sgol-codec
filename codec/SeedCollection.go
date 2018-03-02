package codec

import (
	"strings"
)

type SeedCollection struct {
	Vertices []string `json:"verticies" bson:"verticies" yaml:"verticies" hcl:"vertitices"`
}

func NewSeedCollection(text string) SeedCollection {
	sc := SeedCollection{}
	sc.Vertices = strings.Split(text, ",")
	return sc
}

func (sc SeedCollection) Sgol() string {
	return strings.Join(sc.Vertices, ",")
}
