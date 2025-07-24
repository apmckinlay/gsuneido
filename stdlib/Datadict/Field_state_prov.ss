// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Field_string
	{
	Heading: 'State/\nProv'
	Prompt: 'State/Prov'
	Control: (StateProv)
	Format: (Text width: 4 justify: "center")
	Encode(val)
		{
		val = super.Encode(val)
		return val.Upper()
		}
	}
