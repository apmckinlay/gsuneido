// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// hwnd argument is used to center on parent
function (msg, title = "ALERT", _hwnd = 0, flags = 0)
	{
	msg = TranslateLanguage(msg)
	title = TranslateLanguage(title)
	// Suneido.Alert optionally overrides AlertMessageBox
	// Use Global so this record is compilable on suneido.js and
	// also avoid loading too many server Controller related code to browser
	alert = Suneido.GetDefault(#Alert, { Global('AlertMessageBox') })
	msg = String(msg).LeftTrim().Ellipsis(10000 /*=maxMsgSize*/, atEnd:).RightTrim()
	return alert(hwnd, msg, title, flags)
	}
