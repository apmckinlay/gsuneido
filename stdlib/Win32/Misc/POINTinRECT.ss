// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (r, p)
	{
	return r.left <= p.x and p.x <= r.right and
		r.top <= p.y and p.y <= r.bottom
	}