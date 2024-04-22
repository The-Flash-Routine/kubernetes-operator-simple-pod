/*
Copyright 2024.

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

package controller

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	podv1 "routine.kat/simple-pod-operator/api/v1"
)

// SimplePodReconciler reconciles a SimplePod object
type SimplePodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pod.routine.kat,resources=simplepods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pod.routine.kat,resources=simplepods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pod.routine.kat,resources=simplepods/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SimplePod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *SimplePodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting Reconciliation loop at: " + time.Now().String())
	// Declare SimplePod Custom Object
	var simplePod podv1.SimplePod

	// Load SimplePod Custom Object
	if err := r.Get(ctx, req.NamespacedName, &simplePod); err != nil {
		log.Error(err, "Unable to fetch SimplePod object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch list of all pods
	var podList corev1.PodList

	if err := r.List(ctx, &podList, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "Unable to list Pods")
		return ctrl.Result{}, err
	}

	// This variable will denote any managed pods
	var foundPod *corev1.Pod
	// Iterating over existing podList to find any managed pod
	for _, pod := range podList.Items {
		labelMap := pod.ObjectMeta.Labels
		resourceOwner, found := labelMap["resourceOwner"]
		if found && resourceOwner == simplePod.ObjectMeta.Name {
			foundPod = &pod
			// Always update SimplePod Custom Object with POd's IP to keep ip for latest managed pod
			simplePod.Status.PodIp = foundPod.Status.PodIP
			if errUpdate := r.Status().Update(ctx, &simplePod); errUpdate != nil {
				log.Error(errUpdate, "Cannot Update Status of SimplePod Object")
				return ctrl.Result{}, errUpdate
			}
			log.Info("Updated status for SimplePod custom object: " + simplePod.Name)

			break
		}
	}

	// If no managed pod is found then create pod
	if foundPod == nil {
		log.Info("Found no managed Pod (and thus creating) for SimplePod custom object: " + simplePod.Name)
		// Construct Pod Object
		toCreatePod, err := createPod(&simplePod)
		if err != nil {
			log.Error(err, "Cannot construct Pod Object")
			return ctrl.Result{}, err
		}
		// Create Pod Object
		if errPodCreate := r.Create(ctx, toCreatePod); errPodCreate != nil {
			log.Error(errPodCreate, "Cannot create pod")
			return ctrl.Result{}, errPodCreate
		}
		log.Info("Created Pod for object: " + simplePod.Name)

	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func createPod(simplePod *podv1.SimplePod) (*corev1.Pod, error) {

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      simplePod.Name,
			Namespace: simplePod.Namespace,
			Labels:    simplePod.Labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{},
		},
	}

	pod.ObjectMeta.Labels["resourceOwner"] = simplePod.Name

	for _, container := range simplePod.Spec.Containers {
		pod.Spec.Containers = append(pod.Spec.Containers, container)
	}

	return pod, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimplePodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podv1.SimplePod{}).
		Complete(r)
}
