
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)


var GroupVersion = schema.GroupVersion{Group: "example.com", Version: "v1"}


var (
	SchemeBuilder = &runtime.SchemeBuilder{AddKnownTypes}
	AddToScheme   = SchemeBuilder.AddToScheme
)


func AddKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&ClusterScan{},
		&ClusterScanList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
