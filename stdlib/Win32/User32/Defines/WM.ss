// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Window Messages
#(
NULL:                         0x0000
CREATE:                       0x0001
DESTROY:                      0x0002
MOVE:                         0x0003
SIZE:                         0x0005

ACTIVATE:                     0x0006

SETFOCUS:                     0x0007
KILLFOCUS:                    0x0008
ENABLE:                       0x000A
SETREDRAW:                    0x000B
SETTEXT:                      0x000C
GETTEXT:                      0x000D
GETTEXTLENGTH:                0x000E
PAINT:                        0x000F
CLOSE:                        0x0010
QUERYENDSESSION:              0x0011
QUIT:                         0x0012
QUERYOPEN:                    0x0013
ERASEBKGND:                   0x0014
SYSCOLORCHANGE:               0x0015
ENDSESSION:                   0x0016
SHOWWINDOW:                   0x0018
CTLCOLOR:					  0x0019
WININICHANGE:                 0x001A
SETTINGCHANGE:                0x001A

DEVMODECHANGE:                0x001B
ACTIVATEAPP:                  0x001C
FONTCHANGE:                   0x001D
TIMECHANGE:                   0x001E
CANCELMODE:                   0x001F
SETCURSOR:                    0x0020
MOUSEACTIVATE:                0x0021
CHILDACTIVATE:                0x0022
QUEUESYNC:                    0x0023

GETMINMAXINFO:                0x0024

PAINTICON:                    0x0026
ICONERASEBKGND:               0x0027
NEXTDLGCTL:                   0x0028
SPOOLERSTATUS:                0x002A
DRAWITEM:                     0x002B
MEASUREITEM:                  0x002C
DELETEITEM:                   0x002D
VKEYTOITEM:                   0x002E
CHARTOITEM:                   0x002F
SETFONT:                      0x0030
GETFONT:                      0x0031
SETHOTKEY:                    0x0032
GETHOTKEY:                    0x0033
QUERYDRAGICON:                0x0037
COMPAREITEM:                  0x0039

COMPACTING:                   0x0041
WINDOWPOSCHANGING:            0x0046
WINDOWPOSCHANGED:             0x0047

POWER:                        0x0048

COPYDATA:                     0x004A
CANCELJOURNAL:                0x004B

NOTIFY:                       0x004E
INPUTLANGCHANGEREQUEST:       0x0050
INPUTLANGCHANGE:              0x0051
TCARD:                        0x0052
HELP:                         0x0053
USERCHANGED:                  0x0054
NOTIFYFORMAT:                 0x0055

CONTEXTMENU:                  0x007B
STYLECHANGING:                0x007C
STYLECHANGED:                 0x007D
DISPLAYCHANGE:                0x007E
GETICON:                      0x007F
SETICON:                      0x0080

NCCREATE:                     0x0081
NCDESTROY:                    0x0082
NCCALCSIZE:                   0x0083
NCHITTEST:                    0x0084
NCPAINT:                      0x0085
NCACTIVATE:                   0x0086
GETDLGCODE:                   0x0087

NCMOUSEMOVE:                  0x00A0
NCLBUTTONDOWN:                0x00A1
NCLBUTTONUP:                  0x00A2
NCLBUTTONDBLCLK:              0x00A3
NCRBUTTONDOWN:                0x00A4
NCRBUTTONUP:                  0x00A5
NCRBUTTONDBLCLK:              0x00A6
NCMBUTTONDOWN:                0x00A7
NCMBUTTONUP:                  0x00A8
NCMBUTTONDBLCLK:              0x00A9

KEYFIRST:                     0x0100
KEYDOWN:                      0x0100
KEYUP:                        0x0101
CHAR:                         0x0102
DEADCHAR:                     0x0103
SYSKEYDOWN:                   0x0104
SYSKEYUP:                     0x0105
SYSCHAR:                      0x0106
SYSDEADCHAR:                  0x0107
KEYLAST:                      0x0108

IME_STARTCOMPOSITION:         0x010D
IME_ENDCOMPOSITION:           0x010E
IME_COMPOSITION:              0x010F
IME_KEYLAST:                  0x010F

