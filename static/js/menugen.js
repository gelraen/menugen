function replaceContent(elem, uri) {
	$.get(uri,
		function(data) {
			elem.innerHTML = data
		})
}

function genBreakfast(t) {
	replaceContent(document.getElementById(t), "/gen/breakfast")
}

function genLunch(t, day) {
	replaceContent(document.getElementById(t), "/gen/lunch?day=" + day)
}

function genDinner(t) {
	replaceContent(document.getElementById(t), "/gen/dinner")
}
