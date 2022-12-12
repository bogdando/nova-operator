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

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	novav1 "github.com/openstack-k8s-operators/nova-operator/api/v1beta1"
)

var _ = Describe("Nova controller", func() {
	var novaName types.NamespacedName
	var novaNames NovaNames
	var cell0 CellNames

	BeforeEach(func() {
		// Uncomment this if you need the full output in the logs from gomega
		// matchers
		// format.MaxLength = 0

		novaName = types.NamespacedName{
			Namespace: namespace,
			Name:      uuid.New().String(),
		}
		novaNames = GetNovaNames(novaName, []string{"cell0"})
		cell0 = novaNames.Cells["cell0"]
	})

	When("Nova CR instance is created without cell0", func() {
		BeforeEach(func() {

			DeferCleanup(DeleteInstance, CreateNovaWithoutCell0(novaName))
		})

		It("is not Ready", func() {
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
			)
		})

		It("has no hash and no services ready", func() {
			instance := GetNova(novaName)
			Expect(instance.Status.Hash).To(BeEmpty())
			Expect(instance.Status.APIServiceReadyCount).To(Equal(int32(0)))
			Expect(instance.Status.SchedulerServiceReadyCount).To(Equal(int32(0)))
			Expect(instance.Status.MetadataServiceReadyCount).To(Equal(int32(0)))
		})

		It("reports that cell0 is missing from the spec", func() {
			th.ExpectConditionWithDetails(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NovaCell creation failed for cell0(missing cell0 specification from Spec.CellTemplates)",
			)
		})
	})

	When("Nova CR instance is created", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(namespace, MessageBusSecretName),
			)
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"openstack",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			keystoneAPIName := th.CreateKeystoneAPI(namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			DeferCleanup(DeleteInstance, CreateNovaWithCell0(novaName))
		})

		It("registers nova service to keystone", func() {
			// assert that the KeystoneService for nova is created
			th.GetKeystoneService(novaNames.KeystoneServiceName)
			// and simulate that it becomes ready i.e. the keystone-operator
			// did its job and registered the nova service
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			keystone := th.GetKeystoneService(novaNames.KeystoneServiceName)
			Expect(keystone.Status.Conditions).ToNot(BeNil())

			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.KeystoneServiceReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("creates nova_api DB", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIDBReadyCondition,
				corev1.ConditionFalse,
			)
			th.GetMariaDBDatabase(novaNames.APIMariaDBDatabaseName)

			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIDBReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("creates nova-api MQ", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIMQReadyCondition,
				corev1.ConditionFalse,
			)
			th.GetTransportURL(cell0.TransportURLName)

			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIMQReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("creates nova_cell0 DB", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsDBReadyCondition,
				corev1.ConditionFalse,
			)
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.GetMariaDBDatabase(cell0.MariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsDBReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("creates cell0 NovaCell", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			// assert that cell related CRs are created
			cell := GetNovaCell(cell0.CellName)
			Expect(cell.Spec.CellMessageBusSecretName).To(Equal("rabbitmq-secret"))
			Expect(cell.Spec.ServiceUser).To(Equal("nova"))

			conductor := GetNovaConductor(cell0.CellConductorName)
			Expect(conductor.Spec.CellMessageBusSecretName).To(Equal("rabbitmq-secret"))
			Expect(conductor.Spec.ServiceUser).To(Equal("nova"))

			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionFalse,
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
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("create NovaAPI", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
                        th.SimulateStatefulSetReplicaReady(novaNames.SchedulerStatefulSetName)
                        th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)

			api := GetNovaAPI(novaNames.APIName)
			Expect(api.Spec.APIMessageBusSecretName).To(Equal("rabbitmq-secret"))
			Expect(api.Spec.ServiceUser).To(Equal("nova"))

			th.SimulateStatefulSetReplicaReady(novaNames.APIDeploymentName)
			th.SimulateKeystoneEndpointReady(novaAPIKeystoneEndpointName)
			th.ExpectCondition(
				novaNames.APIName,
				ConditionGetterFunc(NovaAPIConditionGetter),
				condition.DeploymentReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("create NovaScheduler", func() {
			th.SimulateKeystoneServiceReady(novaKeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForAPI)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForCell0)
			th.SimulateTransportURLReady(apiTransportURLName)
			th.SimulateJobSuccess(cell0DBSyncJobName)
			th.SimulateStatefulSetReplicaReady(novaCell0ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaAPIdeploymentName)
			th.SimulateKeystoneEndpointReady(novaAPIKeystoneEndpointName)
			th.SimulateStatefulSetReplicaReady(novaMetadataStatefulSetName)

			scheduler := GetNovaScheduler(novaSchedulerName)
			Expect(scheduler.Spec.APIMessageBusSecretName).To(Equal("rabbitmq-secret"))
			th.SimulateStatefulSetReplicaReady(novaSchedulerStatefulSetName)

			th.ExpectCondition(
				novaSchedulerName,
				ConditionGetterFunc(NovaSchedulerConditionGetter),
				condition.DeploymentReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaSchedulerReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
		It("create NovaMetadata", func() {
			th.SimulateKeystoneServiceReady(novaKeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForAPI)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForCell0)
			th.SimulateTransportURLReady(apiTransportURLName)
			th.SimulateJobSuccess(cell0DBSyncJobName)
			th.SimulateStatefulSetReplicaReady(novaCell0ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaAPIdeploymentName)
			th.SimulateKeystoneEndpointReady(novaAPIKeystoneEndpointName)
			th.SimulateStatefulSetReplicaReady(novaSchedulerStatefulSetName)

			metadata := GetNovaMetadata(novaMetadataName)
			Expect(metadata.Spec.APIMessageBusSecretName).To(Equal("rabbitmq-secret"))
			th.SimulateStatefulSetReplicaReady(novaMetadataStatefulSetName)

			th.ExpectCondition(
				novaMetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.DeploymentReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaMetadataReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})

	When("Nova CR instance is created but cell0 DB sync fails", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(namespace, MessageBusSecretName),
			)
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"openstack",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			keystoneAPIName := th.CreateKeystoneAPI(namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			DeferCleanup(DeleteInstance, CreateNovaWithCell0(novaName))

		It("does not create NovaAPI", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)
			GetNovaCell(cell0.CellName)
			GetNovaConductor(cell0.CellConductorName)

			th.SimulateJobFailure(cell0.CellDBSyncJobName)
			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionFalse,
			)
		})

		It("does not set the cell0 ready condtion to ture", func() {
			th.ExpectCondition(
				cell0.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionFalse,
			)
		})

		It("does not set the all cell ready condtion", func() {
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionFalse,
			)
		})

			NovaAPINotExists(novaNames.APIName)
		})

		It("does not create NovaScheduler", func() {
			NovaSchedulerNotExists(novaSchedulerName)
		})
	})

	When("Nova CR instance with different DB Services for nova_api and cell0 DBs", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(namespace, MessageBusSecretName),
			)

			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"db-for-cell0",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"db-for-api",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			keystoneAPIName := th.CreateKeystoneAPI(namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			spec := GetDefaultNovaSpec()
			cell0 := GetDefaultNovaCellTemplate()
			cell0["cellDatabaseInstance"] = "db-for-cell0"
			spec["cellTemplates"] = map[string]interface{}{"cell0": cell0}
			spec["apiDatabaseInstance"] = "db-for-api"

			DeferCleanup(DeleteInstance, CreateNova(novaName, spec))
		})

		It("uses the correct hostnames to access the different DB services", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(novaNames.APIMariaDBDatabaseName)
			th.SimulateMariaDBDatabaseCompleted(cell0.MariaDBDatabaseName)
			th.SimulateTransportURLReady(cell0.TransportURLName)

			cell0DBSync := th.GetJob(cell0.CellDBSyncJobName)
			cell0DBSyncJobEnv := cell0DBSync.Spec.Template.Spec.InitContainers[0].Env
			Expect(cell0DBSyncJobEnv).To(
				ContainElements(
					[]corev1.EnvVar{
						{Name: "CellDatabaseHost", Value: "hostname-for-db-for-cell0"},
						{Name: "APIDatabaseHost", Value: "hostname-for-db-for-api"},
					},
				),
			)

			th.SimulateJobSuccess(cell0.CellDBSyncJobName)
			th.SimulateStatefulSetReplicaReady(cell0.ConductorStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaSchedulerStatefulSetName)
			th.SimulateStatefulSetReplicaReady(novaMetadataStatefulSetName)

			configDataMap := th.GetConfigMap(
				types.NamespacedName{
					Namespace: namespace,
					Name:      fmt.Sprintf("%s-config-data", novaNames.APIName.Name),
				},
			)
			Expect(configDataMap.Data).Should(HaveKey("01-nova.conf"))
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring("[database]\nconnection = mysql+pymysql://nova_cell0:12345678@hostname-for-db-for-cell0/nova_cell0"),
			)
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring("[api_database]\nconnection = mysql+pymysql://nova_api:12345678@hostname-for-db-for-api/nova_api"),
			)

			configDataMap = th.GetConfigMap(
				types.NamespacedName{
					Namespace: namespace,
					Name:      fmt.Sprintf("%s-config-data", novaSchedulerName.Name),
				},
			)
			Expect(configDataMap.Data).Should(HaveKey("01-nova.conf"))
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring("[database]\nconnection = mysql+pymysql://nova_cell0:12345678@hostname-for-db-for-cell0/nova_cell0"),
			)
			Expect(configDataMap.Data["01-nova.conf"]).To(
				ContainSubstring("[api_database]\nconnection = mysql+pymysql://nova_api:12345678@hostname-for-db-for-api/nova_api"),
			)

			th.SimulateStatefulSetReplicaReady(novaNames.APIName)
			th.SimulateKeystoneEndpointReady(novaAPIKeystoneEndpointName)

			th.ExpectCondition(
				cell0.CellConductorName,
				ConditionGetterFunc(NovaConductorConditionGetter),
				condition.DBSyncReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				cell0.CellName,
				ConditionGetterFunc(NovaCellConditionGetter),
				novav1.NovaConductorReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAllCellsReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaAPIReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				novav1.NovaSchedulerReadyCondition,
				corev1.ConditionTrue,
			)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})

	When("Nova CR instance is deleted", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(namespace, MessageBusSecretName),
			)
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"openstack",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			keystoneAPIName := th.CreateKeystoneAPI(namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			DeferCleanup(DeleteInstance, CreateNovaWithCell0(novaName))
		})

		It("removes the finalizer from KeystoneService", func() {
			th.SimulateKeystoneServiceReady(novaNames.KeystoneServiceName)
			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.KeystoneServiceReadyCondition,
				corev1.ConditionTrue,
			)

			service := th.GetKeystoneService(novaNames.KeystoneServiceName)
			Expect(service.Finalizers).To(ContainElement("Nova"))

			DeleteNova(novaName)
			service = th.GetKeystoneService(novaNames.KeystoneServiceName)
			Expect(service.Finalizers).NotTo(ContainElement("Nova"))
		})

		It("removes the finalizers from the nova dbs", func() {
			th.SimulateKeystoneServiceReady(novaKeystoneServiceName)

			apiDB := th.GetMariaDBDatabase(mariaDBDatabaseNameForAPI)
			Expect(apiDB.Finalizers).To(ContainElement("Nova"))
			cell0DB := th.GetMariaDBDatabase(mariaDBDatabaseNameForCell0)
			Expect(cell0DB.Finalizers).To(ContainElement("Nova"))

			DeleteInstance(GetNova(novaName))

			apiDB = th.GetMariaDBDatabase(mariaDBDatabaseNameForAPI)
			Expect(apiDB.Finalizers).NotTo(ContainElement("Nova"))
			cell0DB = th.GetMariaDBDatabase(mariaDBDatabaseNameForCell0)
			Expect(cell0DB.Finalizers).NotTo(ContainElement("Nova"))
		})
	})
	When("Nova CR instance is created with NetworkAttachment and ExternalEndpoints", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaSecret(namespace, SecretName))
			DeferCleanup(
				k8sClient.Delete,
				ctx,
				CreateNovaMessageBusSecret(namespace, MessageBusSecretName),
			)
			DeferCleanup(
				th.DeleteDBService,
				th.CreateDBService(
					namespace,
					"openstack",
					corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Port: 3306}},
					},
				),
			)
			keystoneAPIName := th.CreateKeystoneAPI(namespace)
			DeferCleanup(th.DeleteKeystoneAPI, keystoneAPIName)
			keystoneAPI := th.GetKeystoneAPI(keystoneAPIName)
			keystoneAPI.Status.APIEndpoints["internal"] = "http://keystone-internal-openstack.testing"
			Eventually(func(g Gomega) {
				g.Expect(k8sClient.Status().Update(ctx, keystoneAPI.DeepCopy())).Should(Succeed())
			}, timeout, interval).Should(Succeed())

			internalAPINADName := types.NamespacedName{Namespace: namespace, Name: "internalapi"}
			nad := CreateNetworkAttachmentDefinition(internalAPINADName)
			DeferCleanup(DeleteInstance, nad)

			var externalEndpoints []interface{}
			externalEndpoints = append(
				externalEndpoints, map[string]interface{}{
					"endpoint":        "internal",
					"ipAddressPool":   "osp-internalapi",
					"loadBalancerIPs": []string{"10.1.0.1", "10.1.0.2"},
				},
			)
			rawSpec := map[string]interface{}{
				"secret": SecretName,
				"cellTemplates": map[string]interface{}{
					"cell0": map[string]interface{}{
						"cellDatabaseUser": "nova_cell0",
						"hasAPIAccess":     true,
						"conductorServiceTemplate": map[string]interface{}{
							"networkAttachments": []string{"internalapi"},
						},
					},
				},
				"apiServiceTemplate": map[string]interface{}{
					"networkAttachments": []string{"internalapi"},
					"externalEndpoints":  externalEndpoints,
				},
				"schedulerServiceTemplate": map[string]interface{}{
					"networkAttachments": []string{"internalapi"},
				},
			}
			DeferCleanup(DeleteInstance, CreateNova(novaName, rawSpec))

			th.SimulateKeystoneServiceReady(novaKeystoneServiceName)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForAPI)
			th.SimulateMariaDBDatabaseCompleted(mariaDBDatabaseNameForCell0)
			th.SimulateTransportURLReady(apiTransportURLName)
			th.SimulateJobSuccess(cell0DBSyncJobName)
		})

		It("creates all the sub CRs and passes down the network parameters", func() {
			SimulateStatefulSetReplicaReadyWithPods(
				novaCell0ConductorStatefulSetName,
				map[string][]string{namespace + "/internalapi": {"10.0.0.1"}},
			)
			SimulateStatefulSetReplicaReadyWithPods(
				novaSchedulerStatefulSetName,
				map[string][]string{namespace + "/internalapi": {"10.0.0.1"}},
			)
			SimulateStatefulSetReplicaReadyWithPods(
				novaAPIdeploymentName,
				map[string][]string{namespace + "/internalapi": {"10.0.0.1"}},
			)
			SimulateStatefulSetReplicaReadyWithPods(
				novaMetadataStatefulSetName,
				map[string][]string{namespace + "/internalapi": {"10.0.0.1"}},
			)
			th.SimulateKeystoneEndpointReady(novaAPIKeystoneEndpointName)

			th.ExpectCondition(
				novaName,
				ConditionGetterFunc(NovaConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)

			nova := GetNova(novaName)

			conductor := GetNovaConductor(cell0ConductorName)
			Expect(conductor.Spec.NetworkAttachments).To(
				Equal(nova.Spec.CellTemplates["cell0"].ConductorServiceTemplate.NetworkAttachments))

			api := GetNovaAPI(novaAPIName)
			Expect(api.Spec.NetworkAttachments).To(Equal(nova.Spec.APIServiceTemplate.NetworkAttachments))
			Expect(api.Spec.ExternalEndpoints).To(Equal(nova.Spec.APIServiceTemplate.ExternalEndpoints))

			scheduler := GetNovaScheduler(novaSchedulerName)
			Expect(scheduler.Spec.NetworkAttachments).To(Equal(nova.Spec.APIServiceTemplate.NetworkAttachments))
		})

	})
})
