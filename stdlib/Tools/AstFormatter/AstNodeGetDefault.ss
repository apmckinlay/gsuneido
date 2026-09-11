// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function(node, member, dflt = "")
	{
	try
		return node[member]
	catch (unused, "member not found")
		return dflt
	}
