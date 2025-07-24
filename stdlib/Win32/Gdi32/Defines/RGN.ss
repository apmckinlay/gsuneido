// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Win32/Gdi32/Defines/RGN
// Region combine modes for CombineRgn
#(
ERROR:	0	// RGN_ERROR == ERROR (REGIONTYPE.ERROR)
AND:	1	// Creates the intersection of the two combined regions.
OR:		2	// Creates the union of two combined regions.
XOR:	3	// Creates the union of two combined regions except for any overlapping areas.
DIFF:	4	// Combines the parts of hrgnSrc1 that are not part of hrgnSrc2.
COPY:	5	// Creates a copy of the region identified by hrgnSrc1.
MIN:	1	// #define RGN_MIN	RGN_AND
MAX:	5	// #define RGN_MAX	RGN_COPY
)
