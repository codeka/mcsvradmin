// +build dev
//go:generate go run -tags=dev static_generate.go

package static

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("./static/files")
