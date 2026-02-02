// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(mode, page, remote_user = '', table = 'wiki')
		{
		if '' isnt valid = WikiTitleValid?(page)
			return valid
		if true isnt locker = WikiLock(page, remote_user)
			return .pageAlreadyLocked(page, locker)

		text = .getPageText(mode, page, table)
		return .pageLocked(mode, page, text)
		}

	pageAlreadyLocked(page, locker)
		{
		return '<html><head><title>Page Locked</title></head>
			<body>
			<form method="post" action="Wiki?unlock=' $ page $ '">
			<div align=center><p>&nbsp;</p>
			<p>' $ locker.user $
				' started editing this page at ' $ locker.date.ShortDateTime() $ '</p>
			<input type="submit" value="Unlock">
			</div>
			</form></body></html>'
		}

	getPageText(mode, page, table)
		{
		text = ''
		if mode is 'edit' and false isnt x = Query1(table, name: page)
			text = x.text
		return text
		}

	pageLocked(mode, page, text)
		{
		title = (mode is 'edit' ? 'Edit ' : 'Add to ') $
			page.Replace('(.)([A-Z])', '\1 \2')
		return Razor(.template, Object(:page, :title, :mode, set: Base64.Encode(text)))
		}

	template: `
<html>
	<head>
		<title>@.title</title>
		<script type=text/javascript src="/CodeRes?name=overtype_min_2_1_1.js"></script>
		<script language="JavaScript">
			let needToConfirm = true;
			let content = atob("@.set");
			let init = false;
			window.onbeforeunload = confirmExit;
			function confirmExit() {
				if (needToConfirm)
					return "Please use Save or Cancel or you will leave the page locked.";
			}
			window.onload = onReady;
			function onReady() {
				const preview = document.getElementById('preview');
				disablePreviewLink();
				let debounceTimer;
				function disablePreviewLink() {
					const doc = preview.contentDocument || preview.contentWindow.document;
					doc.open();
					doc.write('<!DOCTYPE html>' +
'<html><head><style>a { pointer-events: none; }</style></head><body></body></html>');
					doc.close();
				}
				function updatePreview(s) {
					const doc = preview.contentDocument || preview.contentWindow.document;
					const el = doc.documentElement;

					const isAtBottom = el.scrollTop > 0 &&
						(el.scrollHeight - el.clientHeight - el.scrollTop) <= 1;
					const scrollPos = el.scrollTop;
					fetch('/Wiki?preview', {
						method: 'POST',
						body: s
					})
					.then(response => {
						if (!response.ok)
							throw new Error('Network response was not ok');
						return response.text();
					})
					.then(htmlString => {
						doc.body.innerHTML = htmlString;
						requestAnimationFrame(() => {
							el.scrollTop = isAtBottom ? el.scrollHeight : scrollPos;
						});
					});
				}
				function onChange(value) {
					if (!init || value === content) {
						return;
					}
					content = value;
					clearTimeout(debounceTimer);
					debounceTimer = setTimeout(() => {
						updatePreview(value)
					}, 500);
				}

				const [editor] = new OverType("#editor", {
					toolbar: true,
					toolbarButtons: [
						toolbarButtons.bold,
						toolbarButtons.italic,
						toolbarButtons.code,
						toolbarButtons.separator,
						toolbarButtons.link,
						toolbarButtons.separator,
						toolbarButtons.h1,
						toolbarButtons.h2,
						toolbarButtons.h3,
						toolbarButtons.separator,
						toolbarButtons.bulletList,
						toolbarButtons.orderedList,
						toolbarButtons.separator,
						toolbarButtons.quote],
					onChange: onChange});
				editor.setValue(content);
				init = true;
				updatePreview(content);

				const form = document.querySelector('form');
				form.addEventListener('submit', (event) => {
					const hiddenInput = document.createElement('input');
					hiddenInput.type = 'hidden';
					hiddenInput.name = 'text';
					hiddenInput.value = editor.getValue();
					form.appendChild(hiddenInput);
				});
			}
		</script>
		<style>
			#layout {
				display: flex;
				flex-direction: column;
				height: 100%;
				width: 100%;
			}
			#container {
				flex-grow: 1;
				display: flex;
				flex-direction: row;
			}
			#editor {
				flex: 1;
				font-size:10pt;
				border: 1px solid black;
			}
			#preview {
				flex: 1;
				border: 1px solid black;
			}
		</style>
	</head>
	<body>
		<form method="post" action="Wiki?@.page">
			<div id="layout">
				<h1>@.title
					<input type="submit" name="Save" value="Save"
						onclick="needToConfirm = false;">
					<input type="submit" name="Cancel" value="Cancel"
						onclick="needToConfirm = false;">
					<input type="hidden" name="editmode" value="@.mode">
				</h1>
				<div id="container">
					<div id="editor"></div>
					<iframe id="preview"></iframe>
				</div>
			</div>
		</form>
	</body>
</html>
`
	}
