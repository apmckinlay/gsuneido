// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		.input = .FindControl('input')
		.diff = .FindControl('diff')
		}
	Controls: (Vert
		(CodeControl name: input)
		(Diff2Control "", "", "", "", " Original", " Formatted", name: diff)
		)
	EN_CHANGE(source)
		{
		if source isnt .input
			return
		src = .input.Get()
		.diff.UpdateList(src, FmtCode(src))
		}
	}