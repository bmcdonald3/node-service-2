package reconcilers

import (
	"context"
	"time"

	"github.com/openchami/fabrica/pkg/reconcile"
	"github.com/OpenCHAMI/node-service/pkg/clients"
	"github.com/OpenCHAMI/node-service/pkg/resources/node"
)

type NodeReconciler struct {
	*NodeReconcilerGenerated // Embedding the generated boilerplate
    Clients *clients.ServiceClients
}

func (r *NodeReconciler) ReconcileNode(ctx context.Context, res *node.Node) (reconcile.Result, error) {
	if r.Clients == nil {
        // Hardcoded URLs for now, or fetch from env vars
        r.Clients = clients.NewServiceClients("http://localhost:8081", "http://localhost:8082")
    }
	// This function is the "Brain" of the Node resource.
    // It runs whenever a Node is created, updated, or the requeue timer expires.

	// 1. Fetch External Data
    // We use the xname to look up data in the other services
    bootInfo, err := r.Clients.FetchBootConfig(res.Spec.XName)
    if err != nil {
        // If we can't talk to boot-service, we log it but don't crash
        // We might want to requeue quickly to retry
        return reconcile.Result{RequeueAfter: 30 * time.Second}, err
    }

    metaInfo, err := r.Clients.FetchMetaConfig(res.Spec.XName)
    if err != nil {
        return reconcile.Result{RequeueAfter: 30 * time.Second}, err
    }

    // 2. Update Status (The Composed View)
    // We compare current status with new data to avoid unnecessary writes
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
    // (Note: deeper slice comparison would be better for production)

    // 3. Save if changed
    if changed {
        // Determine effective profile (Phase 4 placeholder)
        if res.Status.EffectiveProfile == "" {
            res.Status.EffectiveProfile = "default"
        }
        
        // The generated wrapper handles the actual .Save() call if we modify the object
        // but we explicitly return it here to be safe/clear
    }

	// 4. Periodic Sync
    // We want to poll frequently because the boot-service might change 
    // without us knowing.
	return reconcile.Result{RequeueAfter: 1 * time.Minute}, nil
}