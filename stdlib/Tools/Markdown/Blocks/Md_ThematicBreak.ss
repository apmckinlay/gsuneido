// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	Closed?: true
	New()
		{
		}

	Match(line, container = false, checkingContinuationText? = false)
		{
		if container is false and not checkingContinuationText? and
			Md_Paragraph.IsSetextHeadingUnderline?(line)
			return false
		if false is n = .IgnoreLeadingSpaces(line)
			return false
		sub = line[n..].Tr(' \t')
		if sub =~ `^(\-\-\-+|___+|\*\*\*+)$`
			return new this()
		return false
		}
	}