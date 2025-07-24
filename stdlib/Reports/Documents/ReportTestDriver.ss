// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Document_Builder
	{
	New(trackAccesses = false)
		{
		super()
		.fmtOb = Object()
		if trackAccesses is true
			.fmtOb.Accesses = Object()
		}

	PlainText?() { return true }

	GetLineSpecs(@unused)
		{
		return Object(height: 10, descent: 3)
		}

	GetTextWidth(@unused)
		{
		return 10  /*= fake size info for testing */
		}

	GetCharWidth(@unused)
		{
		return 10 /*= fake size info for testing */
		}

	GetTextHeight(@unused)
		{
		return 10 /*= fake size info for testing */
		}

	Process(fmt)
		{
		if not Instance?(fmt)
			return
		csv = fmt.ExportCSV()
		if String?(csv)
			{
			str = fmt.ExportCSV().Trim()

			if str isnt ''
				.fmtOb.Add(str)
			}
		else if Object?(csv)
			.fmtOb.Merge(csv)

		if .fmtOb.Member?('Accesses')
			.verifyAccessPoints(fmt)
		}

	verifyAccessPoints(fmts)
		{
		for idx in fmts.MembersIf({ it.Suffix?('_formats') })
			for fmt in fmts[idx]
				if fmt.Member?('Field') and fmt.Field isnt false
					for i in fmts.Data.MembersIf({ it is '_access_' $ fmt.Field })
						.fmtOb.Accesses.Add(fmts.Data[i])
		}

	Finish(@unused)
		{
		return .fmtOb
		}
	}
