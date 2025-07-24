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
		title = page.Replace('(.)([A-Z])', '\1 \2')
		return '<html><head><title>' $
			(mode is 'edit' ? 'Edit ' : 'Add to ') $ title $
			'</title>
			<script language="JavaScript">
				var needToConfirm = true;
				window.onbeforeunload = confirmExit;
				function confirmExit()
					{
					if (needToConfirm)
						return "Please use Save or Cancel ' $
							'or you will leave the page locked.";
					}
			</script>
			</head>
			<body>
			<form method="post" action="Wiki?' $ page $ '">
			<h1>' $
			(mode is 'edit' ? 'Edit ' : 'Add to ') $ title $
			' <input type="submit" name="Save" value="Save"
				onclick="needToConfirm = false;">
			<input type="submit" name="Cancel" value="Cancel"
				onclick="needToConfirm = false;">
			<input type="hidden" name="editmode" value="' $ mode $ '">
			</h1>
			<textarea name="text" style="width:100%;height:85%" rows=20 ' $
				'wrap="virtual" style="font-size:10pt">' $
			text $
			'</textarea><p>
			</form>
			</body></html>'
		}
	}