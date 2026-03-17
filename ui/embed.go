package ui

import "embed"

//go:embed index.html app.js graph.js style.css vendor/vis-network.min.js
var Assets embed.FS
