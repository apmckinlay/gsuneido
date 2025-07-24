// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	controlStart: "^[ \t]*" $
		"(if|else|for|function|while|do|switch|case|class|struct|try|catch)"
	CharAdded(c)
		{
		if c isnt '\n'
			return
		curLine = .LineFromPosition()
		if curLine <= 0
			return
		lineLength = .GetLine(curLine).Tr('\r\n').Size()
		prevLine = .GetLine(curLine - 1).Tr('\r\n')
		if prevLine is "" // Enter at start of line
			return
		indent = .getIndent(prevLine, lineLength)
		.Paste(indent)
		}
	getIndent(prevLine, lineLength)
		{
		indent = prevLine.Extract("^[ \t]*")
		if ((lineLength > 0 and prevLine.Trim() isnt "{") or // Enter in middle of line
			.opening?(prevLine) or
			prevLine =~ .controlStart)
			indent $= "\t" // indent more
		else if .closing?(prevLine)
			indent = indent.Replace("\t", "", 1) // indent less
		return indent
		}
	opening?(line)
		{
		return (line.Tr('^{}') is "{" and not line.BeforeFirst('{').White?()) or
			(line.Tr('^()') is "(" and not line.BeforeFirst('(').White?())
		}
	closing?(line)
		{
		return line.Tr('^{}') is "}" or line.Tr('^()') is ")"
		}
	}