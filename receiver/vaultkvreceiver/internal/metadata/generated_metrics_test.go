// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testConfigCollection int

const (
	testSetDefault testConfigCollection = iota
	testSetAll
	testSetNone
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		configSet testConfigCollection
	}{
		{
			name:      "default",
			configSet: testSetDefault,
		},
		{
			name:      "all_set",
			configSet: testSetAll,
		},
		{
			name:      "none_set",
			configSet: testSetNone,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := receivertest.NewNopCreateSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadMetricsBuilderConfig(t, test.name), settings, WithStartTime(start))

			expectedWarnings := 0

			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordVaultkvCreatedOnDataPoint(ts, 1, "key-val", "mount-val", AttributeTypeDigitaloceanSpaces)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordVaultkvMetadataDataPoint(ts, 1, "key-val", "mount-val", "versions-val", "current_version-val")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordVaultkvMetadataErrorDataPoint(ts, 1, "key-val", "mount-val", AttributeMetadataErrorTypeMissingType)

			res := pcommon.NewResource()
			metrics := mb.Emit(WithResource(res))

			if test.configSet == testSetNone {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			assert.Equal(t, res, rm.Resource())
			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if test.configSet == testSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if test.configSet == testSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "vaultkv.created_on":
					assert.False(t, validatedMetrics["vaultkv.created_on"], "Found a duplicate in the metrics slice: vaultkv.created_on")
					validatedMetrics["vaultkv.created_on"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The epoch time in seconds the key was created at.", ms.At(i).Description())
					assert.Equal(t, "seconds", ms.At(i).Unit())
					assert.Equal(t, true, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("key")
					assert.True(t, ok)
					assert.EqualValues(t, "key-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("mount")
					assert.True(t, ok)
					assert.EqualValues(t, "mount-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("type")
					assert.True(t, ok)
					assert.EqualValues(t, "digitalocean.spaces", attrVal.Str())
				case "vaultkv.metadata":
					assert.False(t, validatedMetrics["vaultkv.metadata"], "Found a duplicate in the metrics slice: vaultkv.metadata")
					validatedMetrics["vaultkv.metadata"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Metadata about the key.", ms.At(i).Description())
					assert.Equal(t, "1", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("key")
					assert.True(t, ok)
					assert.EqualValues(t, "key-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("mount")
					assert.True(t, ok)
					assert.EqualValues(t, "mount-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("versions")
					assert.True(t, ok)
					assert.EqualValues(t, "versions-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("current_version")
					assert.True(t, ok)
					assert.EqualValues(t, "current_version-val", attrVal.Str())
				case "vaultkv.metadata.error":
					assert.False(t, validatedMetrics["vaultkv.metadata.error"], "Found a duplicate in the metrics slice: vaultkv.metadata.error")
					validatedMetrics["vaultkv.metadata.error"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Errors reported while trying to fetch metrics.", ms.At(i).Description())
					assert.Equal(t, "1", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("key")
					assert.True(t, ok)
					assert.EqualValues(t, "key-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("mount")
					assert.True(t, ok)
					assert.EqualValues(t, "mount-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("metadata_error_type")
					assert.True(t, ok)
					assert.EqualValues(t, "missing_type", attrVal.Str())
				}
			}
		})
	}
}
