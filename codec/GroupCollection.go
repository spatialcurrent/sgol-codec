package codec

import (
	"strings"
)

type GroupCollection struct {
	Input  bool     `json:"input" bson:"input" yaml:"input" hcl:"input"`
	All    bool     `json:"all" bson:"all" yaml:"all" hcl:"all"`
	Groups []string `json:"groups" bson:"groups" yaml:"groups" hcl:"groups"`
}

func NewGroupCollection(text string) GroupCollection {
	gc := GroupCollection{All: false, Groups: make([]string, 0)}

	if len(text) > 0 {
		if text == "INPUT" {
			gc.Input = true
		} else {
			if strings.HasPrefix(text, "$") {
				if len(text) > 1 {
					if strings.Contains(text, ",") {
						parts := strings.Split(text, ",")
						result := make([]string, len(parts))
						for i, x := range parts {
							result[i] = x[1:len(x)]
						}
						gc.Groups = result
					} else {
						gc.Groups = []string{text[1:len(text)]}
					}
				} else {
					gc.All = true
				}
			}
		}
	}

	return gc
}

func (gc GroupCollection) GetGroups() []string {
	return gc.Groups
}

func (gc GroupCollection) HasGroup(group string) bool {
	return gc.All || StringSliceContains(gc.Groups, group)
}

func (gc GroupCollection) Sgol() string {
	if gc.Input {
		return "INPUT"
	} else {
		if gc.All {
			return "$"
		} else {
			s := ""
			for i, x := range gc.Groups {
				if i > 0 {
					s += ","
				}
				s += "$" + x
			}
			return s
		}
	}
}
