// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// PageSetupDlg defines
#(
DEFAULTMINMARGINS:             0x00000000 // default (printer's)
INWININIINTLMEASURE:           0x00000000 // 1st of 4 possible

MINMARGINS:                    0x00000001 // use caller's
MARGINS:                       0x00000002 // use caller's
INTHOUSANDTHSOFINCHES:         0x00000004 // 2nd of 4 possible
INHUNDREDTHSOFMILLIMETERS:     0x00000008 // 3rd of 4 possible
DISABLEMARGINS:                0x00000010
DISABLEPRINTER:                0x00000020
NOWARNING:                     0x00000080 // must be same as PD_*
DISABLEORIENTATION:            0x00000100
RETURNDEFAULT:                 0x00000400 // must be same as PD_*
DISABLEPAPER:                  0x00000200
SHOWHELP:                      0x00000800 // must be same as PD_*
ENABLEPAGESETUPHOOK:           0x00002000 // must be same as PD_*
ENABLEPAGESETUPTEMPLATE:       0x00008000 // must be same as PD_*
ENABLEPAGESETUPTEMPLATEHANDLE: 0x00020000 // must be same as PD_*
ENABLEPAGEPAINTHOOK:           0x00040000
DISABLEPAGEPAINTING:           0x00080000
NONETWORKBUTTON:               0x00200000 // must be same as PD_*
)