# Meeting Mood

![screenshot.png](screenshot.png)

Consensus hand signals for online meetings

To install, either:
* Download the appropriate binary from the [latest release](https://github.com/WheresAlice/meeting-mood/releases/latest)
* Run `brew install wheresalice/meeting-mood/meeting-mood` to install via Homebrew
* Clone this repo and run `go build`

Then run it with `meeting-mood` and open http://localhost:8844

Specify the environment variable `PORT=8080` or pass the flag `--port 8080` to override the port number

Put it on the internet and share the link with everybody in the meeting (use [Ngrok](https://ngrok.com/) and plain HTTP)

## Known Issues

- HTTPS through Ngrok doesn't work
- Only a single meeting is supported
- There is no authentication or authorisation
- Hardcoded set of moods (for consensus decision making)
- No sorting of moods
- Not very mobile friendly