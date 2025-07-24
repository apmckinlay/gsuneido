// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Control: (CheckBox)
	Format: (Boolean width: 7)
	Encode(val)
		{
		vals = #('true':, yes:, y:,
			'false': false, no: false, n: false)
		return String?(val) ? vals.GetDefault(val.Lower(), val) : val
		}
	}

