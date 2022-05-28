// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// MetricSettings provides common settings for a particular metric.
type MetricSettings struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsSettings provides settings for vaultkvreceiver metrics.
type MetricsSettings struct {
	VaultkvCreatedOn     MetricSettings `mapstructure:"vaultkv.created_on"`
	VaultkvMetadata      MetricSettings `mapstructure:"vaultkv.metadata"`
	VaultkvMetadataError MetricSettings `mapstructure:"vaultkv.metadata.error"`
}

func DefaultMetricsSettings() MetricsSettings {
	return MetricsSettings{
		VaultkvCreatedOn: MetricSettings{
			Enabled: true,
		},
		VaultkvMetadata: MetricSettings{
			Enabled: true,
		},
		VaultkvMetadataError: MetricSettings{
			Enabled: true,
		},
	}
}

// AttributeMetadataErrorType specifies the a value metadata_error_type attribute.
type AttributeMetadataErrorType int

const (
	_ AttributeMetadataErrorType = iota
	AttributeMetadataErrorTypeMissingType
	AttributeMetadataErrorTypeInvalidType
)

// String returns the string representation of the AttributeMetadataErrorType.
func (av AttributeMetadataErrorType) String() string {
	switch av {
	case AttributeMetadataErrorTypeMissingType:
		return "missing_type"
	case AttributeMetadataErrorTypeInvalidType:
		return "invalid_type"
	}
	return ""
}

// MapAttributeMetadataErrorType is a helper map of string to AttributeMetadataErrorType attribute value.
var MapAttributeMetadataErrorType = map[string]AttributeMetadataErrorType{
	"missing_type": AttributeMetadataErrorTypeMissingType,
	"invalid_type": AttributeMetadataErrorTypeInvalidType,
}

// AttributeType specifies the a value type attribute.
type AttributeType int

const (
	_ AttributeType = iota
	AttributeTypeDigitaloceanSpaces
	AttributeTypeDigitaloceanAPI
	AttributeTypeTailscaleAPI
	AttributeTypeConsulEncryption
	AttributeTypeNomadEncryption
)

// String returns the string representation of the AttributeType.
func (av AttributeType) String() string {
	switch av {
	case AttributeTypeDigitaloceanSpaces:
		return "digitalocean.spaces"
	case AttributeTypeDigitaloceanAPI:
		return "digitalocean.api"
	case AttributeTypeTailscaleAPI:
		return "tailscale.api"
	case AttributeTypeConsulEncryption:
		return "consul.encryption"
	case AttributeTypeNomadEncryption:
		return "nomad.encryption"
	}
	return ""
}

// MapAttributeType is a helper map of string to AttributeType attribute value.
var MapAttributeType = map[string]AttributeType{
	"digitalocean.spaces": AttributeTypeDigitaloceanSpaces,
	"digitalocean.api":    AttributeTypeDigitaloceanAPI,
	"tailscale.api":       AttributeTypeTailscaleAPI,
	"consul.encryption":   AttributeTypeConsulEncryption,
	"nomad.encryption":    AttributeTypeNomadEncryption,
}

