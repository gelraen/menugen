function replaceContent(elem, uri) {
	$.get(uri,
		function(data) {
			elem.innerHTML = data
		})
}

function genBreakfast(t) {
	replaceContent(document.getElementById(t), "/gen/breakfast")
}

function genLunch(t) {
	replaceContent(document.getElementById(t), "/gen/lunch")
}

function genDinner(t) {
	replaceContent(document.getElementById(t), "/gen/dinner")
}
