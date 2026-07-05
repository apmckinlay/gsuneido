// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Comby Playgroud'
	Controls()
		{
		return Object('Vert',
			Object('Horz',
				Object(#WorkSpaceCode, name: "source")
				Object('Vert',
					Object('Horz'
						Object('Edit', name: 'search', width: 60),
						 'Skip',
						Object('Button', 'Find'))
					Object('Static', 'Double click lines to locate the result')
					Object('WorkSpaceCode', name: 'results', readonly:))))
		}
	On_Find()
		{
		search = .FindControl('search').Get()
		source = .FindControl('source').Get()
		results = .FindControl('results')
		results.Set('')
		if search.Blank?() or source.Blank?()
			return

		try
			{
			matches = CombyMatch(search, source)
			if matches.Size() is 0
				results.Set('No matches found.')
			else
				{
				output = ''
				for m in matches
					{
					output $= 'Match at position ' $ m.pos $ ' to ' $ m.end $ '\n'
					for hole_name, hole_val in m.holes
						output $= '  :[' $ hole_name $ '] = ' $ hole_val $ '\n'
					output $= '=============================\n'
					}
				results.Set(output)
				}
			}
		catch (e)
			{
			results.Set('Error: ' $ e)
			}
		}

	Scintilla_DoubleClick(source)
		{
		if source isnt .FindControl('results')
			return

		results = source
		text = results.Get()
		if text.Blank?()
			return

		lineNum = results.LineFromPosition()
		lineText = text.NthLine(lineNum)

		// Find nearest "Match at position" line above
		while lineNum >= 0 and not lineText.Prefix?('Match at position')
			{
			--lineNum
			if lineNum < 0
				return
			lineText = text.NthLine(lineNum)
			}

		// Parse position, format: "Match at position N to M"
		parts = lineText.Split(' ')
		if parts.Size() < 4
			return
		pos = Number(parts[3])

		source = .FindControl('source')
		sourceText = source.Get()
		line = sourceText.LineFromPosition(pos)

		source.GotoLine(line)
		}
	}
