// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (s)
	{
	return s.
		Replace('&lt;', '<').
		Replace('&gt;', '>').
		Replace('&quot;', '"').
		Replace('&apos;', "'").
		Replace('&amp;', '\&')
	}