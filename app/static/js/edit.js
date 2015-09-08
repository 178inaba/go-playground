$(document).ready(function() {
	playground({
		'codeEl':       '#code',
		'outputEl':     '#output',
		'runEl':        '#run',
		'fmtEl':        '#fmt',
		'fmtImportEl':  '#imports',
		'shareEl':      '#share',
		'shareURLEl':   '#shareURL',
		'enableHistory': true
	});
	playgroundEmbed({
		'codeEl':       '#code',
		'shareEl':      '#share',
		'embedEl':      '#embed',
		'embedLabelEl': '#embedLabel',
		'embedHTMLEl':  '#shareURL'
	});
	$('#code').linedtextarea();

	// login
	$('#loginButton').click(function() {
		location.href = $(this).attr("url");
	});

	// open post gist modal
	$("#postGistModalButton").click(function() {
		if (typeof Cookies.get("access_token") === "undefined") {
			$("#loginAlertModal").modal("show");
		} else {
			$("#postGistModal").modal("show");
		}
	});

	// create gist
	$("#postGistButton").click(function() {
		$.ajax({
			url: "/gist",
			type: "POST",
			data: $('#code').val()
		});
	});

	// Avoid line wrapping.
	$('#code').attr('wrap', 'off');
	var about = $('#about');
	about.click(function(e) {
		if ($(e.target).is('a')) {
			return;
		}
		about.hide();
	});
	$('#aboutButton').click(function() {
		if (about.is(':visible')) {
			about.hide();
			return;
		}
		about.show();
	});
	// Preserve "Imports" checkbox value between sessions.
	if (readCookie('playgroundImports') == 'true')
		$('#imports').attr('checked','checked');
	$('#imports').change(function() {
		createCookie('playgroundImports', $(this).is(':checked') ? 'true' : '');
	});
	// Fire Google Analytics events for Run/Share button clicks.
	if (window.trackEvent) {
		$('#run').click(function() {
			window.trackEvent('playground', 'click', 'run-button');
		});
		$('#share').click(function() {
			window.trackEvent('playground', 'click', 'share-button');
		});
	}
});

function createCookie(name, value) {
	document.cookie = name+"="+value+"; path=/";
}

function readCookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for(var i=0;i < ca.length;i++) {
		var c = ca[i];
		while (c.charAt(0)==' ') c = c.substring(1,c.length);
		if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
	}
	return null;
}
