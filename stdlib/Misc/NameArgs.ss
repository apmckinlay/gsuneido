// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// Note: modifies the args object passed in, return value is for convenience
function (args, names, defs = #())
	{
	nlist = args.Size(list:)
	if nlist > names.Size()
		throw "too many arguments"
	idefs = names.Size() - defs.Size()
	for (i = names.Size() - 1; i >= 0; --i)
		{
		if not args.Member?(names[i])
			if i < nlist
				args[names[i]] = args[i]
			else if i >= idefs
				args[names[i]] = defs[i - idefs]
			else
				throw "missing argument: " $ names[i]
		args.Erase(i)
		}
	return args
	}