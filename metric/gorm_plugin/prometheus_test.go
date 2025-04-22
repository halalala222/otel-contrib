package gorm_plugin

import (
	"database/sql"
	"testing"

	"github.com/ncuhome/otel-contrib/metric/test"

	. "github.com/bytedance/mockey"
	"github.com/go-playground/assert/v2"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

const (
	testMaxOpenConnections = iota
	testOpenConnections
	testInUse
	testIdle
	testWaitCount
	testWaitDuration
	testMaxIdleClosed
	testMaxLifetimeClosed
	testMaxIdleTimeClosed
)

func TestPrometheus_Initialize(t *testing.T) {
	testProvider := &test.Provider{}
	otel.SetMeterProvider(testProvider)
	PatchConvey("TestPrometheus_Initialize", t, func() {
		Mock((*sql.DB).Stats).Return(sql.DBStats{
			MaxOpenConnections: testMaxOpenConnections,
			OpenConnections:    testOpenConnections,
			InUse:              testInUse,
			Idle:               testIdle,
			WaitCount:          testWaitCount,
			WaitDuration:       testWaitDuration,
			MaxIdleClosed:      testMaxIdleClosed,
			MaxLifetimeClosed:  testMaxLifetimeClosed,
			MaxIdleTimeClosed:  testMaxIdleTimeClosed,
		}).Build()
		Mock((*gorm.DB).DB).Return(new(sql.DB), nil).Build()

		p := New()

		if err := p.Initialize(new(gorm.DB)); err != nil {
			t.Error(err)
		}
	})

	maxOpenConnections := testProvider.GetMeterFloat64RecordData("gorm-db-status-max-open-connections")
	openConnections := testProvider.GetMeterFloat64RecordData("gorm-db-status-open-connections")
	inUse := testProvider.GetMeterFloat64RecordData("gorm-db-status-in-use")
	idle := testProvider.GetMeterFloat64RecordData("gorm-db-status-idle")
	waitCount := testProvider.GetMeterFloat64RecordData("gorm-db-status-wait-count")
	waitDuration := testProvider.GetMeterFloat64RecordData("gorm-db-status-wait-duration")
	maxIdleClosed := testProvider.GetMeterFloat64RecordData("gorm-db-status-max-idle-closed")
	maxLifetimeClosed := testProvider.GetMeterFloat64RecordData("gorm-db-status-max-lifetime-closed")
	maxIdleTimeClosed := testProvider.GetMeterFloat64RecordData("gorm-db-status-max-idle-time-closed")

	assert.Equal(t, maxOpenConnections, float64(testMaxOpenConnections))
	assert.Equal(t, openConnections, float64(testOpenConnections))
	assert.Equal(t, inUse, float64(testInUse))
	assert.Equal(t, idle, float64(testIdle))
	assert.Equal(t, waitCount, float64(testWaitCount))
	assert.Equal(t, waitDuration, float64(testWaitDuration))
	assert.Equal(t, maxIdleClosed, float64(testMaxIdleClosed))
	assert.Equal(t, maxLifetimeClosed, float64(testMaxLifetimeClosed))
	assert.Equal(t, maxIdleTimeClosed, float64(testMaxIdleTimeClosed))
}
