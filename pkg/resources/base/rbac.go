/*
Copyright 2019 Banzai Cloud.

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

package base

import (
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/banzaicloud/istio-operator/pkg/resources/templates"
)

func (r *Reconciler) serviceAccountReader() runtime.Object {
	return &apiv1.ServiceAccount{
		ObjectMeta: templates.ObjectMetaWithRevision(istioReaderServiceAccountName, istioReaderLabel, r.Config),
	}
}

func (r *Reconciler) clusterRoleReader() runtime.Object {
	return &rbacv1.ClusterRole{
		ObjectMeta: templates.ObjectMetaClusterScopeWithRevision(istioReaderName, istioReaderLabel, r.Config),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"config.istio.io", "security.istio.io", "networking.istio.io", "authentication.istio.io"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints", "pods", "services", "nodes", "replicationcontrollers", "namespaces"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"replicasets"},
				Verbs:     []string{"get", "watch", "list"},
			},
		},
	}
}

func (r *Reconciler) clusterRoleBindingReader() runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMetaClusterScopeWithRevision(istioReaderName, istioReaderLabel, r.Config),
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     r.Config.WithNamespacedRevision(istioReaderName),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      r.Config.WithRevision(istioReaderServiceAccountName),
				Namespace: r.Config.Namespace,
			},
		},
	}
}
