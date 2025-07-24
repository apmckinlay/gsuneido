// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (html)
	{
	blob = SuUI.MakeWebObject('Blob', [html], [type: 'text/html'])
	clipboardItem = SuUI.MakeWebObject('ClipboardItem', ['text/html': blob])

	// Support for multiple ClipboardItems is not implemented by Chrome as of 09/10/2021
	SuUI.GetCurrentWindow().navigator.clipboard.Write(
		[clipboardItem])
	}
