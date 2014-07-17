var wsUri = "ws://192.168.1.203:8080/ws/";
var output;
var websocket;

function writeToScreen(message) {
	$("#content").append("<p>"+message+"</p>")
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
	console.log(evt.data);
	var msg = JSON.parse(evt.data);
	if (msg.R = "list-compR") {
		make_computer_table();
		for (i=0; i < msg.D.length; i++)  {
			$('#comp-table').append(
				$('<tr>').attr("id", "comp-"+msg.D[i].MAC).append(
					$('<td>').append(msg.D[i].CN)));
		}
	}
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
