// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// GlobalAlloc flags
#(
FIXED:          0x0000
MOVEABLE:       0x0002
NOCOMPACT:      0x0010
NODISCARD:      0x0020
ZEROINIT:       0x0040
MODIFY:         0x0080
DISCARDABLE:    0x0100
NOT_BANKED:     0x1000
SHARE:          0x2000
DDESHARE:       0x2000
NOTIFY:         0x4000
LOWER:          0x1000 // GMEM_NOT_BANKED
VALID_FLAGS:    0x7F72
INVALID_HANDLE: 0x8000
)