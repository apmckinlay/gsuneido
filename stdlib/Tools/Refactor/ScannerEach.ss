// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (src, block)
	{
	scan = ScannerWithContext(src, wantWhitespace:)
	while scan isnt scan.Next()
		block(scan.Prev2(), scan.Prev(), scan.Token(), scan.Ahead())
	}