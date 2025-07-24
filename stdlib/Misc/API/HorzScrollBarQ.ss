// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// from MSDN:
// "The system hides and disables a standard scroll bar
// when equal minimum and maximum values are specified.
// The system also hides and disables a standard scroll bar
// if you specify a page size that includes the entire
// scroll range of the scroll bar."
function (hwnd)
	{
	GetScrollInfo(hwnd, SB.HORZ,
		si = Object(cbSize: SCROLLINFO.Size(), fMask: SIF.ALL))
	hidden = si.nMin is si.nMax or si.nPage > (si.nMax - si.nMin)
	return not hidden
	}