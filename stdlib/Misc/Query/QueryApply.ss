// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
function (@args) //query, block, update = false, dir = 'Next')
	{
	update = args.Extract(#update, false)
	// t.QueryApply (in Transactions) does QueryAddKeySort
	// but we do it here as well to make it mandatory
	if update is true
		args[0] = QueryEnsureKeySort(args[0])
	Transaction(:update)
		{|t|
		t.QueryApply(@args)
		}
	}