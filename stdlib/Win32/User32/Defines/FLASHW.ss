// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Flags for the FLASHWINFO uFlags field (used for FlashWindowEx)
#(
ALL:       0x00000003 //Flash both the window caption and taskbar button. This is equivalent to setting the FLASHW_CAPTION | FLASHW_TRAY flags.
CAPTION:   0x00000001 // Flash the window caption.
STOP:      0          // Stop flashing. The system restores the window to its original state.
TIMER:     0x00000004 // Flash continuously, until the FLASHW_STOP flag is set.
TIMERNOFG: 0x0000000C // Flash continuously until the window comes to the foreground.
TRAY:      0x00000002 // Flash the taskbar button.
)