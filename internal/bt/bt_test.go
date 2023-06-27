package bt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/model"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScan(t *testing.T) {
	t.Run("errors if new device fails", func(t *testing.T) {
		oldNewDev := newBtDevice
		newBtDevice = func(opts ...ble.Option) (*linux.Device, error) {
			return nil, errors.New("test error")
		}
		defer func() {
			newBtDevice = oldNewDev
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := Scan(ctx, nil, nil)
		assert.ErrorContains(t, err, "test error")
	})

	t.Run("calls filter and handler", func(t *testing.T) {
		var newDevCalled bool
		var setDevCalled bool
		oldNewDev := newBtDevice
		newBtDevice = func(opts ...ble.Option) (*linux.Device, error) {
			newDevCalled = true
			return &linux.Device{}, nil
		}
		oldSetDev := setBtDevice
		setBtDevice = func(_ ble.Device) {
			setDevCalled = true
		}
		oldScan := scan
		addr := &mockAddr{}
		addr.On("String").Return("myaddr")
		adv := &mockAdv{}
		adv.On("Addr").Return(addr).
			On("RSSI").Return(42).
			On("ManufacturerData").Return([]byte{})

		scan = func(_ context.Context, _ bool, h ble.AdvHandler, f ble.AdvFilter) error {
			f(adv)
			h(adv)
			return nil
		}
		defer func() {
			newBtDevice = oldNewDev
			setBtDevice = oldSetDev
			scan = oldScan
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		hdlr := &stubHdlr{}
		filter := &stubFilter{}

		err := Scan(ctx, hdlr, filter)
		assert.NoError(t, err)
		assert.True(t, newDevCalled)
		assert.True(t, setDevCalled)
		assert.True(t, hdlr.called)
		assert.True(t, filter.called)
	})

}

type stubHdlr struct {
	called bool
}

func (h *stubHdlr) Handle(_ model.DataPacket) {
	h.called = true
}

type stubFilter struct {
	called bool
}

func (f *stubFilter) Filter(_ model.DataPacket) bool {
	f.called = true
	return true
}

type mockAdv struct {
	mock.Mock
}

func (m *mockAdv) LocalName() string {
	return m.Called().String(0)
}

func (m *mockAdv) ManufacturerData() []byte {
	return m.Called().Get(0).([]byte)
}

func (m *mockAdv) ServiceData() []ble.ServiceData {
	return m.Called().Get(0).([]ble.ServiceData)
}

func (m *mockAdv) Services() []ble.UUID {
	return m.Called().Get(0).([]ble.UUID)
}

func (m *mockAdv) OverflowService() []ble.UUID {
	return m.Called().Get(0).([]ble.UUID)
}

func (m *mockAdv) TxPowerLevel() int {
	return m.Called().Int(0)
}

func (m *mockAdv) Connectable() bool {
	return m.Called().Bool(0)
}

func (m *mockAdv) SolicitedService() []ble.UUID {
	return m.Called().Get(0).([]ble.UUID)
}

func (m *mockAdv) RSSI() int {
	return m.Called().Int(0)
}

func (m *mockAdv) Addr() ble.Addr {
	return m.Called().Get(0).(ble.Addr)
}

type mockAddr struct {
	mock.Mock
}

func (m *mockAdv) String() string {
	return m.Called().String(0)
}
