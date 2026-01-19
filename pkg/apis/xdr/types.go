package xdr

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XDRConfig defines the configuration for the xdr controller.
type XDRConfig struct {
	metav1.TypeMeta
	NoProxy   bool
	CustomTag string
	Tenant    string
}
