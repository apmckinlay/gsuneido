// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Win32/User32/Defines/DFCS
// State flags for DrawFrameControl
#(
CAPTIONCLOSE:		0x0000	// Close caption button [X]
CAPTIONMIN:			0x0001	// Minimize caption button
CAPTIONMAX:			0x0002	// Maximize caption button
CAPTIONRESTORE:		0x0003	// Restore caption button
CAPTIONHELP:		0x0004	// Context help caption button [?]

MENUARROW:			0x0000	// Submenu arrow
MENUCHECK:			0x0001	// Menu check mark
MENUBULLET:			0x0002	// Menu bullet
MENUARROWRIGHT:	0x0004	// Submenu arrow pointing left (for right-to-left languages
SCROLLUP:			0x0000	// Up arrow of scroll bar
SCROLLDOWN:		0x0001	// Down arrow of scroll bar
SCROLLLEFT:			0x0002	// Left scroll bar arrow
SCROLLRIGHT:		0x0003	// Right scroll bar arrow
SCROLLCOMBOBOX:	0x0005	// Combo box scroll bar
SCROLLSIZEGRIP:		0x0008	// Size grip in bottom-right corner of window
SCROLLSIZEGRIPRIGHT:	0x0010	// Size grip in bottom-left corner of window (right-to-left languages)

BUTTONCHECK:		0x0000	// Check box
BUTTONRADIOIMAGE:	0x0001	// Image for radio button (non-square needs image)
BUTTONRADIOMASK:	0x0002	// Mask for radio button (non-square needs mask)
BUTTONRADIO:		0x0004	// Radio button
BUTTON3STATE:		0x0008	// Three-state button
BUTTONPUSH:			0x0010	// Push button

INACTIVE:			0x0100	// Button is inactive (grayed out)
PUSHED:				0x0200	// Button is pushed
CHECKED:			0x0400	// Button is checked

TRANSPARENT:		0x0800	// if WinVer > 0x0500 ; Background remains untouched
HOT:				0x1000	// if WinVer > 0x0500 ; Button is hot-tracked

ADJUSTRECT:			0x2000	// Bounding rectangle of push button is adjusted to exclude surrounding edge
FLAT:				0x4000	// Button has flat border
MONO:				0x8000	// Button has monochrome border
)
