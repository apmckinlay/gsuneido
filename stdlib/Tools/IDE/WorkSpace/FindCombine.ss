// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
// Specifically for use by FindInLibraries to combine lists of line numbers
// Needs to be global to use from query
// WARNING: Find References won't see the usage in FindInLibraries
class
	{
	CallClass(@args)
		{
		if .stop?(.extractUid(args))
			return #(0)
		results = Object()
		args.Values(list:).Each()
			{
			res = .evaluate(it)
			if res.Empty?()
				return #()
			if String?(res[0]) // error
				return Object(res[0])
			results.Append(res)
			}
		return results.Sort!().Unique!()
		}
	evaluate(args)
		{
		return Global(args[0])(@+1args)
		}
	stop?(findUid)
		{
		return ServerSuneido.GetAt(#workSpaceFindStop, findUid, false)
		}
	extractUid(args)
		{
		if false is idx = args.FindIf(Date?)
			return false
		return args.Extract(idx)
		}
	}
