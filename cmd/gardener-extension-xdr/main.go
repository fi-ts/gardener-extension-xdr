package main

import (
	"os"

	logger "github.com/gardener/gardener/pkg/logger"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/fi-ts/gardener-extension-xdr/cmd/gardener-extension-xdr/app"
)

func main() {
	runtimelog.SetLogger(logger.MustNewZapLogger(logger.InfoLevel, logger.FormatJSON))
	cmd := app.NewControllerManagerCommand(signals.SetupSignalHandler())

	if err := cmd.Execute(); err != nil {
		runtimelog.Log.Error(err, "error executing the main controller command")
		os.Exit(1)
	}
}
