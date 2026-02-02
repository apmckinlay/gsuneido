// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
function (component, pasteCallable)
	{
	SuUI.GetCurrentWindow().navigator.clipboard.ReadText().Then({ |s|
		if component.HasFocus?() and component.GetReadOnly() isnt true
			pasteCallable(s)
		}).Catch({ |err| SuClipboardHandleErr(err, 'Paste') })
	}
