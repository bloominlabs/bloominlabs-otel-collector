// Code generated by "go.opentelemetry.io/collector/cmd/builder". DO NOT EDIT.

package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	loggingexporter "go.opentelemetry.io/collector/exporter/loggingexporter"
	otlpexporter "go.opentelemetry.io/collector/exporter/otlpexporter"
	otlphttpexporter "go.opentelemetry.io/collector/exporter/otlphttpexporter"
	zpagesextension "go.opentelemetry.io/collector/extension/zpagesextension"
	basicauthextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension"
	pprofextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension"
	filestorage "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage"
	healthcheckextension "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension"
	batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
	resourceprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	resourceconversionprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceconversionprocessor"
	memorylimiterprocessor "go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	filterprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	lokiprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/lokiprocessor"
	resourcedetectionprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	attributesprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	nomadprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/nomadprocessor"
	metricstransformprocessor "github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor"
	hostmetricsreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver"
	otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
	journaldreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver"
	prometheusreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	filelogreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	vaultkvreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver"
	jaegerreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver"
	postgresqlreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver"
	chronyreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver"
	digitaloceanreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver"
	userstatsreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver"
	certificatesreceiver "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver"
)

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	factories.Extensions, err = extension.MakeFactoryMap(
		zpagesextension.NewFactory(),
		basicauthextension.NewFactory(),
		pprofextension.NewFactory(),
		filestorage.NewFactory(),
		healthcheckextension.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ExtensionModules = make(map[component.Type]string, len(factories.Extensions))
	factories.ExtensionModules[zpagesextension.NewFactory().Type()] = "go.opentelemetry.io/collector/extension/zpagesextension v0.106.1"
	factories.ExtensionModules[basicauthextension.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.106.1"
	factories.ExtensionModules[pprofextension.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.106.1"
	factories.ExtensionModules[filestorage.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.106.1"
	factories.ExtensionModules[healthcheckextension.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.106.1"

	factories.Receivers, err = receiver.MakeFactoryMap(
		hostmetricsreceiver.NewFactory(),
		otlpreceiver.NewFactory(),
		journaldreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		vaultkvreceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		postgresqlreceiver.NewFactory(),
		chronyreceiver.NewFactory(),
		digitaloceanreceiver.NewFactory(),
		userstatsreceiver.NewFactory(),
		certificatesreceiver.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ReceiverModules = make(map[component.Type]string, len(factories.Receivers))
	factories.ReceiverModules[hostmetricsreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.106.1"
	factories.ReceiverModules[otlpreceiver.NewFactory().Type()] = "go.opentelemetry.io/collector/receiver/otlpreceiver v0.106.1"
	factories.ReceiverModules[journaldreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.106.1"
	factories.ReceiverModules[prometheusreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.106.1"
	factories.ReceiverModules[filelogreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.106.1"
	factories.ReceiverModules[vaultkvreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver v0.106.1"
	factories.ReceiverModules[jaegerreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.106.1"
	factories.ReceiverModules[postgresqlreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.106.1"
	factories.ReceiverModules[chronyreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver v0.106.1"
	factories.ReceiverModules[digitaloceanreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver v0.106.1"
	factories.ReceiverModules[userstatsreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver v0.106.1"
	factories.ReceiverModules[certificatesreceiver.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver v0.106.1"

	factories.Exporters, err = exporter.MakeFactoryMap(
		loggingexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ExporterModules = make(map[component.Type]string, len(factories.Exporters))
	factories.ExporterModules[loggingexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/loggingexporter v0.106.1"
	factories.ExporterModules[otlpexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/otlpexporter v0.106.1"
	factories.ExporterModules[otlphttpexporter.NewFactory().Type()] = "go.opentelemetry.io/collector/exporter/otlphttpexporter v0.106.1"

	factories.Processors, err = processor.MakeFactoryMap(
		batchprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		resourceconversionprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		filterprocessor.NewFactory(),
		lokiprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		attributesprocessor.NewFactory(),
		nomadprocessor.NewFactory(),
		metricstransformprocessor.NewFactory(),
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ProcessorModules = make(map[component.Type]string, len(factories.Processors))
	factories.ProcessorModules[batchprocessor.NewFactory().Type()] = "go.opentelemetry.io/collector/processor/batchprocessor v0.106.1"
	factories.ProcessorModules[resourceprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.106.1"
	factories.ProcessorModules[resourceconversionprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceconversionprocessor v0.106.1"
	factories.ProcessorModules[memorylimiterprocessor.NewFactory().Type()] = "go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.106.1"
	factories.ProcessorModules[filterprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.106.1"
	factories.ProcessorModules[lokiprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/lokiprocessor v0.106.1"
	factories.ProcessorModules[resourcedetectionprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.106.1"
	factories.ProcessorModules[attributesprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.106.1"
	factories.ProcessorModules[nomadprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/nomadprocessor v0.106.1"
	factories.ProcessorModules[metricstransformprocessor.NewFactory().Type()] = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.106.1"

	factories.Connectors, err = connector.MakeFactoryMap(
	)
	if err != nil {
		return otelcol.Factories{}, err
	}
	factories.ConnectorModules = make(map[component.Type]string, len(factories.Connectors))

	return factories, nil
}
