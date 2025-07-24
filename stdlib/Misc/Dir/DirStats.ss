// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(path)
		{
		return .fn(Paths.ToStd(path))
		}

	fn(path)
		{
		resultOb = Object()
		stats = Object(size: 0, fileN: 0,
			largest: false, mostRecent: false, leastRecent: false)
		.dir(path)
			{ |item|
			if item.name.Suffix?('/')
				{
				resultOb[item.name] = .fn(Paths.Combine(path, item.name))
				stats.size += resultOb[item.name][0].size
				}
			else
				{
				stats.size += item.size
				stats.fileN++
				if stats.largest is false or item.size > stats.largest.size
					stats.largest = item
				if stats.mostRecent is false or item.date > stats.mostRecent.date
					stats.mostRecent = item
				if stats.leastRecent is false or item.date < stats.leastRecent.date
					stats.leastRecent = item
				}
			}
		resultOb.Add(stats)
		return resultOb
		}

	dir(path, block)
		{
		Dir(Paths.Combine(path, '*.*'), details:, :block)
		}
	}
