package resourceconversionprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pmetric"
	// "github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
)

func TestConvertResourceToAttributes(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("service.name", "test")
	rm.Resource().Attributes().PutStr("service.no_copied", "test")

	sm := rm.ScopeMetrics().AppendEmpty()
	m0 := sm.Metrics().AppendEmpty()
	m0g := m0.SetEmptyGauge()
	m0g.DataPoints().AppendEmpty()
	m0.SetName("test")

	assert.Equal(t, 2, md.ResourceMetrics().At(0).Resource().Attributes().Len())
	assert.Equal(t, 0, md.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Gauge().DataPoints().At(0).Attributes().Len())

	processor := resourceProcessor{}
	processor.processMetrics(context.Background(), md)

	assert.Equal(t, 2, md.ResourceMetrics().At(0).Resource().Attributes().Len())
	assert.Equal(t, 1, md.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Gauge().DataPoints().At(0).Attributes().Len())
}
