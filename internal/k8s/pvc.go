package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	DRIVER_NAME    = "csi-s3"
	STORAGE_VOLUME = "1Gi"
)

func CreatePVC(ctx context.Context, clientset *kubernetes.Clientset, pvcName, nameSpace string) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: nameSpace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse(STORAGE_VOLUME),
				},
			},
			StorageClassName: &DRIVER_NAME,
		},
	}

	option := metav1.CreateOptions{}
	_, err := clientset.CoreV1().PersistentVolumeClaims(nameSpace).Create(ctx, pvc, option)
	if err != nil {
		return err
	}

	fmt.Println("PersistentVolumeClaim created successfully.")
	return nil
}
