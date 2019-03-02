var username = "";
var connectedUsers = [];
var ws = new WebSocket('ws://' + window.location.host + '/chat/ws');

ws.addEventListener('message', function(e) {

    var msg = JSON.parse(e.data);

    if(msg.hasOwnProperty("message")) {

        appendMessage(msg);
        scrollToButtom();
        playSound("/sounds/message");

        if(document.hidden) {

            notify(msg.username, msg.message);
        }
    }
    else if(msg.hasOwnProperty("messages")) {

        for(var i = 0; i < msg.messages.length; i++) {

            appendMessage(msg.messages[i]);
        }

        scrollToButtom();
    }
    else {

        if(msg.hasOwnProperty("receiver")) {

            username = msg.receiver;
        }

        if(msg.hasOwnProperty("connections")) {

            connectedUsers = msg.connections;
            document.getElementById('users_list').innerHTML = '';

            for (i = 0; i < connectedUsers.length; i++) {

                var users = document.getElementById('users_list')
                var user = document.createElement('li');
                var href = document.createElement('a');

                href.appendChild(document.createTextNode(connectedUsers[i]));
                href.setAttribute("href", "#");

                user.appendChild(href);
                users.appendChild(user);
            }
        }
    }

});

var message = document.getElementById('message');
message.addEventListener('keyup', function(e) {
    if(e.which === 13) {
        send();
        e.preventDefault();
    }
});

function send() {

    var newMsg = document.getElementById('message').value;

    if (newMsg !== '') {
        ws.send(
            JSON.stringify({
                    message: newMsg // Strip out html
                }
            ));

        document.getElementById('message').value = ''; // Reset newMsg
    }
}

function appendMessage(msg) {

    var user_chip = '';

    if(msg.username === username) {

        user_chip = '<p class="user_chip_self">';
    }
    else {

        user_chip = '<p class="user_chip">';
    }

    var content = '<div class="message">'
        + user_chip
        + msg.username
        + '</p> '
        + '<div class="talk-bubble tri-right round left-in">'
        + '<p class="talktext">'
        + msg.message
        + '</p> '
        + '</div>'
        + '</div>'; // Parse emojis

    var chat = document.getElementById('chat-messages')
    chat.innerHTML += content
}

function scrollToButtom() {

    var element = document.getElementById('chat-messages');
    element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
}

function playSound(filename){

    var mp3Source = '<source src="' + filename + '.mp3" type="audio/mpeg">';
    var oggSource = '<source src="' + filename + '.ogg" type="audio/ogg">';
    var embedSource = '<embed hidden="true" autostart="true" loop="false" src="' + filename +'.mp3">';
    document.getElementById("sound").innerHTML='<audio autoplay="autoplay">' + mp3Source + oggSource + embedSource + '</audio>';
}

function notify(username, message) {

    if (("Notification" in window)) {

        if(Notification.permission === "granted") {

            // If it's okay let's create a notification
            var notification = new Notification("[" + username + "] " + message);
            setTimeout(notification.close.bind(notification), 4000);
        }
        else if (Notification.permission !== "denied") {

            Notification.requestPermission().then(function (permission) {
                // If the user accepts, let's create a notification
                if (permission === "granted") {
                    var notification = new Notification("[" + username + "] " + message);
                    setTimeout(notification.close.bind(notification), 6000);
                }
            });
        }
    }
}

window.onunload = function(){};