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
