// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (x)
	{
	if Object?(x)
		keys = x
	else if String?(x)
		keys = QueryKeys(x)
	else
		keys = x.Keys() // query or cursor
	fieldCount = function (s)
		{ return s.Count(',') + 1 }
	key = keys.MinWith()
		{ |k|
		fieldCount(k) -
			(k.Has?('_num') ? .5 : 0) + /*= prefer _num */
			(k.Has?('_name') ? .5 : 0)  /*= avoid _name */
		}
	if fieldCount(key) > 10 /*= nothing special about 10, just a reasonable threshold */
		ProgrammerError("key too large", [:key])
	return key
	}