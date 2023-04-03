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
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	novav1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func CreateNovaWith3CellsAndEnsureReady(namespace string) NovaNames {
	var novaName types.NamespacedName
	var novaNames NovaNames
	var cell0 CellNames
	var cell1 CellNames
	var cell2 CellNames

	novaName = types.NamespacedName{
		Namespace: namespace,
		Name:      uuid.New().String(),
	}
	novaNames = GetNovaNames(novaName, []string{"cell0", "cell1", "cell2"})
	cell0 = novaNames.Cells["cell0"]
	cell1 = novaNames.Cells["cell1"]
	cell2 = novaNames.Cells["cell2"]

	DeferCleanup(k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
	DeferCleanup(
		k8sClient.Delete,
		ctx,
		CreateNovaMessageBusSecret(namespace, "mq-for-api-secret"),
	)
	DeferCleanup(
		k8sClient.Delete,
		ctx,
		CreateNovaMessageBusSecret(namespace, "mq-for-cell1-secret"),
	)
	DeferCleanup(
		k8sClient.Delete,
		ctx,
		CreateNovaMessageBusSecret(namespace, "mq-for-cell2-secret"),
	)

	serviceSpec := corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 3306}}}
	DeferCleanup(th.DeleteDBService, th.CreateDBService(namespace, "db-for-api", serviceSpec))
	DeferCleanup(th.DeleteDBService, th.CreateDBService(namespace, "db-for-cell1", serviceSpec))
	DeferCleanup(th.DeleteDBService, th.CreateDBService(namespace, "db-for-cell2", serviceSpec))

	spec := GetDefaultNovaSpec()
	cell0Template := GetDefaultNovaCellTemplate()
	cell0Template["cellName"] = "cell0"
	cell0Template["cellDatabaseInstance"] = "db-for-api"
	cell0Template["cellDatabaseUser"] = "nova_cell0"

	cell1Template := GetDefaultNovaCellTemplate()
	cell1Template["cellName"] = "cell1"
	cell1Template["cellDatabaseInstance"] = "db-for-cell1"
	cell1Template["cellDatabaseUser"] = "nova_cell1"
	cell1Template["cellMessageBusInstance"] = "mq-for-cell1"

	cell2Template := GetDefaultNovaCellTemplate()
	cell2Template["cellName"] = "cell2"
	cell2Template["cellDatabaseInstance"] = "db-for-cell2"
	cell2Template["cellDatabaseUser"] = "nova_cell2"
	cell2Template["cellMessageBusInstance"] = "mq-for-cell2"
	cell2Template["hasAPIAccess"] = false

	spec["cellTemplates"] = map[string]interface{}{
		"cell0": cell0Template,
		"cell1": cell1Template,
		"cell2": cell2Template,
	}
	spec["apiDatabaseInstance"] = "db-for-api"
	spec["apiMessageBusInstance"] = "mq-for-api"

	DeferCleanup(DeleteInstance, CreateNova(novaName, spec))
	keystoneAPIName := th.CreateKeystoneAPI(namespace)
	DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
	keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
	keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
	}, timeout, interval).Should(Succeed())

	th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)

	th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
	th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
	th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
	th.SimulateMariaDBDatabaseCompleted(cell2.MariaDBDatabaseName)

	th.SimulateTransportURLReady(cell0.TransportURLName)
	th.SimulateTransportURLReady(cell1.TransportURLName)
	th.SimulateTransportURLReady(cell2.TransportURLName)

	th.SimulateJobSuccess(cell0.CellDBSyncJobName)
	th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)

	th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
	th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)

	th.SimulateJobSuccess(cell1.CellDBSyncJobName)
	th.SimulateStatefulSetReplicaReady(cell1.ConductorStatefulSetName)

	th.SimulateJobSuccess(cell2.CellDBSyncJobName)
	th.SimulateStatefulSetReplicaReady(cell2.ConductorStatefulSetName)
	th.SimulateStatefulSetReplicaReady(novaSchedulerStatefulSetName)
	th.SimulateStatefulSetReplicaReady(novaMetadataStatefulSetName)
	th.ExpectCondition(
		novaName,
		ConditionGetterFunc(NovaConditionGetter),
		novav1.NovaAllCellsReadyCondition,
		corev1.ConditionTrue,
	)
	th.ExpectCondition(
		novaName,
		ConditionGetterFunc(NovaConditionGetter),
		condition.ReadyCondition,
		corev1.ConditionTrue,
	)
	return novaNames
}

