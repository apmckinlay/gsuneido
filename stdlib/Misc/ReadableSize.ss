// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(n)
		{
		.FromInt(n)
		}
	FromInt(n)
		{
		// maxDisplayWidth is the width of the number portion of the result.
		// The decimal part will be increased when the whole number is smaller.
		maxDisplayWidth = 4
		amountPerUnit = 1024
		for unit in #('', kb, mb, gb, tb)
			{
			if n < amountPerUnit
				return Number(String(n)[..maxDisplayWidth]) $ Opt(' ', unit)
			n /= amountPerUnit
			}
		return n.Round(2) $ ' pb'
		}
	ToInt(s)
		{
		if s.Number?()
			return Number(s)
		s = s.Lower()
		Assert(s =~ `^[\d.]+ ?(kb|mb|gb|tb)$`)
		f = #(kb: 1024, mb: 1048576, gb: 1073741824, tb: 1099511627776)[s[-2 ..]]
		return (f * Number(s[.. -2])).Round(0)
		}
	}