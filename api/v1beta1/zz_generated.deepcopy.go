//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Debug) DeepCopyInto(out *Debug) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Debug.
func (in *Debug) DeepCopy() *Debug {
	if in == nil {
		return nil
	}
	out := new(Debug)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Nova) DeepCopyInto(out *Nova) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Nova.
func (in *Nova) DeepCopy() *Nova {
	if in == nil {
		return nil
	}
	out := new(Nova)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Nova) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaAPI) DeepCopyInto(out *NovaAPI) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaAPI.
func (in *NovaAPI) DeepCopy() *NovaAPI {
	if in == nil {
		return nil
	}
	out := new(NovaAPI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaAPI) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaAPIList) DeepCopyInto(out *NovaAPIList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaAPI, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaAPIList.
func (in *NovaAPIList) DeepCopy() *NovaAPIList {
	if in == nil {
		return nil
	}
	out := new(NovaAPIList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaAPIList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaAPISpec) DeepCopyInto(out *NovaAPISpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.NovaServiceBase.DeepCopyInto(&out.NovaServiceBase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaAPISpec.
func (in *NovaAPISpec) DeepCopy() *NovaAPISpec {
	if in == nil {
		return nil
	}
	out := new(NovaAPISpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaAPIStatus) DeepCopyInto(out *NovaAPIStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.APIEndpoints != nil {
		in, out := &in.APIEndpoints, &out.APIEndpoints
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaAPIStatus.
func (in *NovaAPIStatus) DeepCopy() *NovaAPIStatus {
	if in == nil {
		return nil
	}
	out := new(NovaAPIStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaAPITemplate) DeepCopyInto(out *NovaAPITemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaAPITemplate.
func (in *NovaAPITemplate) DeepCopy() *NovaAPITemplate {
	if in == nil {
		return nil
	}
	out := new(NovaAPITemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaCell) DeepCopyInto(out *NovaCell) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaCell.
func (in *NovaCell) DeepCopy() *NovaCell {
	if in == nil {
		return nil
	}
	out := new(NovaCell)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaCell) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaCellList) DeepCopyInto(out *NovaCellList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaCell, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaCellList.
func (in *NovaCellList) DeepCopy() *NovaCellList {
	if in == nil {
		return nil
	}
	out := new(NovaCellList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaCellList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaCellSpec) DeepCopyInto(out *NovaCellSpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.ConductorServiceTemplate.DeepCopyInto(&out.ConductorServiceTemplate)
	in.MetadataServiceTemplate.DeepCopyInto(&out.MetadataServiceTemplate)
	in.NoVNCProxyServiceTemplate.DeepCopyInto(&out.NoVNCProxyServiceTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaCellSpec.
func (in *NovaCellSpec) DeepCopy() *NovaCellSpec {
	if in == nil {
		return nil
	}
	out := new(NovaCellSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaCellStatus) DeepCopyInto(out *NovaCellStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaCellStatus.
func (in *NovaCellStatus) DeepCopy() *NovaCellStatus {
	if in == nil {
		return nil
	}
	out := new(NovaCellStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaCellTemplate) DeepCopyInto(out *NovaCellTemplate) {
	*out = *in
	in.ConductorServiceTemplate.DeepCopyInto(&out.ConductorServiceTemplate)
	in.MetadataServiceTemplate.DeepCopyInto(&out.MetadataServiceTemplate)
	in.NoVNCProxyServiceTemplate.DeepCopyInto(&out.NoVNCProxyServiceTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaCellTemplate.
func (in *NovaCellTemplate) DeepCopy() *NovaCellTemplate {
	if in == nil {
		return nil
	}
	out := new(NovaCellTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaConductor) DeepCopyInto(out *NovaConductor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaConductor.
func (in *NovaConductor) DeepCopy() *NovaConductor {
	if in == nil {
		return nil
	}
	out := new(NovaConductor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaConductor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaConductorList) DeepCopyInto(out *NovaConductorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaConductor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaConductorList.
func (in *NovaConductorList) DeepCopy() *NovaConductorList {
	if in == nil {
		return nil
	}
	out := new(NovaConductorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaConductorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaConductorSpec) DeepCopyInto(out *NovaConductorSpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.NovaServiceBase.DeepCopyInto(&out.NovaServiceBase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaConductorSpec.
func (in *NovaConductorSpec) DeepCopy() *NovaConductorSpec {
	if in == nil {
		return nil
	}
	out := new(NovaConductorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaConductorStatus) DeepCopyInto(out *NovaConductorStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaConductorStatus.
func (in *NovaConductorStatus) DeepCopy() *NovaConductorStatus {
	if in == nil {
		return nil
	}
	out := new(NovaConductorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaConductorTemplate) DeepCopyInto(out *NovaConductorTemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaConductorTemplate.
func (in *NovaConductorTemplate) DeepCopy() *NovaConductorTemplate {
	if in == nil {
		return nil
	}
	out := new(NovaConductorTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaList) DeepCopyInto(out *NovaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Nova, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaList.
func (in *NovaList) DeepCopy() *NovaList {
	if in == nil {
		return nil
	}
	out := new(NovaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaMetadata) DeepCopyInto(out *NovaMetadata) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaMetadata.
func (in *NovaMetadata) DeepCopy() *NovaMetadata {
	if in == nil {
		return nil
	}
	out := new(NovaMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaMetadata) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaMetadataList) DeepCopyInto(out *NovaMetadataList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaMetadata, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaMetadataList.
func (in *NovaMetadataList) DeepCopy() *NovaMetadataList {
	if in == nil {
		return nil
	}
	out := new(NovaMetadataList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaMetadataList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaMetadataSpec) DeepCopyInto(out *NovaMetadataSpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.NovaServiceBase.DeepCopyInto(&out.NovaServiceBase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaMetadataSpec.
func (in *NovaMetadataSpec) DeepCopy() *NovaMetadataSpec {
	if in == nil {
		return nil
	}
	out := new(NovaMetadataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaMetadataStatus) DeepCopyInto(out *NovaMetadataStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaMetadataStatus.
func (in *NovaMetadataStatus) DeepCopy() *NovaMetadataStatus {
	if in == nil {
		return nil
	}
	out := new(NovaMetadataStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaMetadataTemplate) DeepCopyInto(out *NovaMetadataTemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaMetadataTemplate.
func (in *NovaMetadataTemplate) DeepCopy() *NovaMetadataTemplate {
	if in == nil {
		return nil
	}
	out := new(NovaMetadataTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaNoVNCProxy) DeepCopyInto(out *NovaNoVNCProxy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaNoVNCProxy.
func (in *NovaNoVNCProxy) DeepCopy() *NovaNoVNCProxy {
	if in == nil {
		return nil
	}
	out := new(NovaNoVNCProxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaNoVNCProxy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaNoVNCProxyList) DeepCopyInto(out *NovaNoVNCProxyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaNoVNCProxy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaNoVNCProxyList.
func (in *NovaNoVNCProxyList) DeepCopy() *NovaNoVNCProxyList {
	if in == nil {
		return nil
	}
	out := new(NovaNoVNCProxyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaNoVNCProxyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaNoVNCProxySpec) DeepCopyInto(out *NovaNoVNCProxySpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.NovaServiceBase.DeepCopyInto(&out.NovaServiceBase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaNoVNCProxySpec.
func (in *NovaNoVNCProxySpec) DeepCopy() *NovaNoVNCProxySpec {
	if in == nil {
		return nil
	}
	out := new(NovaNoVNCProxySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaNoVNCProxyStatus) DeepCopyInto(out *NovaNoVNCProxyStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaNoVNCProxyStatus.
func (in *NovaNoVNCProxyStatus) DeepCopy() *NovaNoVNCProxyStatus {
	if in == nil {
		return nil
	}
	out := new(NovaNoVNCProxyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaNoVNCProxyTemplate) DeepCopyInto(out *NovaNoVNCProxyTemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaNoVNCProxyTemplate.
func (in *NovaNoVNCProxyTemplate) DeepCopy() *NovaNoVNCProxyTemplate {
	if in == nil {
		return nil
	}
	out := new(NovaNoVNCProxyTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaScheduler) DeepCopyInto(out *NovaScheduler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaScheduler.
func (in *NovaScheduler) DeepCopy() *NovaScheduler {
	if in == nil {
		return nil
	}
	out := new(NovaScheduler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaScheduler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSchedulerList) DeepCopyInto(out *NovaSchedulerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NovaScheduler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSchedulerList.
func (in *NovaSchedulerList) DeepCopy() *NovaSchedulerList {
	if in == nil {
		return nil
	}
	out := new(NovaSchedulerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NovaSchedulerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSchedulerSpec) DeepCopyInto(out *NovaSchedulerSpec) {
	*out = *in
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.NovaServiceBase.DeepCopyInto(&out.NovaServiceBase)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSchedulerSpec.
func (in *NovaSchedulerSpec) DeepCopy() *NovaSchedulerSpec {
	if in == nil {
		return nil
	}
	out := new(NovaSchedulerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSchedulerStatus) DeepCopyInto(out *NovaSchedulerStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSchedulerStatus.
func (in *NovaSchedulerStatus) DeepCopy() *NovaSchedulerStatus {
	if in == nil {
		return nil
	}
	out := new(NovaSchedulerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSchedulerTemplate) DeepCopyInto(out *NovaSchedulerTemplate) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSchedulerTemplate.
func (in *NovaSchedulerTemplate) DeepCopy() *NovaSchedulerTemplate {
	if in == nil {
		return nil
	}
	out := new(NovaSchedulerTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaServiceBase) DeepCopyInto(out *NovaServiceBase) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DefaultConfigOverwrite != nil {
		in, out := &in.DefaultConfigOverwrite, &out.DefaultConfigOverwrite
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaServiceBase.
func (in *NovaServiceBase) DeepCopy() *NovaServiceBase {
	if in == nil {
		return nil
	}
	out := new(NovaServiceBase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaSpec) DeepCopyInto(out *NovaSpec) {
	*out = *in
	if in.CellTemplates != nil {
		in, out := &in.CellTemplates, &out.CellTemplates
		*out = make(map[string]NovaCellTemplate, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	out.PasswordSelectors = in.PasswordSelectors
	out.Debug = in.Debug
	in.APIServiceTemplate.DeepCopyInto(&out.APIServiceTemplate)
	in.SchedulerServiceTemplate.DeepCopyInto(&out.SchedulerServiceTemplate)
	in.MetadataServiceTemplate.DeepCopyInto(&out.MetadataServiceTemplate)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaSpec.
func (in *NovaSpec) DeepCopy() *NovaSpec {
	if in == nil {
		return nil
	}
	out := new(NovaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NovaStatus) DeepCopyInto(out *NovaStatus) {
	*out = *in
	if in.Hash != nil {
		in, out := &in.Hash, &out.Hash
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(condition.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NovaStatus.
func (in *NovaStatus) DeepCopy() *NovaStatus {
	if in == nil {
		return nil
	}
	out := new(NovaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PasswordSelector) DeepCopyInto(out *PasswordSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PasswordSelector.
func (in *PasswordSelector) DeepCopy() *PasswordSelector {
	if in == nil {
		return nil
	}
	out := new(PasswordSelector)
	in.DeepCopyInto(out)
	return out
}
