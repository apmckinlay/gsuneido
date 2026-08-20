// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// Interactive test tool for AstFormatter - pick a library record,
// then see the difference between original and formatted code.
Controller
	{
	Title: 'AstFormatter Tester'

	Controls()
		{
		return Object('Vert'
			Object('Horz'
				'LibLocate'
				#(Button, 'Go To Random')
				)
			#(Diff2Code, '', '', 'Original', 'Formatted', comment: false)
			)
		}

	New()
		{
		.locate = .FindControl('LibLocate')
		.diff = .FindControl('Diff')
		}

	// Received from LibLocate as "lib:Name" e.g. "stdlib:AstFormatter"
	Locate(item)
		{
		name = item.AfterFirst(':')
		lib = item.BeforeFirst(':')
		.formatAndDiff(name, lib)
		}

	formatAndDiff(name, lib)
		{
		rec = Query1(lib $ ' where name is ' $ Display(name) $ ' and group is -1')
		if rec is false
			{
			.FindControl('status').Set('Record not found: ' $ name)
			return
			}
		src = rec.text
		formatted = AstFormatter(src)
		.diff.UpdateList(src, formatted)
		}

	On_Go_To_Random()
		{
		rec = RandomLibraryRecord()
		lib = rec.lib
		name = rec.name
		.locate.Set(name)
		.formatAndDiff(name, lib)
		}
	}
