// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
// useful for debugging
// e.g. BitNames(0x40000002, 'WS', 'TTS') => "WS.CHILD | TTS.NOPREFIX"
// Note: if there are multiple names for the same single bit,
// it's unpredictable which you'll get. However, if NAME_COMBO is a name for
// NAME1|NAME2, then always returns NAME_COMBO.
// If a name should never be returned, add it to the ignore list.
class
	{
	separator: ' | '
	CallClass(@args)
		{
		bits = args.PopFirst()
		if bits is 0
			return "0"
		names = ""
		for defsname in args
			{
			defs = .transform(Global(defsname), defsname)
			for tuple in defs
				{
				mask = tuple[0]
				if mask is (bits & mask)
					{
					names $= .separator $ defsname $ '.' $ tuple[1]
					bits &= ~mask
					}
				}
			}
		return names[.separator.Size()..]
		}
	bitcount(x)
		{
		count = 0
		while x isnt 0
			{
			if 0x1 is (x & 0x1)
				++count
			x >>= 1
			}
		return count
		}
	transform(defs, defsname)
		{
		// Convert a map of form (NAME1: bits, NAME2: bits) into a list of
		// form [ [bits, NAME], [bits, NAME] ] in which ignorable names are
		// removed and the list is sorted so that names combining more bits sort
		// before names combining fewer bits. This allows us to prefer combo
		// names (i.e. if NAME_COMBO = NAME1|NAME2, prefer NAME_COMBO).
		list = Object()
		ignore = .ignore.Member?(defsname) ? .ignore[defsname]: #()
		for m in defs.Members()
			{
			if ignore.Member?(m)
				continue
			mask = defs[m]
			if mask is 0
				continue
			list.Add(Object(mask, m))
			}
		list.Sort!(
			{|x,y|
			.bitcount(x[0]) > .bitcount(y[0])
			}
		)
		}
	ignore: (WM: (MOUSEFIRST:, MOUSELAST:, KEYFIRST:, KEYLAST:))
	}