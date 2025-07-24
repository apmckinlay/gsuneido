// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
function (s)
	{
	if not Suneido.Member?('csDevPrint')
		Suneido.csDevPrint = Object()
	Suneido.csDevPrint.Add(s)
	}
