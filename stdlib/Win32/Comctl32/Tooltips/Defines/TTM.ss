// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Win32/ComCtl32/Tooltip/Defines/TTM
// Tooltip message constants
// Note that for all messages in which the ASCII and Unicode versions differed,
// the ASCII versions were defined
#(
ACTIVATE:			0x0401 // (WM_USER + 1)
SETDELAYTIME:		0x0403 // (WM_USER + 3)
ADDTOOL:			0x0404 // (WM_USER + 4)
DELTOOL:			0x0405 // (WM_USER + 5)
NEWTOOLRECT:		0x0406 // (WM_USER + 6)
RELAYEVENT:			0x0407 // (WM_USER + 7)
GETTOOLINFO:		0x0408 // (WM_USER + 8)
SETTOOLINFO:		0x0409 // (WM_USER + 9)
HITTEST:			0x040A // (WM_USER + 10)
GETTEXT:			0x040B // (WM_USER + 11)
UPDATETIPTEXT:		0x040C // (WM_USER + 12)
GETTOOLCOUNT:		0x040D // (WM_USER + 13)
ENUMTOOLS:			0x040E // (WM_USER + 14)
GETCURRENTTOOL:		0x040F // (WM_USER + 15)
WINDOWFROMPOINT:	0x0410 // (WM_USER + 16)
// wParam = TRUE/FALSE start end  lparam = LPTOOLINFO
TRACKACTIVATE:		0x0411 // (WM_USER + 17)
TRACKPOSITION:		0x0412 // (WM_USER + 18)  // lParam = dwPos
SETTIPBKCOLOR:		0x0413 // (WM_USER + 19)
SETTIPTEXTCOLOR:	0x0414 // (WM_USER + 20)
GETDELAYTIME:		0x0415 // (WM_USER + 21)
GETTIPBKCOLOR:		0x0416 // (WM_USER + 22)
GETTIPTEXTCOLOR:	0x0417 // (WM_USER + 23)
SETMAXTIPWIDTH:		0x0418 // (WM_USER + 24)
GETMAXTIPWIDTH:		0x0419 // (WM_USER + 25)
SETMARGIN:			0x041A // (WM_USER + 26)  // lParam = lprc
GETMARGIN:			0x041B // (WM_USER + 27)  // lParam = lprc
POP:				0x041C // (WM_USER + 28)
UPDATE:				0x041D // (WM_USER + 29)
ADJUSTRECT:			0x041f // (WM_USER + 31)
POPUP:				0x0422 // (WM_USER + 34)
)
