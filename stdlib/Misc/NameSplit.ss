// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (name, split_on = false)
	{
	if split_on isnt false
		names = name.Split(split_on).Reverse!()
	else
		names = name.Has?(',')
			? name.Split(',').Reverse!()
			: name.Split(' ')

	return Object(
		first: names[.. names.Size() > 1 ? -1 : 1].Join(' ').Trim(),
		last: names.Size() > 1 ? names.Last() : "")
	}