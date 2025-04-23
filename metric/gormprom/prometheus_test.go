package gormprom

import (
	"context"
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
		err := p.Initialize(new(gorm.DB))
		assert.Equal(t, err, nil)

		tetsCallbacks := testProvider.GetGlobalMeterCallbacks()
		assert.Equal(t, len(tetsCallbacks), 1)

		testObserver := &test.Observer{}
		err = tetsCallbacks[0](context.Background(), testObserver)
		assert.Equal(t, err, nil)
	})

	maxOpenConnections := testProvider.GetMeterInt64RecordData(maxOpenConnsName)
	openConnections := testProvider.GetMeterInt64RecordData(openConnsName)
	inUse := testProvider.GetMeterInt64RecordData(inUseConnsName)
	idle := testProvider.GetMeterInt64RecordData(idleConnsName)
	waitCount := testProvider.GetMeterInt64RecordData(connsWaitCountName)
	waitDuration := testProvider.GetMeterInt64RecordData(connsWaitDurationName)
	maxIdleClosed := testProvider.GetMeterInt64RecordData(connsClosedMaxIdleName)
	maxLifetimeClosed := testProvider.GetMeterInt64RecordData(connsClosedMaxLifetimeName)
	maxIdleTimeClosed := testProvider.GetMeterInt64RecordData(connsClosedMaxIdleTimeName)

	assert.Equal(t, maxOpenConnections, int64(testMaxOpenConnections))
	assert.Equal(t, openConnections, int64(testOpenConnections))
	assert.Equal(t, inUse, int64(testInUse))
	assert.Equal(t, idle, int64(testIdle))
	assert.Equal(t, waitCount, int64(testWaitCount))
	assert.Equal(t, waitDuration, int64(testWaitDuration))
	assert.Equal(t, maxIdleClosed, int64(testMaxIdleClosed))
	assert.Equal(t, maxLifetimeClosed, int64(testMaxLifetimeClosed))
	assert.Equal(t, maxIdleTimeClosed, int64(testMaxIdleTimeClosed))
}
