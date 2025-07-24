// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	File(args[0], 'w')
		{ |f|
		for s in args[1..]
			{
			if s.Size() > 15.Mb() /*=7mb limit base64 encoded twice (4/3 x 4/3)*/
				ProgrammerError("PutFile bigger than 15mb (" $ ReadableSize(s.Size()) $
					") to " $ args[0])
			f.Write(s)
			}
		}
	}
