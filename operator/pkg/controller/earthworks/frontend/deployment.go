package frontend

import (
	earthworksv1alpha1 "github.com/stephenhillier/geoprojects/operator/pkg/apis/earthworks/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewFrontendDeployment takes an Earthworks object and creates a frontend Deployment resource
// information contained in the Earthworks object includes the name, number of replicas,
// and the image version for the web frontend.
func NewFrontendDeployment(e *earthworksv1alpha1.Earthworks) *appsv1.Deployment {
	ls := labelsForFrontend(e.Name)
	replicas := e.Spec.Size
	image := e.Spec.Image

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      e.Name,
			Namespace: e.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "earthworks-web",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}

	return dep
}

func labelsForFrontend(name string) map[string]string {
	return map[string]string{"heritage": "earthworks", "app": name, "version": "v1", "svc": "frontend"}
}
