// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Num(s)
		{
		return s is "inf" ? 1/0 : s is "-inf" ? -1/0 : Number(s)
		}
	Nums(args, block)
		{
		args.Delete(#str?)
		return block(@args.Map(.Num))
		}
	Binary(op, x, y, z)
		{
		result = op(x, y)
		if result.RoundToPrecision(15) is z.RoundToPrecision(15)
			return true
		Print('\t', x, y, "=>", result, "should be", z)
		return false
		}
	}