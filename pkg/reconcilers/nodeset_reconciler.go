package reconcilers

import (
	"context"
	"fmt"
	"regexp"

	"github.com/OpenCHAMI/node-service/pkg/resources/node"
	"github.com/OpenCHAMI/node-service/pkg/resources/nodeset"
)

// NOTE: Do not define type NodeSetReconciler struct here.
// It is already defined in nodeset_reconciler_generated.go.

// reconcileNodeSet calculates which nodes belong to this set
func (r *NodeSetReconciler) reconcileNodeSet(ctx context.Context, res *nodeset.NodeSet) error {
	// DEBUG LOG: Prove the reconciler is running
	fmt.Printf(">>> RECONCILING NODESET: %s\n", res.Metadata.Name)

	// 1. Fetch ALL Nodes from storage
	// We use the client embedded in the BaseReconciler
	allNodes := []node.Node{}
	if err := r.Client.List(ctx, "Node", &allNodes); err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
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

	// B. Label Selectors (e.g., "role": "compute")
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

	// C. Regex Pattern (e.g., "x1000.*")
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