var _ = Describe("Nova reconfiguration", func() {
	var novaNames NovaNames

	BeforeEach(func() {
		// Uncomment this if you need the full output in the logs from gomega
		// matchers
		// format.MaxLength = 0

		novaNames = CreateNovaWith3CellsAndEnsureReady(namespace)

	})
	When("cell0 conductor replicas is set to 0", func() {
		It("sets the deployment replicas to 0", func() {
			cell0DeploymentName := novaNames.Cells["cell0"].ConductorStatefulSetName

			deployment := th.GetStatefulSet(cell0DeploymentName)
			one := int32(1)
			Expect(deployment.Spec.Replicas).To(Equal(&one))

			// We need this big Eventually block because the Update() call might
			// return a Conflict and then we have to retry by re-reading Nova,
			// and updating the Replicas again.
			Eventually(func(g Gomega) {
				nova := GetNova(novaNames.NovaName)

				// TODO(gibi): Is there a simpler way to achieve this update
				// in golang?
				cell0 := nova.Spec.CellTemplates["cell0"]
				(&cell0).ConductorServiceTemplate.Replicas = int32(0)
				nova.Spec.CellTemplates["cell0"] = cell0

				err := k8sClient.Update(ctx, nova)
				g.Expect(err == nil || k8s_errors.IsConflict(err)).To(BeTrue())

				deployment = &appsv1.StatefulSet{}
				g.Expect(k8sClient.Get(ctx, cell0DeploymentName, deployment)).Should(Succeed())
				zero := int32(0)
				g.Expect(deployment.Spec.Replicas).To(Equal(&zero))
			}, timeout, interval).Should(Succeed())
		})
	})
	When("networkAttachemnt is added to a conductor while the definition is missing", func() {
		It("applys new NetworkAttachments configuration to that Conductor", func() {
			cell1Names := NewCell(novaName, "cell1")

			Eventually(func(g Gomega) {
				nova := GetNova(novaName)

				cell1 := nova.Spec.CellTemplates["cell1"]
				attachments := cell1.ConductorServiceTemplate.NetworkAttachments
				attachments = append(attachments, "internalapi")
				(&cell1).ConductorServiceTemplate.NetworkAttachments = attachments
				nova.Spec.CellTemplates["cell1"] = cell1

				err := k8sClient.Update(ctx, nova)
				g.Expect(err == nil || k8s_errors.IsConflict(err)).To(BeTrue())
			}, timeout, interval).Should(Succeed())

			th.ExpectConditionWithDetails(
				cell1Names.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)
			th.ExpectConditionWithDetails(
				cell1Names.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)

			th.ExpectConditionWithDetails(
				cell1Names.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)
			th.ExpectConditionWithDetails(
				cell1Names.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)

			th.ExpectConditionWithDetails(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NovaCell cell1 is not Ready",
			)
			th.ExpectConditionWithDetails(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NovaCell cell1 is not Ready",
			)

			internalAPINADName := types.NamespacedName{Namespace: namespace, Name: "internalapi"}
			DeferCleanup(DeleteInstance, CreateNetworkAttachmentDefinition(internalAPINADName))

			th.ExpectConditionWithDetails(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NovaCell cell1 is not Ready",
			)

			SimulateStatefulSetReplicaReadyWithPods(
				cell1Names.ConductorStatefulSetName,
				map[string][]string{namespace + "/internalapi": {"10.0.0.1"}},
			)

			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})

})
