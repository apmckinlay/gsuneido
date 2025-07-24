// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	perm = OptContribution('UsernotePermission?', function (unused) { return true })
	return perm(this.path)
	}