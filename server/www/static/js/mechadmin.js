function make_computer_table() {
	$('#content').empty();
	$('#content').append(
		$('<table>').attr("id", "comp-table").append(
			$('<tr>').append(
				$('<th>').append("Computer Name")).append(
				$('<th>').append("Active User")).append(
				$('<th>').append("Status")).append(
				$('<th>').append("IP Address"))));
}

function handle_message(msg) {
	if (msg.R = "list-compR") {
		make_computer_table();
		for (i=0; i < msg.D.length; i++)  {
			$('#comp-table').append(
				$('<tr>').attr("id", "comp-"+msg.D[i].MAC).append(
					$('<td>').append(msg.D[i].CN)));
		}
	}
	else if (evt.data === "bye") {
		websocket.close();
	}
}

function requestCol(e) {
	var req = {R: 'list-comp',A1: e.id};
	doSend(JSON.stringify(req));
}
