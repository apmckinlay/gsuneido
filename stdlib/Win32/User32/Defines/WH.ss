// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// constants for SetWindowsHook and SetWindowsHookEx
#(
MIN:				-1
MSGFILTER:			-1
JOURNALRECORD:		0
JOURNALPLAYBACK:	1
KEYBOARD:			2
GETMESSAGE:			3
CALLWNDPROC:		4
CBT:				5
SYSMSGFILTER:		6
MOUSE:				7
HARDWARE:			8
DEBUG:				9
SHELL:				10
FOREGROUNDIDLE:		11

//#if(WINVER >= 0x0400)
CALLWNDPROCRET:		12
//#endif /* WINVER >= 0x0400 */

KEYBOARD_LL:		13
MOUSE_LL:			14

MAX:				14
//#endif
//MINHOOK:			WH_MIN
//MAXHOOK:			WH_MAX
)
