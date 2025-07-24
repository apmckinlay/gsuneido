// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(folder, block, skipFn? = false)
		{
		if not folder.Suffix?('/')
			folder $= '/'
		for item in .dir(folder)
			{
			cur = folder $ item.name
			if skipFn? isnt false and skipFn?(cur)
				continue

			if item.name.Suffix?('/') /*folder*/
				this(cur, block, :skipFn?)
			else
				{
				item.name = cur
				block(item)
				}
			}
		}

	dir(path)
		{
		return Dir(path $ '*.*', details:)
		}
	}