package codec

import (
	"regexp"
	"strconv"
	"strings"
	//"fmt"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-graph/graph"
	"github.com/spatialcurrent/go-graph/graph/exp"
	//"github.com/spatialcurrent/go-graph/graph/functions"
)

import (
	"github.com/spatialcurrent/sgol-codec/codec/update"
	"github.com/spatialcurrent/sgol-codec/codec/view"
)

type Parser struct {
	Clauses  []string `json:"clauses" hcl:"clauses"`
	Entities []string `json:"entities" hcl:"entities"`
	Edges    []string `json:"edges" hcl:"edges"`
}

func (p *Parser) Rejoin(block []string) string {
	text := ""

	for _, token := range block {
		if strings.Contains(token, " ") && !(strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"")) {
			text += "\"" + token + "\""
		} else {
			text += token + " "
		}
	}

	return text
}

func (p *Parser) ParseGroups(text string) []string {
	if len(text) > 0 {
		if strings.HasPrefix(text, "$") {
			if len(text) > 1 {
				if strings.Contains(text, ",") {
					parts := strings.Split(text, ",")
					result := make([]string, len(parts))
					for i, x := range parts {
						result[i] = x[1:len(x)]
					}
					return result
				} else {
					return []string{text[1:len(text)]}
				}
			} else {
				return p.Entities
			}
		} else if text == "_" || "text" == "-" {
			return []string{}
		} else {
			return []string{text}
		}
	} else {
		return []string{}
	}
}

func (p *Parser) ParseEntities(text string) []string {
	if len(text) > 0 {
		if strings.HasPrefix(text, "$") {
			if len(text) > 1 {
				if strings.Contains(text, ",") {
					parts := strings.Split(text, ",")
					result := make([]string, len(parts))
					for i, x := range parts {
						result[i] = x[1:len(x)]
					}
					return result
				} else {
					return []string{text[1:len(text)]}
				}
			} else {
				return p.Entities
			}
		} else if text == "_" || "text" == "-" {
			return []string{}
		} else {
			return []string{text}
		}
	} else {
		return []string{}
	}
}

func (p *Parser) ParseEdges(text string) []string {
	if len(text) > 0 {
		if strings.HasPrefix(text, "$") {
			if len(text) > 1 {
				if strings.Contains(text, ",") {
					parts := strings.Split(text, ",")
					result := make([]string, len(parts))
					for i, x := range parts {
						result[i] = x[1:len(x)]
					}
					return result
				} else {
					return []string{text[1:len(text)]}
				}
			} else {
				return p.Edges
			}
		} else {
			return []string{}
		}
	} else {
		return []string{}
	}
}

func (p *Parser) ParseUpdate(block []string) (string, map[string]exp.Node, error) {
	updateKey := block[0]
	filterFunctionsByGroup := map[string]exp.Node{}
	if len(block) > 3 {
		ff, err := exp.Parse(p.Rejoin(block[4:]))
		if err != nil {
			return updateKey, filterFunctionsByGroup, err
		}
		if ff != nil {
			filterFunctionsByGroup[block[2][1:]] = ff
		}
	}

	/*if len(block) > 6 {
		filterFunctions, err := p.ParseFilterFunctions(p.Rejoin(block[7:len(block)]))
		if err != nil {
			return op, err
		}
		if len(filterFunctions) > 0 {
			op.UpdateFilterFunctions = map[string][]functions.Filter{}
			op.UpdateFilterFunctions[block[5][1:len(block[5])]] = filterFunctions
		}
	}*/

	return updateKey, filterFunctionsByGroup, nil
}

