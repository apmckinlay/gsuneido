// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
function (query, formatted = false)
	{
	if String?(query)
		WithQuery(query)
			{|q|
			s = q.Strategy(:formatted)
			}
	else
		s = query.Strategy(:formatted)
	if not formatted
		s = s.BeforeFirst('[nrecs~ ').RightTrim()
	return s
	}
