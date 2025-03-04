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
package functional_test

import (
	"os"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	novav1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
)

var _ = Describe("NovaCell controller", func() {
	var namespace string
	var novaCellName types.NamespacedName

	BeforeEach(func() {
		// NOTE(gibi): We need to create a unique namespace for each test run
		// as namespaces cannot be deleted in a locally running envtest. See
		// https://book.kubebuilder.io/reference/envtest.html#namespace-usage-limitation
		namespace = uuid.New().String()
		CreateNamespace(namespace)
		// We still request the delete of the Namespace to properly cleanup if
		// we run the test in an existing cluster.
		DeferCleanup(DeleteNamespace, namespace)
		// NOTE(gibi): ConfigMap generation looks up the local templates
		// directory via ENV, so provide it
		DeferCleanup(os.Setenv, "OPERATOR_TEMPLATES", os.Getenv("OPERATOR_TEMPLATES"))
		os.Setenv("OPERATOR_TEMPLATES", "../../templates")

		// Uncomment this if you need the full output in the logs from gomega
		// matchers
		// format.MaxLength = 0

	})

	When("A NovaCell CR instance is created without any input", func() {
		BeforeEach(func() {
			novaCellName = CreateNovaCell(namespace, novav1.NovaCellSpec{})
			DeferCleanup(DeleteNovaCell, novaCellName)
		})

		It("is not Ready", func() {
			ExpectCondition(
				novaCellName,
				conditionGetterFunc(NovaCellConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionUnknown,
			)
		})

		It("has no hash and no services ready", func() {
			instance := GetNovaCell(novaCellName)
			Expect(instance.Status.Hash).To(BeEmpty())
			Expect(instance.Status.ConductorServiceReadyCount).To(Equal(int32(0)))
			Expect(instance.Status.MetadataServiceReadyCount).To(Equal(int32(0)))
			Expect(instance.Status.NoVNCPRoxyServiceReadyCount).To(Equal(int32(0)))
		})
	})

	When("A NovaCell CR instance is created", func() {
		var novaConductorName types.NamespacedName

		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaConductorSecret(namespace, SecretName),
			)

			novaCellName = CreateNovaCell(
				namespace,
				novav1.NovaCellSpec{
					CellName: "cell0",
					Secret:   SecretName,
					ConductorServiceTemplate: novav1.NovaConductorTemplate{
						ContainerImage: ContainerImage,
						Replicas:       1,
					},
				},
			)
			DeferCleanup(DeleteNovaCell, novaCellName)
			novaConductorName = types.NamespacedName{
				Namespace: namespace,
				Name:      novaCellName.Name + "-conductor",
			}
		})

		It("creates the NovaConductor and tracks its readiness", func() {
			GetNovaConductor(novaConductorName)
			ExpectCondition(
				novaCellName,
				conditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionFalse,
			)
			novaCell := GetNovaCell(novaCellName)
			Expect(novaCell.Status.ConductorServiceReadyCount).To(Equal(int32(0)))
		})

		When("NovaConductor is ready", func() {
			var novaConductorDBSyncJobName types.NamespacedName

			BeforeEach(func() {
				ExpectCondition(
					novaConductorName,
					conditionGetterFunc(NovaConductorConditionGetter),
					condition.DBSyncReadyCondition,
					corev1.ConditionFalse,
				)
				novaConductorDBSyncJobName = types.NamespacedName{
					Namespace: namespace,
					Name:      novaConductorName.Name + "-db-sync",
				}
				SimulateJobSuccess(novaConductorDBSyncJobName)

				ExpectCondition(
					novaConductorName,
					conditionGetterFunc(NovaConductorConditionGetter),
					condition.DBSyncReadyCondition,
					corev1.ConditionTrue,
				)
			})

			It("reports that NovaConductor is ready", func() {
				ExpectCondition(
					novaCellName,
					conditionGetterFunc(NovaCellConditionGetter),
					novav1.NovaConductorReadyCondition,
					corev1.ConditionTrue,
				)
			})

			It("is Ready", func() {
				ExpectCondition(
					novaCellName,
					conditionGetterFunc(NovaCellConditionGetter),
					condition.ReadyCondition,
					corev1.ConditionTrue,
				)
			})
		})
	})
})
