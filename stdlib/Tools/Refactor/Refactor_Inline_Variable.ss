// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.

// TODO error if variable is parameter
// TODO warn if variable is assigned more than once (including e.g. $=, ++)
// TODO handle if variable definition extends over multiple lines
// TODO allow inlining a single occurrence (search for previous definition)
// MAYBE add option to enclose expression in parenthesis

Refactor_Rename_Variable
	{
	Name: "Inline Variable"
	Desc: "Replace all occurrences of a local variable with its definition"
	DiffPos: 6
	Controls: (Vert
		(Pair (Static Replace) (Field font: '@mono', readonly:, name: from))
		Skip
		(Pair (Static With) (Field font: '@mono', xstretch: 1, readonly:, name: to))
		Skip
		(Static Preview)
		(Skip 3)
		(Diff2 '', '', '', '', 'From', 'To')
		name: 'renameVert'
		)
	Init(data)
		{
		super.Init(data, true)
		if data.from !~ .NamePattern
			{
			.Warn("Place the cursor on the variable\n" $
				"at the point where it is initialized")
			return false
			}
		s = data.text
		i = s.FindLast('\n', data.select.cpMin) + 1
		j = s.Find('\n', data.select.cpMax)
		data.to = s[i .. j].AfterLast('=').Trim()
		.line = '^\t*(?q)' $ s[i .. j + 1].LeftTrim()
		.InitializePreview(data, .MethodText)
		return true
		}

	GetRenamed(data, methodText)
		{
		return .Rename(.MethodText.Replace(.line, ''), data.from, data.to)
		}

	Process(data)
		{
		i = .MethodRange.from
		n = .MethodRange.to - i
		text = .Rename(data.text[i :: n].Replace(.line, ''), data.from, data.to)
		data.text = data.text.ReplaceSubstr(i, n, text)
		return true
		}

	Change(member/*unused*/)
		{ }
	Errors(data/*unused*/)
		{ return "" }
	Warnings(data/*unused*/)
		{ return "" }
	}
