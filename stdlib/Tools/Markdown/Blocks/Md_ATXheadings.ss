// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	Closed?: true
	New(.Inline, .Level)
		{
		}

	Match(line)
		{
		if false is n = .IgnoreLeadingSpaces(line)
			return false

		level = .CountLeadingChar(line[n..], '#')
		if level < 1 or level > 6 /*=max*/
			return false

		remain = line[n + level..]
		if remain.Blank?() or remain[0] in (' ', '\t')
			{
			inline = remain
			if false isnt match = remain.Match('\s#+\s*$')
				inline = remain[..match[0][0]]
			return new this(inline.Trim(), level)
			}
		return false
		}
	}