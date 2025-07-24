// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (@args) /*usage: (query, limit = false) */
	{
	maxServerEval = 100
	NameArgs(args, #(query, limit), #(false))
	if args.limit isnt false and args.limit <= maxServerEval and Client?()
		return ServerEval('QueryAll', args.query, args.limit)
	Transaction(read:)
		{|t|
		t.QueryAll(@args)
		}
	}
