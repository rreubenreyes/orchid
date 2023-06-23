package dag

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	Rules  []Rule `json:"rules"`
	IsEnd  bool   `json:"-"`
	IsWait bool   `json:"-"`
}

type Rule struct {
	Next string `json:"next"`
}

type DAG map[string]Node

func detectCycle(dag DAG) (bool, []string) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	cycle := make([]string, 0)

	for node := range dag {
		if !visited[node] {
			if detectCycleRecurse(dag, node, visited, recStack, &cycle) {
				return true, cycle
			}
		}
	}

	return false, nil
}

func detectCycleRecurse(dag DAG, node string, visited, recStack map[string]bool, cycle *[]string) bool {
	visited[node] = true
	recStack[node] = true

	for _, rule := range dag[node].Rules {
		nextNode := rule.Next

		if !visited[nextNode] {
			if detectCycleRecurse(dag, nextNode, visited, recStack, cycle) {
				*cycle = append(*cycle, nextNode)
				return true
			}
		} else if recStack[nextNode] {
			*cycle = append(*cycle, nextNode)
			return true
		}
	}

	recStack[node] = false
	return false
}

func Validate(data string) error {
	var dag DAG

	err := json.Unmarshal([]byte(data), &dag)
	if err != nil {
		return errors.New("invalid DAG definition")
	}

	// DAG cannot be empty
	if len(dag) <= 0 {
		return errors.New(`DAG must contain at least one node`)
	}

	// DAG must contain a node called "start"
	_, ok := dag["start"]
	if !ok {
		return errors.New(`DAG must contain a node called "start"`)
	}

	// DAG cannot contain reserved node names
	_, ok = dag["_wait"]
	if !ok {
		return errors.New(`node cannot have reserved name "_wait"`)
	}
	_, ok = dag["_end"]
	if !ok {
		return errors.New(`node cannot have reserved name "_end"`)
	}

	// all user-provided nodes must contain at least one rule
	for nk, nv := range dag {
		if len(nv.Rules) <= 0 {
			return fmt.Errorf(`node "%s" must contain at least one rule`, nk)
		}
	}

	// add reserved "_wait" and "_end" nodes for use by workflow runner
	dag["_wait"] = Node{
		Rules:  []Rule{},
		IsWait: true,
		IsEnd:  false,
	}
	dag["_end"] = Node{
		Rules:  []Rule{},
		IsWait: false,
		IsEnd:  true,
	}

	// DAG cannot be cyclic
	isCyclic, cycle := detectCycle(dag)
	if isCyclic {
		return fmt.Errorf(`cycle detected: %s`, strings.Join(cycle, "->"))
	}

	// DAG nodes cannot be isolated. A DAG is considered to have isolated nodes when:
	// - The DAG consists of more than a single node
	// - One or more of those nodes cannot be reached by any means

	return nil
}