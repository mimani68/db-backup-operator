package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func int32Ptr(i int32) *int32 { return &i }

func CreateDeployment() {
	// Load Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", "/path/to/kubeconfig")
	if err != nil {
		panic(err.Error())
	}

	//Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Define S3 volume
	s3Volume := &corev1.Volume{
		Name: "s3-volume",
		VolumeSource: corev1.VolumeSource{
			AWSS3: &corev1.AWSS3VolumeSource{
				Bucket: "my-s3-bucket",
				Prefix: "path/to/folder",
				SecretRef: &corev1.SecretReference{
					Name: "s3-secret",
				},
			},
		},
	}

	// Define MySQL deployment
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mysql",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mysql",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:latest",
							Args: []string{
								"sh",
								"-c",
								"mysqldump -u root -p123 user | gzip > /backup/$(date +%Y-%m-%d-%T).sql.gz",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "backup",
									MountPath: "/backup",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name:       "backup",
							Type:       corev1.SecretVolumeSource,
							SecretName: "s3-secret",
						},
						{
							*s3Volume,
						},
					},
				},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments("default").Create(deployment)
	if err != nil {
		panic(err.Error())
	}

	// Define S3 secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "s3-secret",
		},
		StringData: map[string]string{
			"access_key": "YOUR_S3_ACCESS_KEY",
			"secret_key": "YOUR_S3_SECRET_KEY",
		},
		Type: corev1.SecretTypeOpaque,
	}
	_, err = clientset.CoreV1().Secrets("default").Create(secret)
	if err != nil {
		panic(err.Error())
	}

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
}
