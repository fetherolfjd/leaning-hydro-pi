package publish

import (
	"io"

	"github.com/fetherolfjd/leaning-hydro-pi/internal/tilt"
)

type DatapointPublisher interface {
	io.Closer
	Publish(tdp *tilt.TiltDataPoint)
	PublishAll(tdp []*tilt.TiltDataPoint)
}
