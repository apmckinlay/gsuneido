// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
function (component, pasteCallable)
	{
	try
		SuUI.GetCurrentWindow().navigator.clipboard.ReadText().Then({ |s|
			if component.HasFocus?() and component.GetReadOnly() isnt true
				pasteCallable(s)
			})
	catch (unused, 'method not found: Clipboard.ReadText')
		component.UnsupportedFeature(
			'Right-click > Paste is not supported in this ' $
			'web browser.\r\nPlease use Ctrl+V instead.')
	}
