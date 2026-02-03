// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package nodeset

import (
	"context"
	"github.com/openchami/fabrica/pkg/resource"
)

// NodeSet represents a NodeSet resource
type NodeSet struct {
	resource.Resource
	Spec   NodeSetSpec   `json:"spec" validate:"required"`
	Status NodeSetStatus `json:"status,omitempty"`
}

// NodeSetSpec defines the desired state of NodeSet
type NodeSetSpec struct {
	// 1. Explicit Selection
    XNames []string `json:"xnames,omitempty"`
    
    // 2. Dynamic Selection (Labels)
    // e.g. {"role": "compute", "subRole": "worker"}
    Labels map[string]string `json:"labels,omitempty"`
    
    // 3. Regex Selection
    // e.g. "x1000.*"
    XNamePattern string `json:"xNamePattern,omitempty"`
}

// NodeSetStatus defines the observed state of NodeSet
type NodeSetStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	Ready      bool   `json:"ready"`
	// Add your status fields here
	// The calculated list of nodes that match the Spec
    ResolvedNodes []string `json:"resolvedNodes"`
    
    // How many nodes matched
    Count int `json:"count"`
    
    // Hash to track if we need to re-calculate (optimization)
    GenerationHash string  `json:"generationHash"`
}

// Validate implements custom validation logic for NodeSet
func (r *NodeSet) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Name == "forbidden" {
	//     return errors.New("name 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *NodeSet) GetKind() string {
	return "NodeSet"
}

// GetName returns the name of the resource
func (r *NodeSet) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *NodeSet) GetUID() string {
	return r.Metadata.UID
}

func init() {
	// Register resource type prefix for storage
	resource.RegisterResourcePrefix("NodeSet", "nod")
}
