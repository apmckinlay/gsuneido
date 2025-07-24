// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(calls)
		{
		.calls = calls
		.i = 0
		}
	Default(@args)
		{
		if not .calls.Member?(.i)
			throw "No call expected\n" $
				"got: " $ .show(args)

		calli = .calls[.i++]
		call = calli.Member?('result') ? calli[0] : calli

		args.RemoveIf({|x| Type(x) is 'Block' })

		if call isnt args
			throw "Incorrect call:\n" $
				"expected: " $ .show(call) $ "\n" $
				"got: " $ .show(args)

		if calli.Member?('result')
			return calli.result
		else
			return // no return value
		}
	show(args)
		{
		return '.' $ args[0] $ Display(args[1..])[1 ..]
		}
	}