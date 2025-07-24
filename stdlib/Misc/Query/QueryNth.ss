// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// ideally we'd just use Nth but queries don't behave like iterables
function (n, query)
	{
	QueryApply(query)
		{|x|
		if n-- is 0
			return x
		}
	return false
	}