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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	novav1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Nova multicell", func() {
	When("Nova CR instance is created with 3 cells", func() {
		BeforeEach(func() {
			// TODO(bogdando): deduplicate this into CreateNovaWith3CellsAndEnsureReady()
			DeferCleanup(k8sClient.Delete, ctx, CreateNovaSecret(novaNames.NovaName.Namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(cell0.CellName.Namespace, fmt.Sprintf("%s-secret", cell0.TransportURLName.Name)),
			)
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(cell1.CellName.Namespace, fmt.Sprintf("%s-secret", cell1.TransportURLName.Name)),
			)
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(cell2.CellName.Namespace, fmt.Sprintf("%s-secret", cell2.TransportURLName.Name)),
			)

			serviceSpec := corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 3306}}}
			DeferCleanup(th.DeleteDBService, th.CreateDBService(novaNames.APIMariaDBDatabaseName.Namespace, novaNames.APIMariaDBDatabaseName.Name, serviceSpec))
			DeferCleanup(th.DeleteDBService, th.CreateDBService(cell0.MariaDBDatabaseName.Namespace, cell0.MariaDBDatabaseName.Name, serviceSpec))
			DeferCleanup(th.DeleteDBService, th.CreateDBService(cell1.MariaDBDatabaseName.Namespace, cell1.MariaDBDatabaseName.Name, serviceSpec))
			DeferCleanup(th.DeleteDBService, th.CreateDBService(cell2.MariaDBDatabaseName.Namespace, cell2.MariaDBDatabaseName.Name, serviceSpec))

			spec := GetDefaultNovaSpec()
			cell0Template := GetDefaultNovaCellTemplate()
			cell0Template["cellName"] = cell0.CellName.Name
			cell0Template["cellDatabaseInstance"] = cell0.MariaDBDatabaseName.Name
			cell0Template["cellDatabaseUser"] = "nova_cell0"

			cell1Template := GetDefaultNovaCellTemplate()
			cell1Template["cellName"] = cell1.CellName.Name
			cell1Template["cellDatabaseInstance"] = cell1.MariaDBDatabaseName.Name
			cell1Template["cellDatabaseUser"] = "nova_cell1"
			cell1Template["cellMessageBusInstance"] = cell1.TransportURLName.Name

			cell2Template := GetDefaultNovaCellTemplate()
			cell2Template["cellName"] = cell2.CellName.Name
			cell2Template["cellDatabaseInstance"] = cell2.MariaDBDatabaseName.Name
			cell2Template["cellDatabaseUser"] = "nova_cell2"
			cell2Template["cellMessageBusInstance"] = cell2.TransportURLName.Name
			cell2Template["hasAPIAccess"] = false

			spec["cellTemplates"] = map[string]interface{}{
				"cell0": cell0Template,
				"cell1": cell1Template,
				"cell2": cell2Template,
			}
			spec["apiDatabaseInstance"] = novaNames.APIMariaDBDatabaseName.Name
			spec["apiMessageBusInstance"] = cell0.TransportURLName.Name

			DeferCleanup(DeleteInstance, CreateNova(novaNames.NovaName, spec))

			keystoneAPIName := th.CreateKeystoneAPI(novaNames.NovaName.Namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
		})

		It("creates cell0 NovaCell", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)

			// assert that cell related CRs are created pointing to the API MQ
			cell := GetNovaCell(cell0.CellName)
			Expect(cell.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell0.TransportURLName.Name)))
			conductor := GetNovaConductor(cell0.CellConductorName)
			Expect(conductor.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell0.TransportURLName.Name)))

			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionFalse,
			)
			// assert that cell0 using the same DB as the API
			dbSync := th.GetJob(cell0.CellDBSyncJobName)
			dbSyncJobEnv := dbSync.Spec.Template.Spec.InitContainers[0].Env
			Expect(dbSyncJobEnv).To(
				ContainElements(
					[]corev1.EnvVar{
						{Name: "CellDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", cell0.MariaDBDatabaseName.Name)},
						{Name: "APIDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", novaNames.APIMariaDBDatabaseName.Name)},
					},
				),
			)

			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionTrue,
			)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
			th.ExpectCondition(
				cell0.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
			)
		})

		It("creates NovaAPI", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)

			api := GetNovaAPI(novaNames.APIName)
			Expect(api.Spec.Replicas).Should(BeEquivalentTo(1))
			Expect(api.Spec.Cell0DatabaseHostname).To(Equal(fmt.Sprintf("hostname-for-%s", cell0.MariaDBDatabaseName.Name)))
			Expect(api.Spec.Cell0DatabaseHostname).To(Equal(api.Spec.APIDatabaseHostname))
			Expect(api.Spec.APIMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell0.TransportURLName.Name)))

			configDataMap := th.GetConfigMap(
				types.NamespacedName{
					Namespace: novaNames.APIName.Namespace,
					Name:      fmt.Sprintf("%s-config-data", novaNames.APIName.Name),
				},
			)
			Expect(configDataMap.Data).Should(HaveKey("01-nova.conf"))
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring(
					fmt.Sprintf("[database]\nconnection = mysql+pymysql://nova_cell0:12345678@hostname-for-%s/nova_cell0",
						cell0.MariaDBDatabaseName.Name)),
			)
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring(
					fmt.Sprintf("[api_database]\nconnection = mysql+pymysql://nova_api:12345678@hostname-for-%s/nova_api",
						novaNames.APIMariaDBDatabaseName.Name)),
			)

			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)

			th.ExpectCondition(
				novaNames.APIName,
				ConditionGetterFunc(NovaAPIConditionGetter),
				condition.DeploymentReadyCondition,
				corev1.ConditionTrue,
			)

			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("creates all cell DBs", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)

			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsDBReadyCondition,
				corev1.ConditionFalse,
			)
			th.SimulateMariaDBDatabaseCompleted(cell2.MariaDBDatabaseName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsDBReadyCondition,
				corev1.ConditionTrue,
			)

		})

		It("creates all cell MQs", func() {
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateTransportURLReady(cell1.TransportURLName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsMQReadyCondition,
				corev1.ConditionFalse,
			)
			th.SimulateTransportURLReady(cell2.TransportURLName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsMQReadyCondition,
				corev1.ConditionTrue,
			)

		})

		It("creates cell1 NovaCell", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)
			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell1.TransportURLName)

			// assert that cell related CRs are created pointing to the cell1 MQ
			c1 := GetNovaCell(cell1.CellName)
			Expect(c1.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell1.TransportURLName.Name)))
			c1Conductor := GetNovaConductor(cell1.CellConductorName)
			Expect(c1Conductor.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell1.TransportURLName.Name)))

			th.ExpectCondition(
				cell1.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionFalse,
			)
			// assert that cell1 using its own DB but has access to the API DB
			dbSync := th.GetJob(cell1.CellDBSyncJobName)
			dbSyncJobEnv := dbSync.Spec.Template.Spec.InitContainers[0].Env
			Expect(dbSyncJobEnv).To(
				ContainElements(
					[]corev1.EnvVar{
						{Name: "CellDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", cell1.MariaDBDatabaseName.Name)},
						{Name: "APIDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", novaNames.APIMariaDBDatabaseName.Name)},
					},
				),
			)
			th.SimulateJobSuccess(cell1.CellDBSyncJobName)
			th.ExpectCondition(
				cell1.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionTrue,
			)
			th.SimulateStatefulSetReplicaReady(cell1.ConductorStatefulSetName)
			th.ExpectCondition(
				cell1.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
			)
		})
		It("creates cell2 NovaCell", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)
			th.SimulateStatefulSetReplicaReady(novaNames.SchedulerStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)
			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell1.TransportURLName)
			th.SimulateJobSuccess(cell1.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell1.ConductorStatefulSetName)

			th.SimulateMariaDBDatabaseCompleted(cell2.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell2.TransportURLName)

			// assert that cell related CRs are created pointing to the Cell 2 MQ
			c2 := GetNovaCell(cell2.CellName)
			Expect(c2.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell2.TransportURLName.Name)))
			c2Conductor := GetNovaConductor(cell2.CellConductorName)
			Expect(c2Conductor.Spec.CellMessageBusSecretName).To(Equal(fmt.Sprintf("%s-secret", cell2.TransportURLName.Name)))
			th.ExpectCondition(
				cell2.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionFalse,
			)
			// assert that cell2 using its own DB but has *no* access to the API DB
			dbSync := th.GetJob(cell2.CellDBSyncJobName)
			dbSyncJobEnv := dbSync.Spec.Template.Spec.InitContainers[0].Env
			Expect(dbSyncJobEnv).To(
				ContainElements(
					[]corev1.EnvVar{
						{Name: "CellDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", cell2.MariaDBDatabaseName.Name)},
					},
				),
			)
			Expect(dbSyncJobEnv).NotTo(
				ContainElements(
					[]corev1.EnvVar{
						{Name: "APIDatabaseHost", Value: fmt.Sprintf("hostname-for-%s", novaNames.APIMariaDBDatabaseName.Name)},
					},
				),
			)
			th.SimulateJobSuccess(cell2.CellDBSyncJobName)
			th.ExpectCondition(
				cell2.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionTrue,
			)
			th.SimulateStatefulSetReplicaReady(cell2.ConductorStatefulSetName)
			th.ExpectCondition(
				cell2.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			// As cell2 was the last cell to deploy all cells is ready now and
			// Nova becomes ready
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
		It("creates cell2 NovaCell even if everthing else fails", func() {
			// Don't simulate any success for any other DBs MQs or Cells
			// just for cell2
			th.SimulateMariaDBDatabaseCompleted(cell2.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell2.TransportURLName)

			// assert that cell related CRs are created
			GetNovaCell(cell2.CellName)
			GetNovaConductor(cell2.CellConductorName)

			th.SimulateJobSuccess(cell2.CellDBSyncJobName)
			th.ExpectCondition(
				cell2.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionTrue,
			)
			th.SimulateStatefulSetReplicaReady(cell2.ConductorStatefulSetName)
			th.ExpectCondition(
				cell2.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				cell2.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
			// Only cell2 succeeded so Nova is not ready yet
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
			)
		})
		It("creates Nova API even if cell1 and cell2 fails", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)

			// Simulate that cell1 DB sync failed and do not simulate
			// cell2 DB creation success so that will be in Creating state.
			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell1.TransportURLName)
			th.SimulateJobFailure(cell1.CellDBSyncJobName)

			// NovaAPI is still created
			GetNovaAPI(novaNames.APIName)
			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
			)
		})
		It("does not create cell1 if cell0 fails as cell1 needs API access", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateTransportURLReady(cell1.TransportURLName)

			th.SimulateJobFailure(cell0.CellDBSyncJobName)

			NovaCellNotExists(cell1.CellName)
		})
	})

	When("Nova CR instance is created with collapsed cell deployment", func() {
		BeforeEach(func() {
			DeferCleanup(k8sClient.Delete, ctx, CreateNovaSecret(novaNames.NovaName.Namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(cell0.CellName.Namespace, fmt.Sprintf("%s-secret", cell0.TransportURLName.Name)),
			)

			serviceSpec := corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 3306}}}
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(novaNames.APIMariaDBDatabaseName.Namespace, novaNames.APIMariaDBDatabaseName.Name, serviceSpec))

			spec := GetDefaultNovaSpec()
			cell0Template := GetDefaultNovaCellTemplate()
			cell0Template["cellName"] = cell0.CellName.Name
			cell0Template["cellDatabaseInstance"] = cell0.MariaDBDatabaseName.Name
			cell0Template["cellDatabaseUser"] = "nova_cell0"
			cell0Template["hasAPIAccess"] = true
			// disable cell0 conductor
			cell0Template["conductorServiceTemplate"] = map[string]interface{}{
				"replicas": 0,
			}

			cell1Template := GetDefaultNovaCellTemplate()
			// cell1 is configured to have API access and use the same
			// message bus as the top level services. Hence cell1 conductor
			// will act both as a super conductor and as cell1 conductor
			cell1Template["cellName"] = cell1.CellName.Name
			cell1Template["cellDatabaseInstance"] = cell1.MariaDBDatabaseName.Name
			cell1Template["cellDatabaseUser"] = "nova_cell1"
			cell1Template["cellMessageBusInstance"] = cell1.TransportURLName.Name
			cell1Template["hasAPIAccess"] = true

			spec["cellTemplates"] = map[string]interface{}{
				"cell0": cell0Template,
				"cell1": cell1Template,
			}
			spec["apiDatabaseInstance"] = novaNames.APIMariaDBDatabaseName.Name
			spec["apiMessageBusInstance"] = cell0.TransportURLName.Name

			DeferCleanup(DeleteInstance, CreateNova(novaNames.NovaName, spec))
			keystoneAPIName := th.CreateKeystoneAPI(novaNames.NovaName.Namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
		})
		It("cell0 becomes ready with 0 conductor replicas and the rest of nova is deployed", func() {
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell1.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateTransportURLReady(cell1.TransportURLName)

			// We requested 0 replicas from the cell0 conductor so the
			// conductor is ready even if 0 replicas is running but all
			// the necessary steps, i.e. db-sync is run successfully
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)

			th.ExpectCondition(
				cell0.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			ss := th.GetStatefulSet(cell0.ConductorStatefulSetName)
			Expect(ss.Status.Replicas).To(Equal(int32(0)))
			Expect(ss.Status.AvailableReplicas).To(Equal(int32(0)))

			// As cell0 is ready cell1 is deployed
			th.SimulateJobSuccess(cell1.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell1.ConductorStatefulSetName)

			th.ExpectCondition(
				cell1.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)

			th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)
			// As cell0 is ready API is deployed
			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaNames.APIKeystoneEndpointName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIReadyCondition,
				corev1.ConditionTrue,
			)

			// As cell0 is ready scheduler is deployed
			th.SimulateStatefulSetReplicaReady(novaNames.SchedulerStatefulSetName)
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaSchedulerReadyCondition,
				corev1.ConditionTrue,
			)

			// So the whole Nova deployment is ready
			th.ExpectCondition(
				novaNames.NovaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})
})
