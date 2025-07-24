// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Custom draw return flags -- possible return values to NM_CUSTOMDRAW
#(
DODEFAULT:         0x00000000 // The control will draw itself. It will not send any additional NM_CUSTOMDRAW notification codes for this paint cycle. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_PREPAINT.
NEWFONT:           0x00000002 // The application specified a new font for the item; the control will use the new font. For more information about changing fonts, see Changing fonts and colors. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_ITEMPREPAINT.
SKIPDEFAULT:       0x00000004 // The application drew the item manually. The control will not draw the item. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_ITEMPREPAINT.
DOERASE:           0x00000008 // Windows Vista and later. The control will draw the background.
NOTIFYPOSTPAINT:   0x00000010 // The control will notify the parent after painting an item. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_PREPAINT.
NOTIFYITEMDRAW:    0x00000020 // The control will notify the parent of any item-related drawing operations. It will send NM_CUSTOMDRAW notification codes before and after drawing items. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_PREPAINT.
NOTIFYSUBITEMDRAW: 0x00000020 // Internet Explorer 4.0 and later. The control will notify the parent of any item-related drawing operations. It will send NM_CUSTOMDRAW notification codes before and after drawing items. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_PREPAINT. This flag is identical to CDRF_NOTIFYITEMDRAW and its use is context-dependent.
NOTIFYPOSTERASE:   0x00000040 // The control will notify the parent after erasing an item. This occurs when the dwDrawStage of the NMCUSTOMDRAW structure equals CDDS_PREPAINT.
SKIPPOSTPAINT:     0x00000100 // Windows Vista and later. The control will not draw the focus rectangle.
)