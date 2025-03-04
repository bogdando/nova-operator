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

package v1beta1

import (
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NovaConductorTemplate defines the input parameters specified by the user to
// create a NovaConductor via higher level CRDs.
type NovaConductorTemplate struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="quay.io/tripleowallabycentos9/openstack-nova-conductor:current-tripleo"
	// The service specific Container Image URL
	ContainerImage string `json:"containerImage"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1
	// +kubebuilder:validation:Maximum=32
	// +kubebuilder:validation:Minimum=0
	// Replicas of the service to run
	Replicas int32 `json:"replicas"`

	// +kubebuilder:validation:Optional
	// NodeSelector to target subset of worker nodes running this service
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="# add your customization here"
	// CustomServiceConfig - customize the service config using this parameter to change service defaults,
	// or overwrite rendered information using raw OpenStack config format. The content gets added to
	// to /etc/<service>/<service>.conf.d directory as custom.conf file.
	CustomServiceConfig string `json:"customServiceConfig,omitempty"`

	// +kubebuilder:validation:Optional
	// ConfigOverwrite - interface to overwrite default config files like e.g. logging.conf
	// But can also be used to add additional files. Those get added to the service config dir in /etc/<service> .
	DefaultConfigOverwrite map[string]string `json:"defaultConfigOverwrite,omitempty"`

	// +kubebuilder:validation:Optional
	// Resources - Compute Resources required by this service (Limits/Requests).
	// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// NovaConductorSpec defines the desired state of NovaConductor
type NovaConductorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// CellName is the name of the Nova Cell this conductor belongs to.
	CellName string `json:"cellName,omitempty"`

	// +kubebuilder:validation:Required
	// Secret is the name of the Secret instance containing password
	// information for the nova-conductor service.
	Secret string `json:"secret,omitempty"`

	// +kubebuilder:validation:Optional
	// PasswordSelectors - Field names to identify the passwords from the
	// Secret
	PasswordSelectors PasswordSelector `json:"passwordSelectors,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nova
	// ServiceUser - optional username used for this service to register in
	// keystone
	ServiceUser string `json:"serviceUser"`

	// +kubebuilder:validation:Optional
	// KeystoneAuthURL - the URL that the nova-conductor service can use to
	// talk to keystone
	// NOTE(gibi): This is made optional here to allow reusing the
	// NovaConductorSpec struct in the Nova CR via the NovaCell CR where this
	// information is not yet known. We could make this required via multiple
	// options:
	// a) create a NovaConductorTemplate that duplicates NovaConductorSpec
	//    without this field. Use NovaCondcutorTemplate in NovaCellSpec.
	// b) do a) but pull out a the fields to a base struct that are used in
	//    both NovaConductorSpec and NovaCondcutorTemplate
	// c) add a validating webhook here that runs only when NovaConductor CR is
	//    created and does not run when Nova CR is created and make this field
	//    required via that webhook.
	KeystoneAuthURL string `json:"keystoneAuthURL"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nova
	// APIDatabaseUser - username to use when accessing the API DB
	APIDatabaseUser string `json:"apiDatabaseUser"`

	// +kubebuilder:validation:Optional
	// APIDatabaseHostname - hostname to use when accessing the API DB. If not
	// provided then upcalls will be disabled. This filed is Required for
	// cell0.
	// TODO(gibi): Add a webhook to validate cell0 constraint
	APIDatabaseHostname string `json:"apiDatabaseHostname"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nova
	// APIMessageBusUser - username to use when accessing the API message bus
	APIMessageBusUser string `json:"apiMessageBusUser"`

	// +kubebuilder:validation:Optional
	// APIMessageBusHostname - hostname to use when accessing the API message
	// bus. If not provided then upcalls will be disabled. This filed is
	// Required for cell0.
	// TODO(gibi): Add a webhook to validate cell0 constraint.
	APIMessageBusHostname string `json:"apiMessageBusHostname"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nova
	// CellDatabaseUser - username to use when accessing the cell DB
	CellDatabaseUser string `json:"cellDatabaseUser"`

	// +kubebuilder:validation:Optional
	// NOTE(gibi): This should be Required, see notes in KeystoneAuthURL
	// CellDatabaseHostname - hostname to use when accessing the cell DB
	CellDatabaseHostname string `json:"cellDatabaseHostname"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nova
	// CellMessageBusUser - username to use when accessing the cell message bus
	CellMessageBusUser string `json:"cellMessageBusUser"`

	// +kubebuilder:validation:Optional
	// NOTE(gibi): This should be Required, see notes in KeystoneAuthURL
	// CellMessageBusHostname - hostname to use when accessing the cell message
	// bus
	CellMessageBusHostname string `json:"cellMessageBusHostname"`

	// +kubebuilder:validation:Optional
	// Debug - enable debug for different deploy stages. If an init container
	// is used, it runs and the actual action pod gets started with sleep
	// infinity
	Debug Debug `json:"debug,omitempty"`

	// +kubebuilder:validation:Required
	// NovaServiceBase specifies the generic fields of the service
	NovaServiceBase `json:",inline"`
}

// NovaConductorStatus defines the observed state of NovaConductor
type NovaConductorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Map of hashes to track e.g. job status
	Hash map[string]string `json:"hash,omitempty"`

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`

	// ReadyCount defines the number of replicas ready from nova-conductor
	ReadyCount int32 `json:"readyCount,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NovaConductor is the Schema for the novaconductors API
type NovaConductor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NovaConductorSpec   `json:"spec,omitempty"`
	Status NovaConductorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NovaConductorList contains a list of NovaConductor
type NovaConductorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NovaConductor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NovaConductor{}, &NovaConductorList{})
}

// GetConditions returns the list of conditions from the status
func (s NovaConductorStatus) GetConditions() condition.Conditions {
	return s.Conditions
}

// NewNovaConductorSpec constructs a NovaConductorSpec
func NewNovaConductorSpec(
	novaCell NovaCellSpec,
) NovaConductorSpec {
	conductorSpec := NovaConductorSpec{
		CellName:             novaCell.CellName,
		Secret:               novaCell.Secret,
		CellDatabaseHostname: novaCell.CellDatabaseHostname,
		CellDatabaseUser:     novaCell.CellDatabaseUser,
		APIDatabaseHostname:  novaCell.APIDatabaseHostname,
		APIDatabaseUser:      novaCell.APIDatabaseUser,
		Debug:                novaCell.Debug,
		NovaServiceBase:      NovaServiceBase(novaCell.ConductorServiceTemplate),
	}
	return conductorSpec
}
