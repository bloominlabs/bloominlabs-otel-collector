// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type metricDigitaloceanBillingBalance struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills digitalocean.billing.balance metric with initial data.
func (m *metricDigitaloceanBillingBalance) init() {
	m.data.SetName("digitalocean.billing.balance")
	m.data.SetDescription("Balance as of `digitalocean.billing.generate_at` time.")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
}

func (m *metricDigitaloceanBillingBalance) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleValue(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricDigitaloceanBillingBalance) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricDigitaloceanBillingBalance) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricDigitaloceanBillingBalance(cfg MetricConfig) metricDigitaloceanBillingBalance {
	m := metricDigitaloceanBillingBalance{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricDigitaloceanBillingGenerateAt struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills digitalocean.billing.generate_at metric with initial data.
func (m *metricDigitaloceanBillingGenerateAt) init() {
	m.data.SetName("digitalocean.billing.generate_at")
	m.data.SetDescription("The time at which balances were most recently generated.")
	m.data.SetUnit("seconds")
	m.data.SetEmptyGauge()
}

func (m *metricDigitaloceanBillingGenerateAt) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricDigitaloceanBillingGenerateAt) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricDigitaloceanBillingGenerateAt) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricDigitaloceanBillingGenerateAt(cfg MetricConfig) metricDigitaloceanBillingGenerateAt {
	m := metricDigitaloceanBillingGenerateAt{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricDigitaloceanBillingUsage struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills digitalocean.billing.usage metric with initial data.
func (m *metricDigitaloceanBillingUsage) init() {
	m.data.SetName("digitalocean.billing.usage")
	m.data.SetDescription("Amount used in the current billing period as of the `digitalocean.billing.generate_at` time")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
}

func (m *metricDigitaloceanBillingUsage) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val float64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetDoubleValue(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricDigitaloceanBillingUsage) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricDigitaloceanBillingUsage) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricDigitaloceanBillingUsage(cfg MetricConfig) metricDigitaloceanBillingUsage {
	m := metricDigitaloceanBillingUsage{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user config.
type MetricsBuilder struct {
	startTime                           pcommon.Timestamp   // start time that will be applied to all recorded data points.
	metricsCapacity                     int                 // maximum observed number of metrics per resource.
	resourceCapacity                    int                 // maximum observed number of resource attributes.
	metricsBuffer                       pmetric.Metrics     // accumulates metrics data before emitting.
	buildInfo                           component.BuildInfo // contains version information
	metricDigitaloceanBillingBalance    metricDigitaloceanBillingBalance
	metricDigitaloceanBillingGenerateAt metricDigitaloceanBillingGenerateAt
	metricDigitaloceanBillingUsage      metricDigitaloceanBillingUsage
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pcommon.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(mbc MetricsBuilderConfig, settings receiver.CreateSettings, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		startTime:                           pcommon.NewTimestampFromTime(time.Now()),
		metricsBuffer:                       pmetric.NewMetrics(),
		buildInfo:                           settings.BuildInfo,
		metricDigitaloceanBillingBalance:    newMetricDigitaloceanBillingBalance(mbc.Metrics.DigitaloceanBillingBalance),
		metricDigitaloceanBillingGenerateAt: newMetricDigitaloceanBillingGenerateAt(mbc.Metrics.DigitaloceanBillingGenerateAt),
		metricDigitaloceanBillingUsage:      newMetricDigitaloceanBillingUsage(mbc.Metrics.DigitaloceanBillingUsage),
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
			switch metrics.At(i).Type() {
			case pmetric.MetricTypeGauge:
				dps = metrics.At(i).Gauge().DataPoints()
			case pmetric.MetricTypeSum:
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
	ils.Scope().SetName("otelcol/digitaloceanreceiver/billing")
	ils.Scope().SetVersion(mb.buildInfo.Version)
	ils.Metrics().EnsureCapacity(mb.metricsCapacity)
	mb.metricDigitaloceanBillingBalance.emit(ils.Metrics())
	mb.metricDigitaloceanBillingGenerateAt.emit(ils.Metrics())
	mb.metricDigitaloceanBillingUsage.emit(ils.Metrics())

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
// produce metric representation defined in metadata and user config, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit(rmo ...ResourceMetricsOption) pmetric.Metrics {
	mb.EmitForResource(rmo...)
	metrics := mb.metricsBuffer
	mb.metricsBuffer = pmetric.NewMetrics()
	return metrics
}

// RecordDigitaloceanBillingBalanceDataPoint adds a data point to digitalocean.billing.balance metric.
func (mb *MetricsBuilder) RecordDigitaloceanBillingBalanceDataPoint(ts pcommon.Timestamp, val float64) {
	mb.metricDigitaloceanBillingBalance.recordDataPoint(mb.startTime, ts, val)
}

// RecordDigitaloceanBillingGenerateAtDataPoint adds a data point to digitalocean.billing.generate_at metric.
func (mb *MetricsBuilder) RecordDigitaloceanBillingGenerateAtDataPoint(ts pcommon.Timestamp, val int64) {
	mb.metricDigitaloceanBillingGenerateAt.recordDataPoint(mb.startTime, ts, val)
}

// RecordDigitaloceanBillingUsageDataPoint adds a data point to digitalocean.billing.usage metric.
func (mb *MetricsBuilder) RecordDigitaloceanBillingUsageDataPoint(ts pcommon.Timestamp, val float64) {
	mb.metricDigitaloceanBillingUsage.recordDataPoint(mb.startTime, ts, val)
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}
