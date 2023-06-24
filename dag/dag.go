package dag

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Next string `json:"next"`
	Wait bool   `json:"wait"`
	End  bool   `json:"end"`
}

type DAG map[string]Node

func detectCycle(dag DAG, node string, visited, recStack map[string]bool, cycle *[]string) bool {
	visited[node] = true  // keep track of all nodes which were ever visited
	recStack[node] = true // keep track of nodes that were visited as part of recursion starting from an iteration in `detectCycle`

	// rules without a "next" do not specify a traversal, so ignore them
	var nonTerminalRules []Rule
	for _, rule := range dag[node].Rules {
		if rule.Next != "" {
			nonTerminalRules = append(nonTerminalRules, rule)
		}
	}

	// for every rule that specifies a traversal, perform a depth-first search to find cycles
	for _, rule := range nonTerminalRules {
		nextNode := rule.Next

		if !visited[nextNode] {
			if detectCycle(dag, nextNode, visited, recStack, cycle) {
				// if the next node has not been visited, then visit it, and continue the traversal
				*cycle = append(*cycle, nextNode)
				return true
			}
		} else if recStack[nextNode] {
			// else if the next node has been visited, we may not have yet seen it as part of _this_ traversal,
			// and this may be a cross-edge, which does not indicate a cycle.
			*cycle = append(*cycle, nextNode)
			return true
		}
	}

	recStack[node] = false
	return false
}

func assertAcyclic(dag DAG) error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	cycle := make([]string, 0)

	for node := range dag {
		if visited[node] {
			continue
		}
		if detectCycle(dag, node, visited, recStack, &cycle) {
			return fmt.Errorf(`cycle detected: %s`, strings.Join(cycle, "->"))
		}
	}

	return nil
}

func assertNotEmptyRuleSet(k string, rs []Rule) error {
	if len(rs) <= 0 {
		return fmt.Errorf(`node .%s must contain at least one rule`, k)
	}

	return nil
}

func assertNotConflictingTraversal(i int, k string, r Rule) error {
	if (r.Next != "" && r.Wait) ||
		(r.Next != "" && r.End) ||
		(r.End && r.Wait) ||
		(r.Next == "" && !r.Wait && !r.End) {
		return fmt.Errorf(`node .%s.rules[%d] must uniquely contain one of ("next", "wait", "and")`, k, i)
	}

	return nil
}

func assertValidTraversal(dag DAG, i int, k string, r Rule) error {
	_, ok := dag[r.Next]
	if !ok {
		return fmt.Errorf(`node .%s.rules[%d] traverses to invalid vertex %s`, k, i, r.Next)
	}

	return nil
}

func assertNoIsolatedNodes(dag DAG, visited map[string]struct{}) error {
	if len(visited) != len(dag) {
		var isolated []string
		for k := range dag {
			_, ok := visited[k]
			if !ok {
				isolated = append(isolated, fmt.Sprintf(".%s", k))
			}
		}
		return fmt.Errorf("DAG contains isolated nodes: [%s]", strings.Join(isolated, ","))
	}

	return nil
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

	// except "start", all vertices must have at least one incoming edge;
	// keep track of all vertices which have been visited
	visited := map[string]struct{}{"start": {}}
	for nk, nv := range dag {
		// all user-provided nodes must contain at least one rule
		err = assertNotEmptyRuleSet(nk, nv.Rules)
		if err != nil {
			return err
		}

		for i, rule := range nv.Rules {
			err = assertNotConflictingTraversal(i, nk, rule)
			if err != nil {
				return err
			}

			// all rules which specify a traversal must specify a valid vertex
			if rule.Next != "" {
				err = assertValidTraversal(dag, i, nk, rule)
				if err != nil {
					return err
				}

				visited[rule.Next] = struct{}{}
			}
		}
	}

	// DAG must not contain isolated nodes
	err = assertNoIsolatedNodes(dag, visited)
	if err != nil {
		return err
	}

	err = assertAcyclic(dag)
	if err != nil {
		return err
	}

	return nil
}