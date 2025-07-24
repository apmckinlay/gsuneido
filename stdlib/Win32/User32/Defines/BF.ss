// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Border flags
#(
LEFT:                    0x0001
TOP:                     0x0002
RIGHT:                   0x0004
BOTTOM:                  0x0008

TOPLEFT:                 0x0003 // BF.TOP | BF.LEFT
TOPRIGHT:                0x0006 // BF.TOP | BF.RIGHT
BOTTOMLEFT:              0x0009 // BF.BOTTOM | BF.LEFT
BOTTOMRIGHT:             0x000c // BF.BOTTOM | BF.RIGHT
RECT:                    0x000f // BF.LEFT | BF.TOP | BF.RIGHT | BF.BOTTOM

DIAGONAL:                0x0010
// For diagonal lines, the BF_RECT flags specify the end point of the vector
// bounded by the rectangle parameter
DIAGONAL_ENDTOPLEFT:     0x0013 // BF.DIAGONAL | BF.TOP | BF.LEFT
DIAGONAL_ENDTOPRIGHT:    0x0016 // BF.DIAGONAL | BF.TOP | BF.RIGHT
DIAGONAL_ENDBOTTOMLEFT:  0x0019 // BF.DIAGONAL | BF.BOTTOM | BF.LEFT
DIAGONAL_ENDBOTTOMRIGHT: 0x001c // BF.DIAGONAL | BF.BOTTOM | BF.RIGHT

MIDDLE:                  0x0800 // Fill in the middle
SOFT:                    0x1000 // For softer buttons
ADJUST:                  0x2000 // Calculate the space left over
FLAT:                    0x4000 // For flat rather than 3D borders
MONO:                    0x8000 // For monochrome borders
)