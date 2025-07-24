// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// contributions by Claudio Mascioni
// BUG should free pidl with CoTaskMemFree
function (hwnd = 0, title = 'Browse for Folder', flags = "",
	pidlRoot = 0, callbackfn = 0, returnpidl = false, initialPath = false)
	{
	if hwnd is 0
		try hwnd = _hwnd
	title = TranslateLanguage(title)

	if initialPath is false
		callBack = callbackfn
	else
		callBack = { |hWnd, uMsg, lp /*unused*/, pData /*unused*/|
			if uMsg is BFFM.INITIALIZED
				{
				// to set initial path, not tested with network path
				// path without last '\': yes C:\Temp, no C:\Temp\
				initialPath = initialPath.RightTrim('/\\')
				// no dir $= '\x00', seems that no require a '\x00' as last character
				// if dir is '' or is an invalid path, select Computer resources
				SendMessageText(hWnd, WM.USER + BFFM.SETSELECTION, 1, initialPath)
				}
			0 // block return
			}
	bi = Object(hwndOwner: hwnd,
		:pidlRoot,
		pszDisplayName: 0,
		lpszTitle: title,
		ulFlags: flags is ""
			? BIF.USENEWUI | BIF.RETURNONLYFSDIRS | BIF.STATUSTEXT
			: flags,
		lpfn: callBack,
		)
	pidl = 0
	CenterDialog(hwnd)
		{ pidl = SHBrowseForFolder(bi) }

	if initialPath isnt false
		ClearCallback(callBack)
	if returnpidl
		return pidl  // return 0 if invalid pidl

	// otherwise return path
	if 0 is pidl
		return ""     // return "" for an invalid pidl

	return SHGetPathFromIDList(pidl)
	}
