package k8s

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var CLAIM_NAME = "s3-db-backup"

func int32Ptr(i int32) *int32 { return &i }

func CreateDeployment(ctx context.Context, clientset *kubernetes.Clientset, dbType, connectionAddress, nameSpace string) error {
	// Define PVC
	pvcObjectCreationError := CreatePVC(ctx, clientset, CLAIM_NAME, nameSpace)
	if pvcObjectCreationError != nil {
		panic(pvcObjectCreationError.Error())
	}

	var backUpCommand string
	switch dbType {
	case strings.ToUpper("mysql"):
		backUpCommand = fmt.Sprintf("mysqldump %s | gzip > /backup/$(date +%Y-%m-%d-%T).sql.gz", connectionAddress)
	case strings.ToUpper("postgresql"):
	case strings.ToUpper("postgres"):
		backUpCommand = fmt.Sprintf("pgdump %s | gzip > /backup/$(date +%Y-%m-%d-%T).sql.gz", connectionAddress)
	}

	// Define MySQL deployment
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "backup-agent",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "backup-agent",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "backup-agent",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "backup-agent",
							Image: "mysql:latest",
							Args: []string{
								"sh",
								"-c",
								backUpCommand,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "backup-folder",
									MountPath: "/backup",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						// {
						// 	Name:       "backup",
						// 	Type:       corev1.SecretVolumeSource,
						// 	SecretName: "s3-secret",
						// },
						// *s3Volume,
						{
							Name: "backup-folder",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: CLAIM_NAME,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := clientset.AppsV1().Deployments(nameSpace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Create Secret
	// secretData := map[string]string{
	// 	"access_key": "YOUR_S3_ACCESS_KEY",
	// 	"secret_key": "YOUR_S3_SECRET_KEY",
	// }
	// secretCreationError := CreateSecret(ctx, clientset, "s3-secret", nameSpace, secretData)
	// if secretCreationError != nil {
	// 	panic(secretCreationError.Error())
	// }

	// // Check if the deployment already exists, if not create a new one
	// found := &appsv1.Deployment{}
	// err = r.Get(ctx, req.NamespacedName, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	// Define a new deployment
	// 	dep := r.deploymentFordbBackup(dbBackup)
	// 	log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 	err = r.Create(ctx, dep)
	// 	if err != nil {
	// 		log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Deployment created successfully - return and requeue
	// 	return ctrl.Result{Requeue: true}, nil
	// } else if err != nil {
	// 	log.Error(err, "Failed to get Deployment")
	// 	return ctrl.Result{}, err
	// }

	// // Ensure the deployment size is the same as the spec
	// size := dbBackup.Spec.DbUrl
	// if *found.Spec.Replicas != size {
	// 	found.Spec.Replicas = &size
	// 	err = r.Update(ctx, found)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Spec updated - return and requeue
	// 	return ctrl.Result{Requeue: true}, nil
	// }

	// // Update the dbBackup status with the pod names
	// // List the pods for this dbBackup's deployment
	// podList := &corev1.PodList{}
	// listOpts := []client.ListOption{
	// 	client.InNamespace(dbBackup.Namespace),
	// 	client.MatchingLabels(labelsFordbBackup(dbBackup.Name)),
	// }
	// if err = r.List(ctx, podList, listOpts...); err != nil {
	// 	log.Error(err, "Failed to list pods", "dbBackup.Namespace", dbBackup.Namespace, "dbBackup.Name", dbBackup.Name)
	// 	return ctrl.Result{}, err
	// }
	// podNames := getPodNames(podList.Items)

	// // Update status.Nodes if needed
	// if !reflect.DeepEqual(podNames, dbBackup.Status.Nodes) {
	// 	dbBackup.Status.Nodes = podNames
	// 	err := r.Status().Update(ctx, dbBackup)
	// 	if err != nil {
	// 		log.Error(err, "Failed to update dbBackup status")
	// 		return ctrl.Result{}, err
	// 	}
	// }
	return nil
}
