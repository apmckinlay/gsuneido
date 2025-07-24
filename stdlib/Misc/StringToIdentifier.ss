// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (str)
	{
	return String(str).Map(function (c)
		{ (c.AlphaNum?() or c is '_') ? c : c.Asc().Hex() })
	}