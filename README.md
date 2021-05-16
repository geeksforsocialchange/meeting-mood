# Meeting Mood

![screenshot.png](screenshot.png)

Consensus hand signals for online meetings

## Installation:

To install, either:
* Download the appropriate binary from the [latest release](https://github.com/WheresAlice/meeting-mood/releases/latest)
* Run `brew install wheresalice/meeting-mood/meeting-mood` to install via Homebrew
* Clone this repo and run `go build`

Run the `meeting-mood` binary file somewhere that all your meeting joiners can connect to, which usually means somewhere with a public IP.  By default, it will listen on port 8844.

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

## Known Issues

- HTTPS through Ngrok doesn't work
- There is no authentication or authorisation
- Hardcoded set of moods (for consensus decision making)
- No sorting of moods
- Not very mobile friendly
- Probably not very performant in places

## Donations

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/wheresalice)

## License

meeting-mood is released under the [MIT license](LICENSE).