// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Win32/ComCtl32/Tooltip/Defines/TTF
// Tooltip item flag constants
#(
IDISHWND:	0x0001
// Use this to center around trackpoint in trackmode -OR- to center around tool in normal mode.
CENTERTIP:	0x0002
RTLREADING:	0x0004
SUBCLASS:	0x0010
TRACK:		0x0020
// Use TTF_ABSOLUTE to place the tip exactly at the track coords when  in tracking mode.  TTF_ABSOLUTE
// can be used in conjunction with TTF_CENTERTIP to center the tip absolutely about the track point.
ABSOLUTE:	0x0080
TRANSPARENT:0x0100
DI_SETITEM:	0x8000       // valid only on the TTN_NEEDTEXT callback
)
