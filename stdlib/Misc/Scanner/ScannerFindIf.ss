// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (scan_string, block)
	{
	scanner = Scanner(scan_string)
	for token in scanner
		if block(token)
			return scanner.Position() - token.Size()
	return scan_string.Size()
	}
