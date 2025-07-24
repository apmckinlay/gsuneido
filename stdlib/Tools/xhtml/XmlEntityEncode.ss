// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (s)
	{
	return String(s).
		Replace('&', '\&amp;').
		Replace('<', '\&lt;').
		Replace('>', '\&gt;').
		Replace('"', '\&quot;').
		Replace("'", "\&apos;")
	}