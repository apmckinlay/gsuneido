// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	n = args.Size()
	hinst = ResourceModule()
	imagelist = ImageList_Create(
		ScaleWithDpiFactor(args[0]), ScaleWithDpiFactor(args[1]),
		ILC.MASK | ILC.HIGHQUALITYSCALE | ILC.COLORDDB, n - 2, 1)
	for (i = 2; i < n; ++i)
		ImageList_AddIcon(imagelist, LoadIcon(hinst, args[i]))
	return imagelist
	}