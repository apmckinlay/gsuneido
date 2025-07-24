// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Check(@args)
		{
		fn = args[0]
		expected = args.Last()
		Assert(fn(@args[1..-1].Map(.toOp)) is: .toOp(expected))
		}
	CheckError(@args)
		{
		fn = args[0]
		expected = args.Last()
		Assert({ fn(@args[1..-1].Map(.toOp)) } throws: expected)
		}
	toOp(rec)
		{
		if rec is false
			return false
		return rec.ListToNamed(#type, #value)
		}
	}