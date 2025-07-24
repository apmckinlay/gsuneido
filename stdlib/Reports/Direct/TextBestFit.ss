// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(w, str, measure, report = false) // Called From Reports
		{
		if (str is '')
			return false
		if report isnt false and
			report.PlainText?()
			return str
		return .bestFit(w, str, measure)
		}
	bestFit(w, str, measure)
		{
		line = ''
		lo = 0
		hi = str.Find('\n') + 1
		best = 1
		while lo <= hi
			{
			mid = ((lo + hi) / 2).Int()
			line = str[.. mid]
			if measure(line.RightTrim()) <= w
				{
				if (mid > best)
					best = mid
				lo = mid + 1
				}
			else
				hi = mid - 1
			}
		line = str[.. best]
		if best < str.Size() and not line.Suffix?('\n')
			{
			// back up to word break (if there is one)
			if false isnt pos = .lastbreak(line)
				line = line[.. pos + 1]
			}
		Assert(line isnt '')
		return line
		}
	lastbreak(line)
		{
		while false isnt pos = line.FindLast1of(', ')
			{
			if (.validbreak?(line, pos))
				break
			line = line[.. pos]
			}
		return pos
		}
	validbreak?(line, pos)
		{
		return line[pos] isnt ',' or
			not line[pos - 1].Numeric?() or not line[pos + 1].Numeric?()
		}
	}