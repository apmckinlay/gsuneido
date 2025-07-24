// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// for SetWindowPos
#(
NOSIZE:          0x0001
NOMOVE:          0x0002
NOZORDER:        0x0004
NOREDRAW:        0x0008
NOACTIVATE:      0x0010
FRAMECHANGED:    0x0020  /* The frame changed: send WM_NCCALCSIZE */
SHOWWINDOW:      0x0040
HIDEWINDOW:      0x0080
NOCOPYBITS:      0x0100
NOOWNERZORDER:   0x0200  /* Dont do owner Z ordering */
NOSENDCHANGING:  0x0400  /* Dont send WM_WINDOWPOSCHANGING */

DRAWFRAME:       0x0020 // SWP_FRAMECHANGED
NOREPOSITION:    0x0200 // SWP_NOOWNERZORDER

DEFERERASE:      0x2000
ASYNCWINDOWPOS:  0x4000
)