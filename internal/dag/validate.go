package dag

import (
	"errors"
	"fmt"
	"strings"
)

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

// validateAcyclic uses a depth-first traversal to check that
// a DAG does not contain cycles.
func validateAcyclic(dag DAG) error {
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

// validateNoIsolated iterates over all outgoing edges in a DAG.
// It checks that all nodes in a DAG, besides the "start" node,
// have at least one incoming edge.
func validateNoIsolated(dag DAG) error {
	vertices := map[string]struct{}{"start": {}}
	for _, node := range dag {
		for _, rule := range node.Rules {
			if rule.Next != "" {
				vertices[rule.Next] = struct{}{}
			}
		}
	}

	if len(vertices) != len(dag) {
		isolated := []string{}
		for k := range dag {
			_, ok := vertices[k]
			if !ok {
				isolated = append(isolated, k)
			}
		}
		return fmt.Errorf("DAG contains isolated nodes: [%s]", strings.Join(isolated, ","))
	}

	return nil
}

// validateRulesNotEmpty checks that the given Node contains
// at least one Rule.
func validateRulesNotEmpty(k string, n Node) error {
	if len(n.Rules) <= 0 {
		return fmt.Errorf(`node .%s must contain at least one rule`, k)
	}

	return nil
}

// validateValidTraversal checks the following:
//   - A Rule cannot specify a conflicting traversal; that is, a Rule
//     must specify _exactly one_ of Next, Wait, or End.
//   - A Rule cannot specify a traversal to a nonexistent node.
func validateValidTraversal(dag DAG, i int, k string, r Rule) error {
	if (r.Next != "" && r.Wait) ||
		(r.Next != "" && r.End) ||
		(r.End && r.Wait) ||
		(r.Next == "" && !r.Wait && !r.End) {
		return fmt.Errorf(`node .%s.rules[%d] must uniquely contain one of ("next", "wait", "and")`, k, i)
	}

	if r.Next == "" {
		return nil
	}

	_, ok := dag[r.Next]
	if !ok {
		return fmt.Errorf(`node .%s.rules[%d] traverses to invalid vertex %s`, k, i, r.Next)
	}

	return nil
}

// validateDAGNotEmpty checks that a given DAG is not an empty map.
func validateDAGNotEmpty(dag DAG) error {
	if len(dag) <= 0 {
		return errors.New(`DAG must contain at least one node`)
	}

	return nil
}

// validateDAGNotEmpty checks that a given DAG contains a node called "start".
func validateContainsStartNode(dag DAG) error {
	_, ok := dag["start"]
	if !ok {
		return errors.New(`DAG must contain a node called "start"`)
	}

	return nil
}

func validate(dag DAG) error {
	if err := validateDAGNotEmpty(dag); err != nil {
		return err
	}
	if err := validateContainsStartNode(dag); err != nil {
		return err
	}

	for nk, nv := range dag {
		if err := validateRulesNotEmpty(nk, nv); err != nil {
			return err
		}

		for i, rule := range nv.Rules {
			if err := validateValidTraversal(dag, i, nk, rule); err != nil {
				return err
			}
		}
	}

	if err := validateAcyclic(dag); err != nil {
		return err
	}
	if err := validateNoIsolated(dag); err != nil {
		return err
	}

	return nil
}