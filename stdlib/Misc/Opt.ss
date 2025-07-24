// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	return args.Any?({|s| s is "" }) ? "" : args.Join("")
	}