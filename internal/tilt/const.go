package tilt

type tiltUUID = string
type tiltColor = string

const (
	uuidRed    tiltUUID = "A495BB10C5B14B44B5121370F02D74DE"
	uuidGreen  tiltUUID = "A495BB20C5B14B44B5121370F02D74DE"
	uuidBlack  tiltUUID = "A495BB30C5B14B44B5121370F02D74DE"
	uuidPurple tiltUUID = "A495BB40C5B14B44B5121370F02D74DE"
	uuidOrange tiltUUID = "A495BB50C5B14B44B5121370F02D74DE"
	uuidBlue   tiltUUID = "A495BB60C5B14B44B5121370F02D74DE"
	uuidYellow tiltUUID = "A495BB70C5B14B44B5121370F02D74DE"
	uuidPink   tiltUUID = "A495BB80C5B14B44B5121370F02D74DE"

	colorRed    tiltColor = "RED"
	colorGreen  tiltColor = "GREEN"
	colorBlack  tiltColor = "BLACK"
	colorPurple tiltColor = "PURPLE"
	colorOrange tiltColor = "ORANGE"
	colorBlue   tiltColor = "BLUE"
	colorYellow tiltColor = "YELLOW"
	colorPink   tiltColor = "PINK"

	tempStartByte       int = 20
	tempEndByte         int = 22
	sgStartByte         int = 22
	sgEndByte           int = 24
	transmitDataByte    int = 24
	deviceUUIDStartByte int = 4
	deviceUUIDEndByte   int = 20
)

var uuidToColor = map[tiltUUID]tiltColor{
	uuidRed:    colorRed,
	uuidGreen:  colorGreen,
	uuidBlack:  colorBlack,
	uuidPurple: colorPurple,
	uuidOrange: colorOrange,
	uuidBlue:   colorBlue,
	uuidYellow: colorYellow,
	uuidPink:   colorPink,
}

var colorToUuid = map[tiltColor]tiltUUID{
	colorRed:    uuidRed,
	colorGreen:  uuidGreen,
	colorBlack:  uuidBlack,
	colorPurple: uuidPurple,
	colorOrange: uuidOrange,
	colorBlue:   uuidBlue,
	colorYellow: uuidYellow,
	colorPink:   uuidPink,
}
