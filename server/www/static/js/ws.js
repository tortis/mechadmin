var wsUri = "ws://192.168.1.203:8080/ws/";
var output;
var websocket;

function writeToScreen(message) {
	$("#content").append("<p>"+message+"</p>")
}

function doSend(message) {
	websocket.send(message);
}

function onOpen(evt) {
	writeToScreen("CONNECTED");
}

function onClose(evt) {
	writeToScreen("DISCONNECTED");
}

function onMessage(evt) {
	var msg = JSON.parse(evt.data);
	console.log(msg.R);
	handle_message(msg);
}

function onError(evt) {
	writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
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

window.addEventListener("load", init, false);
