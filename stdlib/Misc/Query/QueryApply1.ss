// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (@args)
	{
	Assert(args.Extract(#update, "") is "",
		msg: "QueryApply1 does not take update: (it's always update)")
	RetryTransaction()
		{ |t|
		t.QueryApply1(@args)
		}
	}