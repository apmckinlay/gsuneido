// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
// NOTE: Do not use this in application code.
// It is intended for use by stdlib e.g. for QueryColumns
function (query)
	{
	return query $ "
		/* CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE */
		/* CHECKQUERY SUPPRESS: UNION NOT DISJOINT */
		/* CHECKQUERY SUPPRESS: JOIN MANY TO MANY */"
	}