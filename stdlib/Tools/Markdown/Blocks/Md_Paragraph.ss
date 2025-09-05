// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_Base
	{
	HeadingLevel: false
	New(line)
		{
		.raws = [line]
		}

	Continue(line, checkingContinuationText? = false)
		{
		if .BlankLine?(line)
			return false
		if Md_ContainerBlock.MatchParagraphInteruptableBlockItem(
			line, :checkingContinuationText?) isnt false
			return false
		return line
		}

	IsContinuationText?(line)
		{
		return false isnt .Continue(line, checkingContinuationText?:)
		}

	Add(line, lazyContinuation? = false)
		{
		if not lazyContinuation? and .IsSetextHeadingUnderline?(line)
			{
			.HeadingLevel = line.Has?('=') ? 1 : 2
			.Close()
			}
		else
			.raws.Add(line)
		}

	IsSetextHeadingUnderline?(line)
		{
		return false isnt .IgnoreLeadingSpaces(line) and line.Trim() =~ '^(=+|-+)$'
		}

	Close()
		{
		.Inline = .raws.Join('\n')
		super.Close()
		}
	}