type metricVaultkvCreatedOn struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills vaultkv.created_on metric with initial data.
func (m *metricVaultkvCreatedOn) init() {
	m.data.SetName("vaultkv.created_on")
	m.data.SetDescription("The epoch time in seconds the key was created at.")
	m.data.SetUnit("s")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(true)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricVaultkvCreatedOn) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, typeAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("key", pcommon.NewValueString(keyAttributeValue))
	dp.Attributes().Insert("mount", pcommon.NewValueString(mountAttributeValue))
	dp.Attributes().Insert("type", pcommon.NewValueString(typeAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricVaultkvCreatedOn) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricVaultkvCreatedOn) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricVaultkvCreatedOn(settings MetricSettings) metricVaultkvCreatedOn {
	m := metricVaultkvCreatedOn{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricVaultkvMetadata struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills vaultkv.metadata metric with initial data.
func (m *metricVaultkvMetadata) init() {
	m.data.SetName("vaultkv.metadata")
	m.data.SetDescription("Metadata about the key.")
	m.data.SetUnit("1")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricVaultkvMetadata) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, versionsAttributeValue string, currentVersionAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("key", pcommon.NewValueString(keyAttributeValue))
	dp.Attributes().Insert("mount", pcommon.NewValueString(mountAttributeValue))
	dp.Attributes().Insert("versions", pcommon.NewValueString(versionsAttributeValue))
	dp.Attributes().Insert("current_version", pcommon.NewValueString(currentVersionAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricVaultkvMetadata) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricVaultkvMetadata) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricVaultkvMetadata(settings MetricSettings) metricVaultkvMetadata {
	m := metricVaultkvMetadata{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricVaultkvMetadataError struct {
	data     pmetric.Metric // data buffer for generated metric.
	settings MetricSettings // metric settings provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills vaultkv.metadata.error metric with initial data.
func (m *metricVaultkvMetadataError) init() {
	m.data.SetName("vaultkv.metadata.error")
	m.data.SetDescription("The epoch time in seconds the key was created at.")
	m.data.SetUnit("1")
	m.data.SetDataType(pmetric.MetricDataTypeSum)
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricVaultkvMetadataError) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, metadataErrorTypeAttributeValue string) {
	if !m.settings.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntVal(val)
	dp.Attributes().Insert("key", pcommon.NewValueString(keyAttributeValue))
	dp.Attributes().Insert("mount", pcommon.NewValueString(mountAttributeValue))
	dp.Attributes().Insert("type", pcommon.NewValueString(metadataErrorTypeAttributeValue))
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricVaultkvMetadataError) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricVaultkvMetadataError) emit(metrics pmetric.MetricSlice) {
	if m.settings.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricVaultkvMetadataError(settings MetricSettings) metricVaultkvMetadataError {
	m := metricVaultkvMetadataError{settings: settings}
	if settings.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user settings.
type MetricsBuilder struct {
	startTime                  pcommon.Timestamp   // start time that will be applied to all recorded data points.
	metricsCapacity            int                 // maximum observed number of metrics per resource.
	resourceCapacity           int                 // maximum observed number of resource attributes.
	metricsBuffer              pmetric.Metrics     // accumulates metrics data before emitting.
	buildInfo                  component.BuildInfo // contains version information
	metricVaultkvCreatedOn     metricVaultkvCreatedOn
	metricVaultkvMetadata      metricVaultkvMetadata
	metricVaultkvMetadataError metricVaultkvMetadataError
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pcommon.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(settings MetricsSettings, buildInfo component.BuildInfo, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                  pcommon.NewTimestampFromTime(time.Now()),
		metricsBuffer:              pmetric.NewMetrics(),
		buildInfo:                  buildInfo,
		metricVaultkvCreatedOn:     newMetricVaultkvCreatedOn(settings.VaultkvCreatedOn),
		metricVaultkvMetadata:      newMetricVaultkvMetadata(settings.VaultkvMetadata),
		metricVaultkvMetadataError: newMetricVaultkvMetadataError(settings.VaultkvMetadataError),
	}
	for _, op := range options {
		op(mb)
	}
	return mb
}

// updateCapacity updates max length of metrics and resource attributes that will be used for the slice capacity.
func (mb *MetricsBuilder) updateCapacity(rm pmetric.ResourceMetrics) {
	if mb.metricsCapacity < rm.ScopeMetrics().At(0).Metrics().Len() {
		mb.metricsCapacity = rm.ScopeMetrics().At(0).Metrics().Len()
	}
	if mb.resourceCapacity < rm.Resource().Attributes().Len() {
		mb.resourceCapacity = rm.Resource().Attributes().Len()
	}
}

// ResourceMetricsOption applies changes to provided resource metrics.
type ResourceMetricsOption func(pmetric.ResourceMetrics)

// WithStartTimeOverride overrides start time for all the resource metrics data points.
// This option should be only used if different start time has to be set on metrics coming from different resources.
func WithStartTimeOverride(start pcommon.Timestamp) ResourceMetricsOption {
	return func(rm pmetric.ResourceMetrics) {
		var dps pmetric.NumberDataPointSlice
		metrics := rm.ScopeMetrics().At(0).Metrics()
		for i := 0; i < metrics.Len(); i++ {
			switch metrics.At(i).DataType() {
			case pmetric.MetricDataTypeGauge:
				dps = metrics.At(i).Gauge().DataPoints()
			case pmetric.MetricDataTypeSum:
				dps = metrics.At(i).Sum().DataPoints()
			}
			for j := 0; j < dps.Len(); j++ {
				dps.At(j).SetStartTimestamp(start)
			}
		}
	}
}

// EmitForResource saves all the generated metrics under a new resource and updates the internal state to be ready for
// recording another set of data points as part of another resource. This function can be helpful when one scraper
// needs to emit metrics from several resources. Otherwise calling this function is not required,
// just `Emit` function can be called instead.
// Resource attributes should be provided as ResourceMetricsOption arguments.
func (mb *MetricsBuilder) EmitForResource(rmo ...ResourceMetricsOption) {
	rm := pmetric.NewResourceMetrics()
	rm.Resource().Attributes().EnsureCapacity(mb.resourceCapacity)
	ils := rm.ScopeMetrics().AppendEmpty()
	ils.Scope().SetName("otelcol/vaultkvreceiver")
	ils.Scope().SetVersion(mb.buildInfo.Version)
	ils.Metrics().EnsureCapacity(mb.metricsCapacity)
	mb.metricVaultkvCreatedOn.emit(ils.Metrics())
	mb.metricVaultkvMetadata.emit(ils.Metrics())
	mb.metricVaultkvMetadataError.emit(ils.Metrics())
	for _, op := range rmo {
		op(rm)
	}
	if ils.Metrics().Len() > 0 {
		mb.updateCapacity(rm)
		rm.MoveTo(mb.metricsBuffer.ResourceMetrics().AppendEmpty())
	}
}

// Emit returns all the metrics accumulated by the metrics builder and updates the internal state to be ready for
// recording another set of metrics. This function will be responsible for applying all the transformations required to
// produce metric representation defined in metadata and user settings, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit(rmo ...ResourceMetricsOption) pmetric.Metrics {
	mb.EmitForResource(rmo...)
	metrics := pmetric.NewMetrics()
	mb.metricsBuffer.MoveTo(metrics)
	return metrics
}

// RecordVaultkvCreatedOnDataPoint adds a data point to vaultkv.created_on metric.
func (mb *MetricsBuilder) RecordVaultkvCreatedOnDataPoint(ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, typeAttributeValue AttributeType) {
	mb.metricVaultkvCreatedOn.recordDataPoint(mb.startTime, ts, val, keyAttributeValue, mountAttributeValue, typeAttributeValue.String())
}

// RecordVaultkvMetadataDataPoint adds a data point to vaultkv.metadata metric.
func (mb *MetricsBuilder) RecordVaultkvMetadataDataPoint(ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, versionsAttributeValue string, currentVersionAttributeValue string) {
	mb.metricVaultkvMetadata.recordDataPoint(mb.startTime, ts, val, keyAttributeValue, mountAttributeValue, versionsAttributeValue, currentVersionAttributeValue)
}

// RecordVaultkvMetadataErrorDataPoint adds a data point to vaultkv.metadata.error metric.
func (mb *MetricsBuilder) RecordVaultkvMetadataErrorDataPoint(ts pcommon.Timestamp, val int64, keyAttributeValue string, mountAttributeValue string, metadataErrorTypeAttributeValue AttributeMetadataErrorType) {
	mb.metricVaultkvMetadataError.recordDataPoint(mb.startTime, ts, val, keyAttributeValue, mountAttributeValue, metadataErrorTypeAttributeValue.String())
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}
