// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (str, action = 'Copy')
	{
	SuUI.GetCurrentWindow().navigator.clipboard.WriteText(str).Then({}).Catch({ |err|
		SuClipboardHandleErr(err, action)
		})
	}
