// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package node

import (
	"context"
	"github.com/openchami/fabrica/pkg/resource"
)

// Node represents a Node resource
type Node struct {
	resource.Resource
	Spec   NodeSpec   `json:"spec" validate:"required"`
	Status NodeStatus `json:"status,omitempty"`
}

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	// Primary Identity (from SMD)
    XName        string            `json:"xname" fabrica:"filterable,unique"` 
    Role         string            `json:"role" fabrica:"filterable"`
    SubRole      string            `json:"subRole,omitempty" fabrica:"filterable"`
    Architecture string            `json:"architecture,omitempty"`
    BmcAddress   string            `json:"bmcAddress,omitempty"`
    MacAddress   string            `json:"macAddress,omitempty"`
    
    // Inventory Labels (from SMD)
    Labels       map[string]string `json:"labels,omitempty" fabrica:"filterable"`
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	Ready      bool   `json:"ready"`
	// Add your status fields here
	// The "Effective" profile (Calculated by Node Service)
    EffectiveProfile string `json:"effectiveProfile"` 
    
    // Observed state (fetched from backend services)
    Boot struct {
        Kernel string `json:"kernel,omitempty"`
        Params string `json:"params,omitempty"`
        Image  string `json:"image,omitempty"`
    } `json:"boot,omitempty"`

    Config struct {
        Groups []string `json:"groups,omitempty"`
    } `json:"config,omitempty"`
}

// Validate implements custom validation logic for Node
func (r *Node) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Name == "forbidden" {
	//     return errors.New("name 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *Node) GetKind() string {
	return "Node"
}

// GetName returns the name of the resource
func (r *Node) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Node) GetUID() string {
	return r.Metadata.UID
}

func init() {
	// Register resource type prefix for storage
	resource.RegisterResourcePrefix("Node", "nod")
}
