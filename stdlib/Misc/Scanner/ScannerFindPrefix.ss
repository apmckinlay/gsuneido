// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (scan_string, find_string)
	{
	scanner = Scanner(scan_string)
	for token in scanner
		if scan_string[scanner.Position() - token.Size() ..].Prefix?(find_string)
			return scanner.Position() - token.Size()
	return scan_string.Size()
	}
