// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (name)
	{
	return LibCanonicalName(name.
		Replace("^(Field_|Rule_|Trigger_)", "").
		Replace("(Control|Format)$", ""))
	}