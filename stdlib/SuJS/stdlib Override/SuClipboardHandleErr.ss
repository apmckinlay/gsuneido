// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
function (err, action)
	{
	pattern = Object(`request is not allowed by the user agent`, // common on Safari
		`method not found: Clipboard.ReadText`)
	shortcuts = #(
		Cut: 'Ctrl+X',
		Copy: 'Ctrl+C',
		Paste: 'Ctrl+V',)
	if err.message =~ pattern.Join(' | ')
		{
		msg = 'Right-click > ' $ action $ ' is not supported in this ' $
			'web browser.\r\nPlease use ' $ shortcuts[action] $ ' instead.'
			SuRender().Event(false, 'Alert', [msg, 'Unsupported Browser Feature',
				flags: MB.ICONINFORMATION])
		return false
		}
	return true
	}