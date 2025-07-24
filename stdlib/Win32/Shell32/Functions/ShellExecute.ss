// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// use this instead of ShellExecute
// due to bug under Windows 9x
// see Programming Industrial Strength Windows page 216
function (hwnd, lpVerb, lpFile, lpParameters = '', lpDirectory = '', nShowCmd = false,
	fMask = 0)
	{
	sei = Object(cbSize: SHELLEXECUTEINFO.Size(), :fMask,
		:hwnd, :lpVerb, :lpFile, :lpParameters, :lpDirectory,
		nShow: nShowCmd is false ? SW.SHOWNORMAL : nShowCmd)
	return ShellExecuteEx(sei)
	}