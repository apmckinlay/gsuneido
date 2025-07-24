// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// NOTE: do NOT delete this until old formulas are handled
function(val, field)
	{
	if not Number?(val)
		return 0

	dd = Datadict(field)
	return dd.Control.Member?('mask')
		? Number(val.Format(dd.Control.mask))
		: val
	}