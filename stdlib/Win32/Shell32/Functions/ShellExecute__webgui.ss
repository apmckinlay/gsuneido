// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
function(hwnd/*unused*/, lpVerb, lpFile, lpParameters/*unused*/ = '',
	lpDirectory/*unused*/ = '', nShowCmd/*unused*/ = false, fMask/*unused*/ = 0)
	{
	switch
		{
	case lpVerb is 'open':
		SuRenderBackend().RecordAction(false, 'SuJsExecute', [#open, lpFile])
	case lpVerb is NULL and Paths.IsValid?(lpFile):
		url = AjaxBrowserControl.Server() $ 'attachment'
		queryOb = Object(Base64.Encode(lpFile.Xor(EncryptControlKey())), preview:,
			token: SuRenderBackend().Token)
		SuRenderBackend().RecordAction(false, 'SuJsExecute',
			[#open, Url.Encode(url, queryOb)])
	default:
		SuServerPrint("unsupported excute verb: " $ lpVerb)
		return false
		}
	return true
	}