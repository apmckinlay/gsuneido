// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// CSV Import
Import
	{
	Header()
		{
		.Fields = .Getline().SplitCSV()
		}
	Import1( line )
		{
		return line.SplitCSV(.Fields, string_vals:)
		}
	}
