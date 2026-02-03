package reconcilers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/OpenCHAMI/node-service/pkg/resources/node"
	"github.com/OpenCHAMI/node-service/pkg/resources/nodeset"
)

// reconcileNodeSet calculates which nodes belong to this set
func (r *NodeSetReconciler) reconcileNodeSet(ctx context.Context, res *nodeset.NodeSet) error {
	fmt.Printf(">>> RECONCILING NODESET: %s\n", res.Metadata.Name)

	// 1. Fetch ALL Nodes from storage
	// FIX: List returns ([]interface{}, error), not unmarshaling into a pointer
	items, err := r.Client.List(ctx, "Node")
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	// Manually cast the generic items to our concrete Node type
	var allNodes []*node.Node
	for _, item := range items {
		if n, ok := item.(*node.Node); ok {
			allNodes = append(allNodes, n)
		}
	}

	// 2. Filter Nodes based on Spec
	var matchedXNames []string

	// Helper to avoid duplicates
	isMatched := func(xname string) bool {
		for _, m := range matchedXNames {
			if m == xname {
				return true
			}
		}
		return false
	}

	// A. Explicit XNames
	if len(res.Spec.XNames) > 0 {
		matchedXNames = append(matchedXNames, res.Spec.XNames...)
	}

	// B. Label Selectors
	if len(res.Spec.Labels) > 0 {
		for _, n := range allNodes {
			if isMatched(n.Spec.XName) {
				continue
			}
			
			match := true
			for k, v := range res.Spec.Labels {
				// Check if node has the label and it matches
				if val, ok := n.Spec.Labels[k]; !ok || val != v {
					match = false
					break
				}
			}
			if match {
				matchedXNames = append(matchedXNames, n.Spec.XName)
			}
		}
	}

	// C. Regex Pattern
	if res.Spec.XNamePattern != "" {
		re, err := regexp.Compile(res.Spec.XNamePattern)
		if err == nil {
			for _, n := range allNodes {
				if !isMatched(n.Spec.XName) && re.MatchString(n.Spec.XName) {
					matchedXNames = append(matchedXNames, n.Spec.XName)
				}
			}
		}
	}

	// 3. Update Status
	// Only update if the count has changed to avoid infinite loops
	if res.Status.Count != len(matchedXNames) {
		res.Status.ResolvedNodes = matchedXNames
		res.Status.Count = len(matchedXNames)
		fmt.Printf("    -> Resolved %d nodes for %s\n", res.Status.Count, res.Metadata.Name)
		// Returning nil triggers the auto-save in the generated wrapper
	}
	
	return nil
}