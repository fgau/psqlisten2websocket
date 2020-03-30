$(function () {
    var ws;

    if (window.WebSocket === undefined) {
        $("#container").append("your browser does not support webSockets");
        return;
    } else {
        ws = initWS();
    }

    function initWS() {
        var socket = new WebSocket("ws://localhost:8080/ws");
        var container = $("#container");
        socket.onopen = function () {
            container.append("<p>socket is open</p>");
        };
        socket.onmessage = function (e) {
            console.log("DATA: ", (e.data.replace(/[\n\t]/g, '')));

            // converting JSON object to JS object, regex tabs and linebreaks
            var obj = JSON.parse(e.data.replace(/[\n\t]/g, ''));

            container.prepend("<p> got some shit1:" + obj.data.id + ", " + obj.data.url + "</p>");
        }
        socket.onclose = function () {
            setTimeout(initWS, 1000);
            container.append("<p>socket closed</p>");
        }

        return socket;
    }

});
