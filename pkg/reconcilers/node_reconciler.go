package reconcilers

import (
	"context"

	"github.com/OpenCHAMI/node-service/pkg/clients"
	"github.com/OpenCHAMI/node-service/pkg/resources/node"
)

// NOTE: Do not define type NodeReconciler struct here.
// It is already defined in node_reconciler_generated.go.

// reconcileNode is the "Brain" of the Node resource.
// The generated Reconcile() wrapper calls this function.
func (r *NodeReconciler) reconcileNode(ctx context.Context, res *node.Node) error {
	// Initialize the client locally since we can't add fields to the struct
	// In production, we might use a singleton pattern or environment variables
	serviceClients := clients.NewServiceClients("http://localhost:8081", "http://localhost:8082")

	// 1. Fetch External Data
	// We use the xname to look up data in the other services
	bootInfo, err := serviceClients.FetchBootConfig(res.Spec.XName)
	if err != nil {
		// Returning an error here triggers the generated code to:
		// 1. Log the error
		// 2. Set Status "Ready" = "False"
		// 3. Requeue after 30 seconds
		return err
	}

	metaInfo, err := serviceClients.FetchMetaConfig(res.Spec.XName)
	if err != nil {
		return err
	}

	// 2. Update Status (The Composed View)
	changed := false

	// Update Boot Status
	if res.Status.Boot.Kernel != bootInfo.Kernel {
		res.Status.Boot.Kernel = bootInfo.Kernel
		res.Status.Boot.Params = bootInfo.Params
		res.Status.Boot.Image = bootInfo.Image
		changed = true
	}

	// Update Config Status
	if len(res.Status.Config.Groups) != len(metaInfo.Groups) {
		res.Status.Config.Groups = metaInfo.Groups
		changed = true
	}

	// 3. Save Logic
	if changed {
		if res.Status.EffectiveProfile == "" {
			res.Status.EffectiveProfile = "default"
		}
		// We don't need to call Save() manually. 
		// The generated wrapper detects that we modified 'res' and will save it for us
		// if we return nil.
	}

	return nil
}