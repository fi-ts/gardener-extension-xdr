package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XDRConfig defines the configuration for the xdr controller.
type XDRConfig struct {
	metav1.TypeMeta `json:",inline"`
	NoProxy         bool     `json:"noProxy,omitempty"`
	ProxyList       []string `json:"proxyList,omitempty"`
	DistributionId  string   `json:"distributionId,omitempty"`
	CustomTag       string   `json:"customTag,omitempty"`
	Tenant          string   `json:"tenant,omitempty"`
}
