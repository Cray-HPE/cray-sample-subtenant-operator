/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	subtenantshpecomv1alpha1 "github.com/Cray-HPE/cray-sample-subtenant-operator/api/v1alpha1"
	tapmshpecomv1alpha1 "github.com/Cray-HPE/cray-tapms-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// SubTenantReconciler reconciles a SubTenant object
type SubTenantReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=subtenants.hpe.com,resources=subtenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=subtenants.hpe.com,resources=subtenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=subtenants.hpe.com,resources=subtenants/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SubTenant object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SubTenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	subTenantRequest := true
	tenantRequest := true
	log := r.Log.WithValues("subtenants", req.NamespacedName)
	subTenant := &subtenantshpecomv1alpha1.SubTenant{}
	err := r.Get(ctx, req.NamespacedName, subTenant)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("SubTenant resource not part of reconcile request (or deleted)")
			subTenantRequest = false
		}
	}

	if subTenantRequest {
		//
		// The CRD for this operator has changed -- implement reconcile logic
		// for this controller.
		//
		fmt.Printf("Reconciling subtenant: %+v\n", subTenant.Spec.TenantName)
	}

	tenant := &tapmshpecomv1alpha1.Tenant{}
	err = r.Get(ctx, req.NamespacedName, tenant)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("Tenant resource not part of reconcile request (or deleted)")
			tenantRequest = false
		}
	}
	if tenantRequest {
		//
		// The CRD for the tapms operator has changed -- implement appropriate reconcile
		// logic for this controller.
		//
		fmt.Printf("Reconciling tenant: %+v\n", tenant.Spec.TenantName)
		fmt.Printf("Tenant State: %+v\n", tenant.Status.State)
		fmt.Printf("Tenant HSM Partition Name: %+v\n", tenant.Status.HsmPartitionName)
		fmt.Printf("Tenant HSM Group Label: %+v\n", tenant.Status.HsmGroupLabel)
		fmt.Printf("Tenant xnames: %+v\n", tenant.Status.Xnames)
		fmt.Printf("Tenant child namespaces: %+v\n", tenant.Status.ChildNamespaces)
	}

	if !subTenantRequest && !tenantRequest {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SubTenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&subtenantshpecomv1alpha1.SubTenant{}).
		//
		// This is the key to watching the tapms CRD:
		//
		Watches(&source.Kind{Type: &tapmshpecomv1alpha1.Tenant{}},
			&handler.EnqueueRequestForObject{}).
		Complete(r)
}
