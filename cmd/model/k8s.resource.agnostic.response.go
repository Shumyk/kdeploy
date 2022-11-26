package model

import (
	apps "k8s.io/api/apps/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type K8sResourceAgnosticResponse struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec apps.DeploymentSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

func (in *K8sResourceAgnosticResponse) DeepCopy() *K8sResourceAgnosticResponse {
	if in == nil {
		return nil
	}
	out := new(K8sResourceAgnosticResponse)
	in.DeepCopyInto(out)
	return out
}

func (in *K8sResourceAgnosticResponse) DeepCopyInto(out *K8sResourceAgnosticResponse) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

func (in *K8sResourceAgnosticResponse) DeepCopyObject() runtime.Object {
	if out := in.DeepCopy(); out != nil {
		return out
	}
	return nil
}
