var wsUri = "ws://192.168.1.203:8080/ws/";
var output;
var websocket;

function writeToScreen(message) {
	var pre = document.createElement("p");
	pre.style.wordWrap = "break-word";
	pre.innerHTML = message;
	output.appendChild(pre);
}

function doSend(message) {
	writeToScreen("SENT: " + message);
	websocket.send(message);
}

function onOpen(evt) {
	writeToScreen("CONNECTED");
}

function onClose(evt) {
	writeToScreen("DISCONNECTED");
}

function onMessage(evt) {
	writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data + '</span>');
	if (evt.data === "bye") {
		websocket.close();
	}
}

function onError(evt) {
	writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
}

function testit() {
	doSend(document.getElementById("in").value);
}

function bindSocket() {
	websocket = new WebSocket(wsUri);
	websocket.onopen = function (evt) { onOpen(evt); };
	websocket.onclose = function (evt) { onClose(evt); };
	websocket.onmessage = function (evt) { onMessage(evt); };
	websocket.onerror = function (evt) { onError(evt); };
}

function init() {
	output = document.getElementById("content");
	bindSocket();
}

function requestCol(e) {
    var req = {
        R: 'list-comp',
        A1: e.id
    };
    doSend(JSON.stringify(req));
}

window.addEventListener("load", init, false);
