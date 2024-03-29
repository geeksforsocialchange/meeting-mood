# Meeting Mood

[![Release](https://img.shields.io/github/release-pre/geeksforsocialchange/meeting-mood.svg?logo=github&style=flat&v=1)](https://github.com/geeksforsocialchange/meeting-mood/releases)
[![Build Status](https://img.shields.io/github/workflow/status/geeksforsocialchange/meeting-mood/run-go-tests?logo=github&v=1)](https://github.com/geeksforsocialchange/meeting-mood/actions)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://gh.mergify.io/badges/geeksforsocialchange/meeting-mood&style=flat&v=1)](https://mergify.io)
[![Go](https://img.shields.io/github/go-mod/go-version/geeksforsocialchange/meeting-mood?v=1)](https://golang.org/)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/geeksforsocialchange/meeting-mood)

![screenshot.png](screenshot.png)

Consensus hand signals for online meetings

## Installation:

This application needs to be run somewhere with a public IP so that people can connect to it.  There is a [development version](https://meeting-mood.wheresalice.info) running but it may break at any point.

The easiest way of doing this is to deploy this repository as a Dokku app.  We presume that you can also push it as a Heroku app, but that is untested.

There are a number of other ways you can get hold the binary to run:

* Download the appropriate binary from the [latest release](https://github.com/geeksforsocialchange/meeting-mood/releases/latest)
* Clone this repo and run `go build`

If you are running the binary outside of Dokku/Heroku then it will listen on port 8844 by default.

## Configuration Options

You can optionally:

* Specify the environment variable `PORT=8080` or pass the flag `--port 8080` to override the port number
* Put a footer.html in your current directory, and the contents will be displayed as a footer.

## Usage

1. Open a web browser to the web server and you'll see a button to create a room.  This will create a room and put you in it.
2. Share the address in the address bar with other members of the meeting.
3. Set your username (this will be saved in a cookie)
4. Press buttons to make hand signals
5. Press `x` to stop making a hand signal

## Development

For purely-frontend development we provide a handy docker-compose environment which will let you make changes to the html without needing Go.

Running `docker-compose up` will launch a Caddy proxy listening on http://localhost:8800/ as well as building and launching the checked-out Go code.  You can then make changes to index.html, room.html, and anything in the assets directory and refresh the browser without needing to rebuild any of the backend components.

## Donations

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/wheresalice)

## License

meeting-mood is released under the [MIT license](LICENSE).
