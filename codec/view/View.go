package view

import (
	"github.com/spatialcurrent/go-graph/graph/exp"
)

type View struct {
	FilterFunctionForAll  exp.Node            `json:"filter_function_for_all" bson:"filter_function_for_all" yaml:"filter_function_for_all" hcl:"filter_function_for_all"`
	FilterFunctionsByGroup map[string]exp.Node `json:"filter_functions_by_group" bson:"filter_functions_by_group" yaml:"filter_functions_by_group" hcl:"filter_functions_by_group"`
}

func New(filterFunctionForAll exp.Node, filterFunctionsByGroup map[string]exp.Node) View {
	return View{
		FilterFunctionForAll:  filterFunctionForAll,
		FilterFunctionsByGroup: filterFunctionsByGroup,
	}
}

func (v View) HasFilters() bool {

	if v.FilterFunctionForAll != nil {
		return true
	}

	return len(v.FilterFunctionsByGroup) > 0
}

func (v View) Sgol() string {
	s := ""

  if v.FilterFunctionForAll != nil {
		s += "FILTER " + v.FilterFunctionForAll.Sgol()
	}

	if len(v.FilterFunctionsByGroup) > 0 {
		for g, ff := range v.FilterFunctionsByGroup {
			s += " FILTER $" + g + " WITH "+ff.Sgol()
		}
	}

	return s
}
