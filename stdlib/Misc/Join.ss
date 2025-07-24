// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	sep = args.PopFirst()
	return args.Remove("").Join(sep)
	}