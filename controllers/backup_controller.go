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
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	ops1alpha1 "github.com/mimani68/db-backup-operator/api/v1alpha1"
	"github.com/mimani68/db-backup-operator/internal/k8s"
	"github.com/mimani68/db-backup-operator/internal/utils"
)

// BackupReconciler reconciles a Backup object
type BackupReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ops.db.io,resources=backups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ops.db.io,resources=backups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ops.db.io,resources=backups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *BackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	dbBackup := &ops1alpha1.Backup{}
	err := r.Get(ctx, req.NamespacedName, dbBackup)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("dbBackup resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get dbBackup")
		return ctrl.Result{}, err
	} else {
		r.Log.Info("dbBackup resource found.")
		value, err := json.Marshal(dbBackup)
		if err == nil {
			r.Log.Info(string(value))
		}
	}

	dbType, err := utils.FindKeyword(dbBackup.Spec.Type)
	if err != nil {
		r.Log.Info("Database type was not valid.")
		return ctrl.Result{}, err
	}

	client, err := k8s.LoadConfig()
	var nameSpace = "default"

	deploymentError := k8s.CreateDeployment(ctx, client, dbType, dbBackup.Spec.DbConnectionUrl, nameSpace)
	if deploymentError != nil {
		r.Log.Info("Database type was not valid.")
		return ctrl.Result{}, deploymentError
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ops1alpha1.Backup{}).
		Complete(r)
}
