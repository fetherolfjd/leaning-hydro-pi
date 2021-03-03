module leaning-hydro-pi

go 1.15

require github.com/gorilla/mux v1.8.0

require (
	github.com/go-ble/ble v0.0.0-20200407180624-067514cd6e24
	github.com/gorilla/websocket v1.4.2
	internal v0.0.0
)

replace internal => ./internal
