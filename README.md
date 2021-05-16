# Meeting Mood

![screenshot.png](screenshot.png)

Consensus hand signals for online meetings

Run it with `go run *.go` or execute one of the release binaries and then open http://localhost:8844

Put it on the internet and share the link with everybody in the meeting (use [Ngrok](https://ngrok.com/) and plain HTTP)

## Known Issues

- Port is hardcoded
- HTTPS through Ngrok doesn't work
- Only a single meeting is supported
- There is no authentication or authorisation
- Hardcoded set of moods (for consensus decision making)
- No sorting of moods
- Not very mobile friendly