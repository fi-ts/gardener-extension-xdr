package xdr

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fi-ts/gardener-extension-xdr/charts"
	"github.com/fi-ts/gardener-extension-xdr/pkg/apis/xdr/v1alpha1"
	"github.com/fi-ts/gardener-extension-xdr/pkg/imagevector"
	"github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/fi-ts/gardener-extension-xdr/pkg/apis/config"
	gardener "github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/utils/managedresources"
	"github.com/go-logr/logr"
	"github.com/metal-stack/metal-lib/pkg/tag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	managedResourceName = "xdr-resource"
)

// NewActuator returns an actuator responsible for Extension resources.
func NewActuator(mgr manager.Manager, config config.ControllerConfiguration) (extension.Actuator, error) {
	ca, err := gardener.NewChartApplierForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	return &actuator{
		client:       mgr.GetClient(),
		decoder:      serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder(),
		config:       config,
		chartApplier: ca,
	}, nil
}

type actuator struct {
	client       client.Client
	decoder      runtime.Decoder
	config       config.ControllerConfiguration
	chartApplier gardener.ChartApplier
}

// Reconcile the Extension resource.
func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	cluster, err := controller.GetCluster(ctx, a.client, ex.GetNamespace())
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	tm := tag.TagMap(cluster.Shoot.Annotations)
	clusterid, _ := tm.Value(tag.ClusterID)
	clustername, _ := tm.Value(tag.ClusterName)
	tenant, _ := tm.Value(tag.ClusterTenant)

	cortextImage, err := imagevector.ImageVector().FindImage("cortex-agent")
	if err != nil {
		return fmt.Errorf("failed to find cortext-agent image: %w", err)
	}

	ci, err := util.NewChartRendererForShoot(cluster.Shoot.Spec.Kubernetes.Version)
	if err != nil {
		return fmt.Errorf("failed to create chart renderer: %w", err)
	}

	var xdrConfig v1alpha1.XDRConfig
	if ex.Spec.ProviderConfig != nil {
		_, _, err := a.decoder.Decode(ex.Spec.ProviderConfig.Raw, nil, &xdrConfig)
		if err != nil {
			return fmt.Errorf("failed to decode provider config: %w", err)
		}
	}

	endpointTags := fmt.Sprintf("tenant=%s;clusterid=%s", tenant, clusterid)
	distributionId := getValue(xdrConfig.DistributionId, a.config.DefaultDistributionId)
	proxyList := getSliceValue(xdrConfig.ProxyList, a.config.DefaultProxyList)

	if xdrConfig.CustomTag == "" {
		endpointTags = fmt.Sprintf("%s,custom=%s", endpointTags, xdrConfig.CustomTag)
	}

	rc, err := ci.RenderEmbeddedFS(charts.InternalChart, filepath.Join("internal", charts.CortexChartsPath), charts.CortextName, charts.CortexNamespace, map[string]any{
		"namespace": map[string]any{
			"create": false,
			"name":   charts.CortexNamespace,
		},
		"agent": map[string]any{
			"endpointTags":   endpointTags,
			"clusterName":    clustername,
			"distributionId": distributionId,
			"proxyList":      proxyList,
		},
		"daemonset": map[string]any{
			"image": map[string]any{
				"repository": cortextImage.Repository,
				"tag":        cortextImage.Tag,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to render chart: %w", err)
	}

	data := map[string][]byte{
		charts.CortextName: rc.Manifest(),
	}

	log.Info("reconciling extension", "configuration", data)

	err = managedresources.CreateForShoot(ctx, a.client, ex.GetNamespace(), managedResourceName, "", false, data)

	if err != nil {
		return fmt.Errorf("failed to apply chart: %w", err)
	}
	log.Info("reconciling extension", "configuration", xdrConfig)
	return nil

}

// Delete the Extension resource.
func (a *actuator) Delete(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	log.Info("deleting managed resource")
	err := managedresources.DeleteForShoot(ctx, a.client, ex.GetNamespace(), managedResourceName)
	if err != nil {
		log.Error(err, "cannot delete managed resource")
	}
	return err
}

// ForceDelete the Extension resource
func (a *actuator) ForceDelete(_ context.Context, _ logr.Logger, _ *extensionsv1alpha1.Extension) error {
	return nil
}

// Restore the Extension resource.
func (a *actuator) Restore(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, log, ex)
}

// Migrate the Extension resource.
func (a *actuator) Migrate(ctx context.Context, log logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return nil
}

func getValue[T comparable](val T, defVal T) T {
	var zero T
	if val == zero {
		return defVal
	}
	return val
}

func getSliceValue[S ~[]T, T comparable](val S, defVal S) S {
	if len(val) == 0 {
		return defVal
	}
	return val
}
