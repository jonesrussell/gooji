package handler

import (
	"github.com/romsar/gonertia"
)

const devTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gooji</title>
    <script type="module" src="http://localhost:5173/@vite/client"></script>
    <script type="module" src="http://localhost:5173/src/main.ts"></script>
</head>
<body>
    <div id="app" data-page="{{ .InertiaPage }}"></div>
</body>
</html>`

const prodTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gooji</title>
    <script type="module" src="/assets/main.js"></script>
    <link rel="stylesheet" href="/assets/main.css">
</head>
<body>
    <div id="app" data-page="{{ .InertiaPage }}"></div>
</body>
</html>`

func NewInertia(isDev bool) (*gonertia.Inertia, error) {
	if isDev {
		return gonertia.New(devTemplate)
	}
	return gonertia.New(prodTemplate)
}
