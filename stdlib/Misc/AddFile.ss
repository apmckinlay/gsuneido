// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	file = args[0]
	limit = args.Extract(#limit, 1e99) /*= bigger than any real size */
	File(file, 'a')
		{ |f|
		if limit < size = f.Size()
			{
			SuneidoLog.Once("ERROR: AddFile: " $ file $ " size " $ size $
				" > limit " $ limit)
			return
			}
		for s in args[1..]
			{
			if s.Size() > 10.Mb() /*= too big */
				ProgrammerError("AddFile bigger than 10mb to " $ args[0])
			f.Write(s)
			}
		}
	}
