package tilt

import (
	"context"
	"testing"

	mock_ble "github.com/fetherolfjd/leaning-hydro-pi/test/mock"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/golang/mock/gomock"
)

func TestConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newBtDevice = func(opts ...ble.Option) (*linux.Device, error) {
		return &linux.Device{HCI: nil, Server: nil}, nil
	}
	setBtDevice = func(d ble.Device) {}
	btScan = func(ctx context.Context, allowDup bool, advHandler ble.AdvHandler, advFilter ble.AdvFilter) error {
		return nil
	}

	ctx := context.Background()
	tdpChan := make(chan *TiltDataPoint)
	defer close(tdpChan)
	testErr := Connect(ctx, tdpChan)
	if testErr != nil {
		t.Fatalf("Unexpected error from Connect: %v", testErr)
	}
}

func TestAdvertisementHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantAddr := "OMGWTF"

	mockAddr := mock_ble.NewMockAddr(ctrl)
	mockAddr.EXPECT().String().Return(wantAddr)
	mockAdv := mock_ble.NewMockAdvertisement(ctrl)
	mockAdv.EXPECT().Addr().Return(mockAddr)
	mockAdv.EXPECT().RSSI().Return(42)
	mockAdv.EXPECT().ManufacturerData().Return(getTestData())

	dataPointChan = make(chan *TiltDataPoint, 1)
	defer close(dataPointChan)

	advHandler(mockAdv)

	tdp := <-dataPointChan
	if tdp.Address != wantAddr {
		t.Fatalf("Address %s does not match expected %s", tdp.Address, wantAddr)
	}
	wantColor := "RED"
	if tdp.Color != wantColor {
		t.Fatalf("Color %s does not match expected %s", tdp.Color, wantColor)
	}
	wantRSSI := 42
	if tdp.RSSI != wantRSSI {
		t.Fatalf("RSSI %d does not match expected %d", tdp.RSSI, wantRSSI)
	}
	wantSG := float32(1.016)
	if tdp.SpecificGravity != wantSG {
		t.Fatalf("Specific gravity %f does not match expected %f", tdp.SpecificGravity, wantSG)
	}
	wantTx := 197
	if tdp.TXPower != wantTx {
		t.Fatalf("Transmit power %d does not match expected %d", tdp.TXPower, wantTx)
	}
	wantTemp := float32(68.0)
	if tdp.Temperature != wantTemp {
		t.Fatalf("Temperature %f does not match expected %f", tdp.Temperature, wantTemp)
	}
	wantUUID := "A495BB10C5B14B44B5121370F02D74DE"
	if tdp.UUID != wantUUID {
		t.Fatalf("UUID %s does not match expected %s", tdp.UUID, wantUUID)
	}
}

func TestAdvertisementFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdv := mock_ble.NewMockAdvertisement(ctrl)
	mockAdv.EXPECT().ManufacturerData().Return(getTestData())
	got := advFilter(mockAdv)
	if got != true {
		t.Fatalf("Advertisement filter failed to find correct address")
	}

	mockAdv.EXPECT().ManufacturerData().Return(getBadTestData())
	got = advFilter(mockAdv)
	if got != false {
		t.Fatalf("Advertisement filter failed to find correct address")
	}
}
