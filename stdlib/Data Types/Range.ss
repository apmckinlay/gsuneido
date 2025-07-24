// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// NOTE: this class is intentionally immutable
// Do not add any methods that change an instance (e.g. set low or high)
// after it has been constructed.
// It is fine to add methods that return a new instance like Plus and Minus.
class
	{
	UseDeepEquals: true
	New(.low = false, .high = false)
		{
		if low > high
			{
			SuneidoLog("ERROR (CAUGHT): Range, low > high", params: Record(:low, :high),
				calls:, caughtMsg: 'need to look at why invalid values passed to Range')
			if Sys.SuneidoJs?()
				SuRenderBackend().DumpStatus('Range low greater than high')
			.high = low // make it an empty range
			}
		}
	GetLow()
		{
		return .low
		}
	GetHigh()
		{
		return .high
		}
	Contains?(value)
		{
		return .low <= value and value <= .high
		}
	Includes?(range)
		{
		return .low <= range.GetLow() and range.GetHigh() <= .high
		}
	Overlaps?(range)
		{
		return .low <= range.GetHigh() and .high >= range.GetLow()
		}
	Plus(range)
		{
		// returns a new range that is the smallest that includes this and range
		return Range(Min(.low, range.GetLow()), Max(.high, range.GetHigh()))
		}
	Minus(range)
		{
		if ((.low < range.GetLow() and range.GetHigh() < .high) or not .Overlaps?(range))
			return this
		else if range.Includes?(this)
			return Range()
		else if .high > range.GetHigh()
			return Range(range.GetHigh(), .high)
		else
			return Range(.low, range.GetLow())
		}
	Separate(range)
		{
		return Object(Range(.low, range.GetLow()), Range(range.GetHigh(), .high))
		}
	Empty?()
		{
		return .low is false and .high is false
		}
	ToString()
		{
		"Range(" $ Display(.low) $ ", " $ Display(.high) $ ")"
		}
	}