package config

import (
	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerConfiguration defines the configuration for the xdr controller.
type ControllerConfiguration struct {
	metav1.TypeMeta

	// HealthCheckConfig is the config for the health check controller
	HealthCheckConfig     *healthcheckconfig.HealthCheckConfig
	DefaultProxyList      []string
	DefaultDistributionId string
}
