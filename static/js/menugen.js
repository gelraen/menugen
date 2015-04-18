function replaceContent(elem, uri) {
	$.get(uri,
		function(data) {
			elem.innerHTML = data
		})
}

function genBreakfast(t) {
	replaceContent(document.getElementById(t), "/gen/breakfast")
}

function genDinner(t) {
	replaceContent(document.getElementById(t), "/gen/dinner")
}
