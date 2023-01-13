// Code generated by "go.opentelemetry.io/collector/cmd/builder". DO NOT EDIT.

package main

import (
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	loggingexporter "go.opentelemetry.io/collector/exporter/loggingexporter"
	otlpexporter "go.opentelemetry.io/collector/exporter/otlpexporter"
	lokiexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter"
	prometheusremotewriteexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	ballastextension "go.opentelemetry.io/collector/extension/ballastextension"
	zpagesextension "go.opentelemetry.io/collector/extension/zpagesextension"
	pprofextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	filestorage "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	healthcheckextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	memorylimiterprocessor "go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
	attributesprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	resourceprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	resourcedetectionprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	lokiprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/lokiprocessor"
	nomadprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/nomadprocessor"
	metricstransformprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
	filelogreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	vaultkvreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver"
	journaldreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver"
	jaegerreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver"
	prometheusreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	prometheusexecreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver"
	hostmetricsreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	postgresqlreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
)

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Extensions, err = extension.MakeFactoryMap(
		ballastextension.NewFactory(),
		zpagesextension.NewFactory(),
		pprofextension.NewFactory(),
		filestorage.NewFactory(),
		healthcheckextension.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Receivers, err = receiver.MakeFactoryMap(
		otlpreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		vaultkvreceiver.NewFactory(),
		journaldreceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		prometheusexecreceiver.NewFactory(),
		hostmetricsreceiver.NewFactory(),
		postgresqlreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Exporters, err = exporter.MakeFactoryMap(
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		lokiexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Processors, err = processor.MakeFactoryMap(
		memorylimiterprocessor.NewFactory(),
		batchprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		lokiprocessor.NewFactory(),
		nomadprocessor.NewFactory(),
		metricstransformprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}