func (p *Parser) ParseFilter(block []string) (exp.Node, map[string]exp.Node, string, map[string]exp.Node, error) {

  var filterFunctionForAll exp.Node
	filterFunctionsByGroup := map[string]exp.Node{}
	updateKey := ""
	updatefilterFunctionsByGroup := map[string]exp.Node{}

	if StringSliceContains(block, "UPDATE") {
		blockUpdateIndex := StringSliceIndex(block, "UPDATE")
		ff, err := exp.Parse(p.Rejoin(block[2:blockUpdateIndex]))
		if err != nil {
			return filterFunctionForAll, filterFunctionsByGroup, updateKey, updatefilterFunctionsByGroup, err
		}
		if ff != nil {
			filterFunctionsByGroup[block[0][1:len(block[0])]] = ff
		}
		updateKey, updatefilterFunctionsByGroup, err = p.ParseUpdate(block[blockUpdateIndex+1:])
	} else {
		if len(block) > 2 && block[1] == "WITH" {
			ff, err := exp.Parse(p.Rejoin(block[2:len(block)]))
			if err != nil {
				return filterFunctionForAll, filterFunctionsByGroup, updateKey, updatefilterFunctionsByGroup, err
			}
			if ff != nil {
				filterFunctionsByGroup[block[0][1:len(block[0])]] = ff
			}
		} else {
			ff, err := exp.Parse(p.Rejoin(block))
			if err != nil {
				return filterFunctionForAll, filterFunctionsByGroup, updateKey, updatefilterFunctionsByGroup, err
			}
			filterFunctionForAll = ff
		}
	}
	return filterFunctionForAll, filterFunctionsByGroup, updateKey, updatefilterFunctionsByGroup, nil
}

/*
func (p *Parser) ParseQueryFunctions(text string) ([]QueryFunction, error) {

	functions := make([]QueryFunction, 0)
	if len(text) > 0 {

		re_whitespace, err := regexp.Compile("(\\s+)")
		if err != nil {
			return functions, err
		}

		re, err := regexp.Compile("(\\s*)(?P<names>([a-zA-Z\\s]+))(\\s*)\\((\\s*)(?P<args>(.)*?)(\\s*)\\)(\\s*)")
		if err != nil {
			return functions, err
		}

		matches := re.FindAllStringSubmatch(strings.TrimSpace(text), -1)
		for _, m := range matches {
			g := map[string]string{}
			for i, name := range re.SubexpNames() {
				if i != 0 {
					g[name] = m[i]
				}
			}

			fn := QueryFunction{Names: re_whitespace.Split(strings.TrimSpace(g["names"]), -1)}
			if args_text, ok := g["args"]; ok {
				if len(args_text) > 0 {
					args := []string{}

					re2, err := regexp.Compile("(\\s*)(?P<value>((\"([^\"]+?)\")|([^,\\s]+)))(\\s*)")
					if err != nil {
						return functions, err
					}

					matches2 := re2.FindAllStringSubmatch(args_text, -1)
					for _, m2 := range matches2 {
						g2 := map[string]string{}
						for i, name := range re2.SubexpNames() {
							if i != 0 {
								g2[name] = m2[i]
							}
						}
						if value, ok := g2["value"]; ok {
							value = strings.TrimSpace(value)
							if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
								args = append(args, value[1:len(value)-1])
							} else {
								args = append(args, value)
							}
						}
					}
					fn.Args = args
				}
			}
			functions = append(functions, fn)
		}

	}

	return functions, nil
}*/

/*func (p *Parser) ParseFilterFunctions(text string) ([]functions.Filter, error) {

	filterFunctions := []functions.Filter{}
	queryFunctions, err := p.ParseQueryFunctions(text)
	if err != nil {
		return filterFunctions, err
	}
	if len(queryFunctions) > 0 {
		for _, qf := range queryFunctions {
			qf_name_lc := strings.ToLower(qf.Names)
			if qf_name_lc == "collectioncontains" {
				ff := functions.NewFilter([]string{qf.Args[0]}, qf.Name, []string{"value"}, map[string]string{
					"value": qf.Args[1],
				})
				filterFunctions = append(filterFunctions, ff)
			} else if qf_name_lc == "bbox" {
				ff := functions.NewFilter([]string{qf.Args[0]}, qf.Name, []string{"value"}, map[string]string{
					"value": qf.Args[1],
				})
				filterFunctions = append(filterFunctions, ff)
			}
		}
		return filterFunctions, nil
	} else {
		return []functions.Filter{}, nil
	}
}*/

/*func (p *Parser) ConvertToFilter(qf QueryFunction, i int) functions.Filter {
	qf_name_lc := strings.ToLower(qf.Names[i])
	if qf_name_lc == "not" {
		ff := functions.NewFilter([]string{qf.Args[0]}, qf.Name, []string{"value"}, map[string]string{
			"value": qf.Args[1],
		})
		filterFunctions = append(filterFunctions, ff)
	} else {

	}
}*/

