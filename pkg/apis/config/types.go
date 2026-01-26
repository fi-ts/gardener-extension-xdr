package config

import (
	"strings"

	healthcheckconfig "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControllerConfiguration defines the configuration for the xdr controller.
type ControllerConfiguration struct {
	metav1.TypeMeta

	// HealthCheckConfig is the config for the health check controller
	HealthCheckConfig *healthcheckconfig.HealthCheckConfig
	TenantConfigs     []TenantConfig
}

type TenantConfig struct {
	TenantId       string   `json:"tenantId"`
	ProxyList      []string `json:"proxyList,omitempty"`
	DistributionId string   `json:"distributionId,omitempty"`
	CustomTag      string   `json:"customTag,omitempty"`
}

type TenantConfigs []TenantConfig

func (tc TenantConfigs) GetTenantConfig(tenantId string) *TenantConfig {
	for _, config := range tc {
		// compare tenant-ids without case, so "ft" is same as "Ft"
		if strings.EqualFold(config.TenantId, tenantId) {
			return &config
		}
	}
	return nil
}
