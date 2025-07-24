// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Refactor_Rename
	{
	Name: 'Rename Member'
	Desc:
'Replace all occurrences of a member variable with another.
Does not replace an identical local variable name
Warns if To is already used.'
	NamePattern: '^[a-zA-Z][a-zA-Z0-9_]*[!?]?$'

	Init(data)
		{
		if false is super.Init(data)
			return false
		if not ClassHelp.Class?(data.text)
			{
			if not OkCancel("Rename Member is normally used on a class.\n\nContinue?",
				.Name, flags: MB.ICONQUESTION)
				return false
			.nest = 99
			}
		.data = data
		.MethodText = .data.text
		.InitializePreview(data, .MethodText)
		return true
		}
	Change(member)
		{
		if member is 'from'
			.MethodText = .data.text
		else
			.renamed = .Rename(.data.text, .data.from, .data.to, .data.inComments)

		if .data.to =~ .NamePattern and .data.from =~ .NamePattern
			.Timer.Reset()
		}
	Warnings(data)
		{
		if .MethodText is .renamed
			return 'There is no change. Please check your member and/or renamed member.'

		return ''
		}
	nest: 0
	Matches?(name, prev2, prev, token, next)
		{
		if token is '{' or token is '(' or token is '['
			++.nest
		else if token is '}' or token is ')' or token is ']'
			--.nest

		if .nest is 1 and token is name and
			(next is ':' or next is '(')
			return true

		if token isnt name or prev isnt '.' or next is ':'
			return false
		if prev2 is 'this'
			return true
		if prev2 is ']' or prev2 is ')'
			return false
		if prev2 =~ .NamePattern
			return false
		return true
		}
	}