INITDIALOG:                   0x0110
COMMAND:                      0x0111
SYSCOMMAND:                   0x0112
TIMER:                        0x0113
HSCROLL:                      0x0114
VSCROLL:                      0x0115
INITMENU:                     0x0116
INITMENUPOPUP:                0x0117
MENUSELECT:                   0x011F
MENUCHAR:                     0x0120
ENTERIDLE:                    0x0121

CTLCOLORMSGBOX:               0x0132
CTLCOLOREDIT:                 0x0133
CTLCOLORLISTBOX:              0x0134
CTLCOLORBTN:                  0x0135
CTLCOLORDLG:                  0x0136
CTLCOLORSCROLLBAR:            0x0137
CTLCOLORSTATIC:               0x0138

MOUSEFIRST:                   0x0200
MOUSEMOVE:                    0x0200
LBUTTONDOWN:                  0x0201
LBUTTONUP:                    0x0202
LBUTTONDBLCLK:                0x0203
RBUTTONDOWN:                  0x0204
RBUTTONUP:                    0x0205
RBUTTONDBLCLK:                0x0206
MBUTTONDOWN:                  0x0207
MBUTTONUP:                    0x0208
MBUTTONDBLCLK:                0x0209
MOUSEWHEEL:					  0x020a
MOUSELAST:					  0x020a

PARENTNOTIFY:                 0x0210
ENTERMENULOOP:                0x0211
EXITMENULOOP:                 0x0212

NEXTMENU:                     0x0213

SIZING:                       0x0214
CAPTURECHANGED:               0x0215
MOVING:                       0x0216
POWERBROADCAST:               0x0218
DEVICECHANGE:                 0x0219

IME_SETCONTEXT:               0x0281
IME_NOTIFY:                   0x0282
IME_CONTROL:                  0x0283
IME_COMPOSITIONFULL:          0x0284
IME_SELECT:                   0x0285
IME_CHAR:                     0x0286
IME_KEYDOWN:                  0x0290
IME_KEYUP:                    0x0291

MDICREATE:                    0x0220
MDIDESTROY:                   0x0221
MDIACTIVATE:                  0x0222
MDIRESTORE:                   0x0223
MDINEXT:                      0x0224
MDIMAXIMIZE:                  0x0225
MDITILE:                      0x0226
MDICASCADE:                   0x0227
MDIICONARRANGE:               0x0228
MDIGETACTIVE:                 0x0229

MDISETMENU:                   0x0230
ENTERSIZEMOVE:                0x0231
EXITSIZEMOVE:                 0x0232
DROPFILES:                    0x0233
MDIREFRESHMENU:               0x0234

MOUSEHOVER:                   0x02A1
MOUSELEAVE:                   0x02A3

CUT:                          0x0300
COPY:                         0x0301
PASTE:                        0x0302
CLEAR:                        0x0303
UNDO:                         0x0304
RENDERFORMAT:                 0x0305
RENDERALLFORMATS:             0x0306
DESTROYCLIPBOARD:             0x0307
DRAWCLIPBOARD:                0x0308
PAINTCLIPBOARD:               0x0309
VSCROLLCLIPBOARD:             0x030A
SIZECLIPBOARD:                0x030B
ASKCBFORMATNAME:              0x030C
CHANGECBCHAIN:                0x030D
HSCROLLCLIPBOARD:             0x030E
QUERYNEWPALETTE:              0x030F
PALETTEISCHANGING:            0x0310
PALETTECHANGED:               0x0311
HOTKEY:                       0x0312

PRINT:                        0x0317
PRINTCLIENT:                  0x0318

THEMECHANGED:				  0x031a

HANDHELDFIRST:                0x0358
HANDHELDLAST:                 0x035F

AFXFIRST:                     0x0360
AFXLAST:                      0x037F

PENWINFIRST:                  0x0380
PENWINLAST:                   0x038F

USER:						  0x0400
APP:						  0x8000
APP_SETFOCUS:				  0x8001
)