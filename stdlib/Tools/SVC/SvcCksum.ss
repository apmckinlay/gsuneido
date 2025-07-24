// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (rec)
	{
	return Adler32().
		Update(rec.name).
		Update(rec.path).
		Update(String(rec.order)).
		Update(rec.lib_current_text.Trim()).
		Update(String(rec.lib_committed)).
		Value() & 0xffffffff /*= 32 bits */
			// need mask to be compatible with BuiltDate > 2025-05-09
	}