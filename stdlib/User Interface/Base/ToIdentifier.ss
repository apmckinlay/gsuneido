// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (s)
	{
	return s.Tr("^a-zA-Z0-9_ ").Trim().Tr(' ', '_')
	}