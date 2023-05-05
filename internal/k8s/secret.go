package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateSecret(ctx context.Context, clientset *kubernetes.Clientset, secretName, nameSpace string, secretData map[string]string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Type: corev1.SecretTypeOpaque,
		// StringData: map[string]string{
		// 	"access_key": "YOUR_S3_ACCESS_KEY",
		// 	"secret_key": "YOUR_S3_SECRET_KEY",
		// },
		StringData: secretData,
	}
	_, err := clientset.CoreV1().Secrets(nameSpace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
}
