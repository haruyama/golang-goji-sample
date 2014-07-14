# Golang Goji Sample

Sample Web Application based on https://github.com/elcct/defaultproject

- Goji - A web microframework for Golang - http://goji.io/
- Gorilla web toolkit sessions - cookie (and filesystem) sessions - http://www.gorillatoolkit.org/pkg/sessions

# Dependencies

Default Project requires `Go`, `MySQL` and few other tools installed.

# Project structure

`/controllers`

All your controllers that serve defined routes.

`/helpers`

Helper functions.

`/models`

You database models.

`/public`

It has all your static files mapped to `/assets/*` path except `robots.txt` and `favicon.ico` that map to `/`.

`/system`

Core functions and structs.

`/views`

Your views using standard `Go` template system.

`server.go`

This file starts your web application and also contains routes definition.
