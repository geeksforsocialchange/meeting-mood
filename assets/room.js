const app = new Vue({
    el: '#app',
    data: {
        user: '',
        room: undefined,
        moods: [
            "Agree",
            "Agree and volunteer",
            "Response",
            "Direct response",
            "Technical point",
            "Language",
            "Speak up",
            "Slow down",
            "I'm confused",
            "Veto"
        ],
        userMoods: []
    },
    computed: {
        isReady: function () {
            return this.user.length > 0
        }
    },
    methods: {
        deleteUserMood: function (username = this.user) {
            const body = JSON.stringify({
                username: username,
                roomUser: `${this.room}${username}`
            })
            fetch(`/${this.room}/delete`, {
                method: 'post',
                body: body
            })
        },
        sendMood: function (mood) {
            // 2419200 is 28 days - the username cookie will expire if you don't use it for four weeks
            document.cookie = `username=${this.user};max-age=2419200;SameSite=Strict`
            const message = JSON.stringify({
                username: this.user,
                mood: mood,
                room: this.room,
                roomUser: `${this.room}${this.user}`
            })
            fetch(`/${this.room}/mood`, {
                method: 'post',
                body: message
            })
        },
    }
})

window.onload = function () {
    try {
        app.user = document.cookie.split("; ").find(row => row.startsWith("username=")).split('=')[1]
    } catch (e) {
        console.log(e)
        console.log("New user")
    }
    app.room = location.pathname.split("/")[1]


    const socketAddr = ((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + window.location.pathname + "/ws"
    const moodSocket = new ReconnectingWebSocket(socketAddr);
    const update = function () {
        moodSocket.onopen = function () {
            console.log(`(re)-connected to ${socketAddr}`)
            app.userMoods = []
            fetch(`/${app.room}/all`)
                .then(res => res.json())
                .then((all) => {
                    if (all != null) {
                        all.forEach(mood => {
                            app.userMoods.push(mood)
                        })
                    }
                })
        }
        moodSocket.onmessage = function (event) {
            const moodOperation = JSON.parse(event.data)
            const mood = {username: moodOperation.username, mood: moodOperation.mood}
            switch (moodOperation.operation) {
                case "Delete":
                    // delete
                    app.userMoods = app.userMoods.filter(item => item.username != mood.username)
                    break;
                case "Save":
                    app.userMoods = app.userMoods.filter(item => item.username != mood.username)
                    app.userMoods.push(mood)
            }
        }
    };
    window.setTimeout(update);
}