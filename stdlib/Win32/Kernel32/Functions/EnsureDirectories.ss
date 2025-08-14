// Copyright (C) 2010 Axon Development Corporation All rights reserved worldwide.
// creates directories on path
// path is assumed to end with file name unless it ends with '/'
// WARNING: does not check for success or failure
class
	{
	CallClass(path)
		{
		path = Paths.ToStd(path)
		segs = path.Split('/')
		if path.Prefix?('//') // unc path
			{
			segs = segs[2..]
			segs[0] = '//' $ segs[0]
			}
		if not path.Suffix?('/')
			segs.Delete(segs.Size() - 1)
		dir = segs[0]
		segs.Delete(0)
		for seg in segs
			.ensureDir(dir $= '/' $ seg)
		}
	ensureDir(dir) // overridden by tests
		{
		EnsureDir(dir)
		}
	}