func (p *Parser) ParseTokens(text string) ([]string, error) {
	tokens := []string{}

	re, err := regexp.Compile("(?P<token>((\"([^\"]+)\")|(\\S+)))")
	if err != nil {
		return tokens, err
	}

	matches := re.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		g := map[string]string{}
		for i, name := range re.SubexpNames() {
			if i != 0 {
				g[name] = m[i]
			}
		}
		tokens = append(tokens, strings.TrimSpace(g["token"]))
	}

	return tokens, nil

}

func (p *Parser) ParseBlocks(tokens []string) [][]string {
	blocks := [][]string{}
	block := []string{}

	for _, token := range tokens {
		token_uc := strings.ToUpper(token)
		if StringSliceContains(p.Clauses, token_uc) {
			if len(block) > 0 {
				blocks = append(blocks, block)
			}
			block = []string{token_uc}
		} else {
			block = append(block, token)
		}
	}

	if len(block) > 0 {
		blocks = append(blocks, block)
	}

	return blocks
}

func (p *Parser) ParseRun(block []string) (*OperationRun, error) {
	op := NewOperationRun()

	re, err := regexp.Compile("(\\s*)(?P<name>([a-zA-Z_\\d]+))(\\s*)\\((\\s*)(?P<args>(.)*?)(\\s*)\\)(\\s*)")
  if err != nil {
    return op, err
  }

  node, err := exp.ParseFunction(p.Rejoin(block[1:]), "", re)
	if err != nil {
		return op, err
	}
	op.SetFunction(node.(*exp.Function))

	return op, nil
}

func (p *Parser) ParseSelect(block []string) (*OperationSelect, error) {
	op := NewOperationSelect()

	if len(block) <= 1 {
		return op, errors.New("Invalid SELECT clause.  Must include at least 2 terms.")
	}

	if strings.HasPrefix(block[1], "$") {
		op.SetGroups(NewGroupCollection(block[1]))
	} else {
		op.SetSeeds(NewSeedCollection(block[1]))
	}

	if len(block) > 2 {
		block_2_uc := strings.ToUpper(block[2])
		if block_2_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[3:])
			if err != nil {
				return op, err
			}
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		} else if block_2_uc == "FILTER" {
			filterFunctionForAll, filterFunctionsByGroup, updateKey, updateFilterFunctions, err := p.ParseFilter(block[3:])
			if err != nil {
				return op, err
			}
			op.SetView(view.New(filterFunctionForAll, filterFunctionsByGroup))
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		}
	}

	return op, nil
}

func (p *Parser) ParseNav(block []string) (*OperationNav, error) {

	op := NewOperationNav()

	if block[1] == "INPUT" {
		op.SetSource(NewGroupCollection(block[1]))
		op.SetDestination(NewGroupCollection(block[3]))
		op.Direction = "OUTGOING"
		op.EdgeIdentifierToExtract = "DESTINATION"
	} else if block[3] == "INPUT" {
		op.SetSource(NewGroupCollection(block[1]))
		op.SetDestination(NewGroupCollection(block[3]))
		op.Direction = "INCOMING"
		op.EdgeIdentifierToExtract = "SOURCE"
	} else if !strings.HasPrefix(block[1], "$") {
		op.SetSeeds(NewSeedCollection(block[1]))
		op.SetDestination(NewGroupCollection(block[3]))
		op.Direction = "OUTGOING"
		op.EdgeIdentifierToExtract = "DESTINATION"
	} else if !strings.HasPrefix(block[3], "$") {
		op.SetSource(NewGroupCollection(block[1]))
		op.SetSeeds(NewSeedCollection(block[3]))
		op.Direction = "INCOMING"
		op.EdgeIdentifierToExtract = "SOURCE"
	} else {
		return op, errors.New("Could not parse NAV block "+strings.Join(block, " "))
	}

	if strings.Contains(block[2], "*") {
		block_2_parts := strings.Split(block[2], "*")
		op.SetEdges(NewGroupCollection(block_2_parts[0]))
		depth, err := strconv.Atoi(block_2_parts[1])
		if err != nil {
			return op, err
		}
		op.Depth = depth
	} else {
		op.SetEdges(NewGroupCollection(block[2]))
		op.Depth = 1
	}

	if len(block) > 4 {
		block_4_uc := strings.ToUpper(block[4])
		if block_4_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[5:])
			if err != nil {
				return op, err
			}
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		} else if block_4_uc == "FILTER" {
			filterFunctionForAll, filterFunctionsByGroup, updateKey, updateFilterFunctions, err := p.ParseFilter(block[5:])
			if err != nil {
				return op, err
			}
			op.SetView(view.New(filterFunctionForAll, filterFunctionsByGroup))
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		}
	}

	return op, nil
}

