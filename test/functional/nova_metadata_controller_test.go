/*
Copyright 2023.

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
	"encoding/json"
	"fmt"

	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openstack-k8s-operators/lib-common/modules/test/helpers"

	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
)

var _ = Describe("NovaMetadata controller", func() {
	When("with standard spec without networkinterface", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaMessageBusSecret(novaNames.MetadataName.Namespace, MessageBusSecretName))

			spec := GetDefaultNovaMetadataSpec()
			spec["customServiceConfig"] = "foo=bar"
			metadata := CreateNovaMetadata(novaNames.MetadataName, spec)
			DeferCleanup(DeleteInstance, metadata)
		})
		When("a NovaMetadata CR is created pointing to a non existent Secret", func() {

			It("is not Ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ReadyCondition, corev1.ConditionFalse,
				)
			})

			It("has empty Status fields", func() {
				instance := GetNovaMetadata(novaNames.MetadataName)
				// NOTE(gibi): Hash has `omitempty` tags so while
				// they are initialized to an empty map that value is omited from
				// the output when sent to the client. So we see nils here.
				Expect(instance.Status.Hash).To(BeEmpty())
				Expect(instance.Status.ReadyCount).To(Equal(int32(0)))
			})
			It("is missing the secret", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.InputReadyCondition,
					corev1.ConditionFalse,
				)
			})
		})

		When("an unrealated Secret is created the CR state does not change", func() {
			BeforeEach(func() {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "not-relevant-secret",
						Namespace: novaNames.MetadataName.Namespace,
					},
				}
				Expect(k8sClient.Create(ctx, secret)).Should(Succeed())
				DeferCleanup(k8sClient.Delete, ctx, secret)
			})

			It("is not Ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ReadyCondition,
					corev1.ConditionFalse,
				)
			})

			It("is missing the secret", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.InputReadyCondition,
					corev1.ConditionFalse,
				)
			})
		})

		When("the Secret is created but some fields are missing", func() {
			BeforeEach(func() {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      SecretName,
						Namespace: novaNames.MetadataName.Namespace,
					},
					Data: map[string][]byte{
						"NovaPassword": []byte("12345678"),
					},
				}
				Expect(k8sClient.Create(ctx, secret)).Should(Succeed())
				DeferCleanup(k8sClient.Delete, ctx, secret)
			})

			It("is not Ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ReadyCondition,
					corev1.ConditionFalse,
				)
			})

			It("reports that the inputes are not ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.InputReadyCondition,
					corev1.ConditionFalse,
				)
			})
		})

		When("the Secret is created with all the expected fields", func() {
			BeforeEach(func() {
				DeferCleanup(
					k8sClient.Delete,
					ctx,
					CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName),
				)
			})

			It("reports that input is ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.InputReadyCondition,
					corev1.ConditionTrue,
				)
			})
			It("generated configs successfully", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ServiceConfigReadyCondition,
					corev1.ConditionTrue,
				)

				configDataMap := th.GetConfigMap(
					types.NamespacedName{
						Namespace: novaNames.MetadataName.Namespace,
						Name:      fmt.Sprintf("%s-config-data", novaNames.MetadataName.Name),
					},
				)
				Expect(configDataMap).ShouldNot(BeNil())
				Expect(configDataMap.Data).Should(HaveKey("01-nova.conf"))
				Expect(configDataMap.Data).Should(
					HaveKeyWithValue("01-nova.conf",
						ContainSubstring("transport_url=rabbit://fake")))
				Expect(configDataMap.Data).Should(
					HaveKeyWithValue("02-nova-override.conf", "foo=bar"))
			})

			It("stored the input hash in the Status", func() {
				Eventually(func(g Gomega) {
					novaMetadata := GetNovaMetadata(novaNames.MetadataName)
					g.Expect(novaMetadata.Status.Hash).Should(HaveKeyWithValue("input", Not(BeEmpty())))
				}, timeout, interval).Should(Succeed())

			})

			When("the NovaMetadata is deleted", func() {
				It("deletes the generated ConfigMaps", func() {
					th.ExpectCondition(
						novaNames.MetadataName,
						ConditionGetterFunc(NovaMetadataConditionGetter),
						condition.ServiceConfigReadyCondition,
						corev1.ConditionTrue,
					)

					DeleteInstance(GetNovaMetadata(novaNames.MetadataName))

					Eventually(func() []corev1.ConfigMap {
						return th.ListConfigMaps(novaNames.MetadataName.Name).Items
					}, timeout, interval).Should(BeEmpty())
				})
			})
		})

		When("NovaMetadata is created with a proper Secret", func() {
			BeforeEach(func() {
				DeferCleanup(
					k8sClient.Delete, ctx, CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName))
			})

			It(" reports input ready", func() {
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.InputReadyCondition,
					corev1.ConditionTrue,
				)
			})

			It("creates a StatefulSet for the nova-metadata service", func() {
				th.ExpectConditionWithDetails(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.DeploymentReadyCondition,
					corev1.ConditionFalse,
					condition.RequestedReason,
					condition.DeploymentReadyRunningMessage,
				)

				ss := th.GetStatefulSet(novaNames.MetadataStatefulSetName)
				Expect(int(*ss.Spec.Replicas)).To(Equal(1))
				Expect(ss.Spec.Template.Spec.Volumes).To(HaveLen(2))
				Expect(ss.Spec.Template.Spec.Containers).To(HaveLen(2))
				Expect(ss.Spec.Selector.MatchLabels).To(Equal(map[string]string{"service": "nova-metadata"}))

				container := ss.Spec.Template.Spec.Containers[0]
				Expect(container.VolumeMounts).To(HaveLen(2))
				Expect(container.Image).To(Equal(ContainerImage))

				Expect(container.LivenessProbe.HTTPGet.Port.IntVal).To(Equal(int32(8775)))
				Expect(container.ReadinessProbe.HTTPGet.Port.IntVal).To(Equal(int32(8775)))

			})

			When("the StatefulSet has at least one Replica ready", func() {
				BeforeEach(func() {
					th.ExpectConditionWithDetails(
						novaNames.MetadataName,
						ConditionGetterFunc(NovaMetadataConditionGetter),
						condition.DeploymentReadyCondition,
						corev1.ConditionFalse,
						condition.RequestedReason,
						condition.DeploymentReadyRunningMessage,
					)
					th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)
				})

				It("reports that the StatefulSet is ready", func() {
					th.GetStatefulSet(novaNames.MetadataStatefulSetName)
					th.ExpectCondition(
						novaNames.MetadataName,
						ConditionGetterFunc(NovaMetadataConditionGetter),
						condition.DeploymentReadyCondition,
						corev1.ConditionTrue,
					)

					novaMetadata := GetNovaMetadata(novaNames.MetadataName)
					Expect(novaMetadata.Status.ReadyCount).To(BeNumerically(">", 0))
				})
			})

			It("exposes the service", func() {
				th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)
				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ExposeServiceReadyCondition,
					corev1.ConditionTrue,
				)
				service := GetService(types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "nova-metadata-internal"})
				Expect(service.Labels["service"]).To(Equal("nova-metadata"))
			})

			It("is Ready", func() {
				th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)

				th.ExpectCondition(
					novaNames.MetadataName,
					ConditionGetterFunc(NovaMetadataConditionGetter),
					condition.ReadyCondition,
					corev1.ConditionTrue,
				)
			})
		})
	})
	When("NovaMetadata is created with networkAttachments", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName))

			spec := GetDefaultNovaMetadataSpec()
			spec["networkAttachments"] = []string{"internalapi"}
			metadata := CreateNovaMetadata(novaNames.MetadataName, spec)
			DeferCleanup(DeleteInstance, metadata)
		})

		It("reports that the definition is missing", func() {
			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)
			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
			)
		})
		It("reports that network attachment is missing", func() {
			internalMetadataName := types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "internalapi"}
			nad := CreateNetworkAttachmentDefinition(internalMetadataName)
			DeferCleanup(DeleteInstance, nad)

			ss := th.GetStatefulSet(novaNames.MetadataStatefulSetName)

			expectedAnnotation, err := json.Marshal(
				[]networkv1.NetworkSelectionElement{
					{
						Name:      "internalapi",
						Namespace: novaNames.MetadataName.Namespace,
					}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ss.Spec.Template.ObjectMeta.Annotations).To(
				HaveKeyWithValue("k8s.v1.cni.cncf.io/networks", string(expectedAnnotation)),
			)

			// We don't add network attachment status annotations to the Pods
			// to simulate that the network attachments are missing.
			SimulateStatefulSetReplicaReadyWithPods(novaNames.MetadataStatefulSetName, map[string][]string{})

			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NetworkAttachments error occured "+
					"not all pods have interfaces with ips as configured in NetworkAttachments: [internalapi]",
			)
		})
		It("reports that an IP is missing", func() {
			internalMetadataName := types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "internalapi"}
			nad := CreateNetworkAttachmentDefinition(internalMetadataName)
			DeferCleanup(DeleteInstance, nad)

			ss := th.GetStatefulSet(novaNames.MetadataStatefulSetName)

			expectedAnnotation, err := json.Marshal(
				[]networkv1.NetworkSelectionElement{
					{
						Name:      "internalapi",
						Namespace: novaNames.MetadataName.Namespace,
					}})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(ss.Spec.Template.ObjectMeta.Annotations).To(
				HaveKeyWithValue("k8s.v1.cni.cncf.io/networks", string(expectedAnnotation)),
			)

			// We simulat that there is no IP associated with the internalapi
			// network attachment
			SimulateStatefulSetReplicaReadyWithPods(
				novaNames.MetadataStatefulSetName,
				map[string][]string{novaNames.MetadataName.Namespace + "/internalapi": {}},
			)

			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NetworkAttachments error occured "+
					"not all pods have interfaces with ips as configured in NetworkAttachments: [internalapi]",
			)
		})
		It("reports NetworkAttachmentsReady if the Pods got the proper annotiations", func() {
			internalMetadataName := types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "internalapi"}
			nad := CreateNetworkAttachmentDefinition(internalMetadataName)
			DeferCleanup(DeleteInstance, nad)

			SimulateStatefulSetReplicaReadyWithPods(
				novaNames.MetadataStatefulSetName,
				map[string][]string{novaNames.MetadataName.Namespace + "/internalapi": {"10.0.0.1"}},
			)

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionTrue,
			)

			Eventually(func(g Gomega) {
				instance := GetNovaMetadata(novaNames.MetadataName)
				g.Expect(instance.Status.NetworkAttachments).To(
					Equal(map[string][]string{novaNames.MetadataName.Namespace + "/internalapi": {"10.0.0.1"}}))

			}, timeout, interval).Should(Succeed())

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})
	When("NovaMetadata is created with externalEndpoints", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName))

			spec := GetDefaultNovaMetadataSpec()
			var externalEndpoints []interface{}
			externalEndpoints = append(
				externalEndpoints, map[string]interface{}{
					"endpoint":        "internal",
					"ipAddressPool":   "osp-internalapi",
					"loadBalancerIPs": []string{"internal-lb-ip-1", "internal-lb-ip-2"},
				},
			)
			spec["externalEndpoints"] = externalEndpoints

			metadata := CreateNovaMetadata(novaNames.MetadataName, spec)
			DeferCleanup(DeleteInstance, metadata)
		})

		It("creates MetalLB service", func() {
			th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)

			// As the internal enpoint is configured in ExternalEndpoints it does not
			// get a Route but a Service with MetalLB annotations instead
			service := GetService(types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "nova-metadata-internal"})
			Expect(service.Annotations).To(
				HaveKeyWithValue("metallb.universe.tf/address-pool", "osp-internalapi"))
			Expect(service.Annotations).To(
				HaveKeyWithValue("metallb.universe.tf/allow-shared-ip", "osp-internalapi"))
			Expect(service.Annotations).To(
				HaveKeyWithValue("metallb.universe.tf/loadBalancerIPs", "internal-lb-ip-1,internal-lb-ip-2"))
			AssertRouteNotExists(types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "nova-metadata-internal"})

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})
	When("NovaMetadata is reconfigured", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName))

			metadata := CreateNovaMetadata(novaNames.MetadataName, GetDefaultNovaMetadataSpec())
			DeferCleanup(DeleteInstance, metadata)

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ServiceConfigReadyCondition,
				corev1.ConditionTrue,
			)

			th.SimulateStatefulSetReplicaReady(novaNames.MetadataStatefulSetName)
			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})

		It("applys new NetworkAttachments configuration", func() {
			Eventually(func(g Gomega) {
				novaMetadata := GetNovaMetadata(novaNames.MetadataName)
				novaMetadata.Spec.NetworkAttachments = append(novaMetadata.Spec.NetworkAttachments, "internalapi")

				err := k8sClient.Update(ctx, novaMetadata)
				g.Expect(err == nil || k8s_errors.IsConflict(err)).To(BeTrue())
			}, timeout, interval).Should(Succeed())

			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)
			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.RequestedReason,
				"NetworkAttachment resources missing: internalapi",
			)

			internalAPINADName := types.NamespacedName{Namespace: novaNames.MetadataName.Namespace, Name: "internalapi"}
			DeferCleanup(DeleteInstance, CreateNetworkAttachmentDefinition(internalAPINADName))

			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NetworkAttachments error occured "+
					"not all pods have interfaces with ips as configured in NetworkAttachments: [internalapi]",
			)
			th.ExpectConditionWithDetails(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionFalse,
				condition.ErrorReason,
				"NetworkAttachments error occured "+
					"not all pods have interfaces with ips as configured in NetworkAttachments: [internalapi]",
			)

			SimulateStatefulSetReplicaReadyWithPods(
				novaNames.MetadataStatefulSetName,
				map[string][]string{novaNames.MetadataName.Namespace + "/internalapi": {"10.0.0.1"}},
			)

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.NetworkAttachmentsReadyCondition,
				corev1.ConditionTrue,
			)

			Eventually(func(g Gomega) {
				novaMetadata := GetNovaMetadata(novaNames.MetadataName)
				g.Expect(novaMetadata.Status.NetworkAttachments).To(
					Equal(map[string][]string{novaNames.MetadataName.Namespace + "/internalapi": {"10.0.0.1"}}))

			}, timeout, interval).Should(Succeed())

			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.ReadyCondition,
				corev1.ConditionTrue,
			)
		})
	})

	When("starts zero replicas", func() {
		BeforeEach(func() {
			DeferCleanup(
				k8sClient.Delete, ctx, CreateNovaMetadataSecret(novaNames.MetadataName.Namespace, SecretName))

			spec := GetDefaultNovaMetadataSpec()
			spec["replicas"] = 0
			metadata := CreateNovaMetadata(novaNames.MetadataName, spec)
			DeferCleanup(DeleteInstance, metadata)
		})
		It("and deployment is Ready", func() {
			ss := th.GetStatefulSet(novaNames.MetadataStatefulSetName)
			Expect(int(*ss.Spec.Replicas)).To(Equal(0))
			th.ExpectCondition(
				novaNames.MetadataName,
				ConditionGetterFunc(NovaMetadataConditionGetter),
				condition.DeploymentReadyCondition,
				corev1.ConditionTrue,
			)

		})
	})
})
