// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
function()
	{
	sigs = Object()
	Plugins().ForeachContribution(#TypeSignatures, #Signatures)
		{|c|
		if c.Member?(#sigs)
			{
			sigs.Append(c.sigs)
			continue
			}
		source = Global(c.source)
		prefix = c.source $ "_sig_"
		for m in source.Members().Sort!()
			{
			if not m.Prefix?(prefix)
				continue
			entry = Object(name: m.RemovePrefix(prefix), sig: source[m])
			if c.Member?(#receiver)
				entry.receiver = c.receiver
			else if c.Member?(#class)
				{
				entry.kind = #static
				entry.class = c.class
				}
			else
				entry.kind = #free
			sigs.Add(entry)
			}
		}
	return sigs
	}