func (p *Parser) ParseHas(block []string) (*OperationHas, error) {

	op := NewOperationHas()

	if block[1] == "INPUT" {
		if strings.HasPrefix(block[3], "$") {
			return op, errors.New("HAS operation requires seeds.")
		}
		op.SetSource(NewGroupCollection(block[1]))
		op.SetSeeds(NewSeedCollection(block[3]))
		op.Direction = "OUTGOING"
		op.EdgeIdentifierToExtract = "SOURCE"
	} else if block[3] == "INPUT" {
		if strings.HasPrefix(block[1], "$") {
			return op, errors.New("HAS operation requires seeds.")
		}
		op.SetSeeds(NewSeedCollection(block[1]))
		op.SetDestination(NewGroupCollection(block[3]))
		op.Direction = "INCOMING"
		op.EdgeIdentifierToExtract = "DESTINATION"
	} else {
		return op, errors.New("Could not parse HAS block "+strings.Join(block, " "))
	}

	if strings.Contains(block[2], "*") {
		block_2_parts := strings.Split(block[2], "*")
		op.SetEdges(NewGroupCollection(block_2_parts[0]))
		depth, err := strconv.Atoi(block_2_parts[1])
		if err != nil {
			return op, err
		}
		op.Depth = depth
	} else {
		op.SetEdges(NewGroupCollection(block[2]))
		op.Depth = 1
	}

	if len(block) > 4 {
		block_4_uc := strings.ToUpper(block[4])
		if block_4_uc == "UPDATE" {
			updateKey, updateFilterFunctions, err := p.ParseUpdate(block[5:])
			if err != nil {
				return op, err
			}
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		} else if block_4_uc == "FILTER" {
			filterFunctionForAll, filterFunctionsByGroup, updateKey, updateFilterFunctions, err := p.ParseFilter(block[5:])
			if err != nil {
				return op, err
			}
			op.SetView(view.New(filterFunctionForAll, filterFunctionsByGroup))
			op.SetUpdate(update.New(updateKey, view.New(nil, updateFilterFunctions)))
		}
	}

	return op, nil
}

func (p *Parser) ParseOperations(blocks [][]string) ([]graph.Operation, string, error) {
	operations := make([]graph.Operation, 0)
	output_type := ""
	for _, block := range blocks {
		//fmt.Println("Parsing Block:", strings.Join(block, "; "))
		switch block[0] {
		case "INIT":
			operations = append(operations, NewOperationInit(block[1]))
		case "DISCARD":
			operations = append(operations, NewOperationDiscard())
		case "ADD":
			operations = append(operations, NewOperationAdd(true))
		case "RELATE":
			operations = append(operations, NewOperationRelate([]string{block[1]}))
		case "LIMIT":
			limit_int, err := strconv.Atoi(block[1])
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, NewOperationLimit(limit_int))
		case "SELECT":
			op, err := p.ParseSelect(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, *op)
		case "NAV":
			op, err := p.ParseNav(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, *op)
		case "HAS":
			op, err := p.ParseHas(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, *op)
		case "FETCH":
			operations = append(operations, NewOperationFetch(block[1]))
		case "SEED":
			operations = append(operations, NewOperationSeed())
		case "RUN":
			op, err := p.ParseRun(block)
			if err != nil {
				return operations, output_type, err
			}
			operations = append(operations, *op)
		case "OUTPUT":
			output_type = block[1]
		}
	}
	return operations, output_type, nil
}

func (p *Parser) ParseQuery(q string) (graph.OperationChain, error) {

	chain := OperationChain{}

	tokens, err := p.ParseTokens(q)
	if err != nil {
		return OperationChain{}, err
	}

	blocks := p.ParseBlocks(tokens)

	operations, outputType, err := p.ParseOperations(blocks)
	if err != nil {
		return chain, err
	}

	return NewOperationChain("chain", operations, outputType), nil
}
