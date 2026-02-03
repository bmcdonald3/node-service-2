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
	Description string `json:"description,omitempty" validate:"max=200"`
	// Add your spec fields here
}

// NodeSetStatus defines the observed state of NodeSet
type NodeSetStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	Ready      bool   `json:"ready"`
	// Add your status fields here
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
