// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// tab delimited import
Import
	{
	Header()
		{
		.Fields = .Getline().Split('\t')
		}
	Import1(line)
		{
		rec = Record()
		values = line.Split('\t')

		n = Min(values.Size(), .Fields.Size())
		for (i = 0; i < n; ++i)
			{
			s = values[i]
			if s.Prefix?('"') and s.Suffix?('"')
				s = s[1 .. -1]
			rec[.Fields[i]] = s
			}
		return rec
		}
	}
