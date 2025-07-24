// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.

// TODO allow renaming by selecting parameter (MethodRange doesn't work)

Refactor_Rename
	{
	Name: 'Rename Variable'
	Desc:
'Replace all occurrences of one local variable name with another.
Does not replace an identical member name e.g. .name
Warns if To is already used.'

	Init(data, omitPreview? = false)
		{
		if false is super.Init(data)
			return false
		if .function?(data.text)
			{
			s = data.text
			.varpos = data.select.cpMin
			.MethodRange = [from: 0, to: s.Size()]
			}
		else
			{
			range = .MethodRange = ClassHelp.MethodRange(data.text, data.select.cpMin)
			while data.text[range.from - 1] is '\t'
				--range.from
			s = data.text[range.from .. range.to]
			indent = s.Extract("^\t+")
			s = s.Replace("^" $ indent, "")
			.varpos = .FindMatch(s, data.from)
			}
		.MethodText = s
		.data = data
		if omitPreview? is false
			.InitializePreview(data, .MethodText)

		return true
		}

	function?(text)
		{
		return ScannerWithContext(text).Next() is 'function'
		}

	Change(member)
		{
		if member is 'from'
			.varpos = (.data.from =~ .NamePattern)
				? .FindMatch(.MethodText, .data.from) : -1

		if .data.to =~ .NamePattern and .data.from =~ .NamePattern
			.Timer.Reset()
		}

	NamePattern: '^[a-z][a-zA-Z0-9_]*[!?]?$'

	Matches?(name, prev2, prev, token, next)
		{
		//FIXME treat .name in parameters as use of name
		return token is name and
			prev isnt '' and prev2 isnt '' and // skip method name
			prev isnt '.' and next isnt ':' and
			prev isnt '#'
		}

	Warnings(data)
		{
		// limit to method range
		src = data.text[.MethodRange.from .. .MethodRange.to]
		if .ToExists?(src, data.to)
			return '"' $ data.to $ '" is already used'

		if .varpos is -1
			return 'Variable not found and will not be renamed'

		return ""
		}

	Process(data)
		{
		i = .MethodRange.from
		n = .MethodRange.to - i
		text = .Rename(data.text[i :: n], data.from, data.to, data.inComments)
		data.text = data.text.ReplaceSubstr(i, n, text)
		return true
		}
	}
