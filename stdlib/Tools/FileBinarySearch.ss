// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
class
	{
	// f should be an open file, the file must be already sorted based on cmp
	// cmp should be a function that is passed a line and returns +1,0,-1 for >,==,<
	// when we return the file position is set to the first line >= (i.e. lower_bound)
	// IF desired, you can supply a function to verify the order.
	//		- This function must take two values, to compare:
	//			first it runs with - lo, mid (lo <= mid)
	//			then it runs with  - mid, hi (mid <= hi)
	//		you can manipulate the lines as wish in the orderCmp function
	CallClass(f, cmp, orderCmp = false)
		{
		lo = curPoint = 0
		hi = FileEndPos(f)
		while lo < hi
			{
			mid = ((lo + hi) / 2).Ceiling()
			if orderCmp isnt false and
				false is .orderMaintained?(f, orderCmp, lo, hi, mid)
					return false
			mid = RewindLines(f, mid)
			line = f.Readline()
			switch cmp(line)
				{
			case 0, 1:
				curPoint = hi = mid			// Prior to the evaluated line
			case -1:
				curPoint = lo = f.Tell() 	// Past the current line
			case false:
				if false isnt tmpHi = .highestValueBeforeInvalidSect(f, cmp, mid)
					{
					hi = tmpHi
					continue
					}
				f.Seek(mid)
				if false isnt tmpLo = .lowestValueAfterInvalidSect(f, cmp)
					{
					lo = tmpLo
					continue
					}
				return false // No valid lines to continue from
				}
			}
		f.Seek(curPoint)
		return true
		}

	orderMaintained?(f, orderCmp, lo, hi, mid)
		{
		loLine = .getLine(f, lo)
		midLine = .getLine(f, mid)
		hiLine = .getLine(f, hi)
		if orderCmp(loLine, midLine) and orderCmp(midLine, hiLine)
			return true
		return false
		}

	highestValueBeforeInvalidSect(f, cmp, pos)
		{
		beforeInvalidSect = false
		while pos > 0
			{
			pos = RewindLines(f, pos, 2)
			line = f.Readline()
			if false isnt beforeInvalidSect = cmp(line)
				break
			}
		if Number?(beforeInvalidSect) and beforeInvalidSect >= 0
			return pos
		return false
		}

	lowestValueAfterInvalidSect(f, cmp)
		{
		afterInvalidSect = false
		while false isnt line = f.Readline()
			if false isnt afterInvalidSect = cmp(line)
				break
		if Number?(afterInvalidSect) and afterInvalidSect <= 0
			return RewindLines(f, f.Tell())
		return false
		}

	getLine(f, pos)
		{
		RewindLines(f, pos)
		return f.Readline()
		}
	}