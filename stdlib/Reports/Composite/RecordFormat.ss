// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
GridGenFormat
	{
	New(rec, fields = false, font = false)
		{
		super(.layout(rec, fields), :font)
		}

	layout(rec, fields)
		{
		formats = Object()
		if fields is false
			fields = rec.Members()

		for field in fields
			{
			prompt = Object('Text' SelectPrompt(field) $ ': ', justify: 'right')
			format = Datadict(field).Format
			format = String?(format) ? Object(format) : format.Copy()
			format[1] = rec[field]
			if format[0] is 'Wrap'
				{
				// most wrap widths are narrow for printing in columns
				// but here we can use wider
				w = 4.InchesInTwips()
				// use WrapFormat to wrap
				wf = _report.Construct(format)
				font = _report.GetFont()
				lines = wf.WrapDataLines(format[1], w, font).Split('\n')
				formats.Add(
					Object(prompt, Object('Text', lines.GetDefault(0, ''))))
				for (i = 1; i < lines.Size(); ++i)
					formats.Add(
						Object(#(Text ''), Object('Text', lines[i])))
				}
			else
				{
				// REFACTOR: should not reference Accountinglib code
				if format[0] is "GlAcctDept"
					{
					format[1] = Record()
					format[1][field] = rec[field]
					deptfield = field.Replace("glacct_num", "gldept_id")
					format[1][deptfield] = rec[deptfield]
					format[2] = field.AfterFirst("glacct_num") //suffix
					}
				formats.Add(Object(prompt, format))
				}
			}
		return formats
		}
	}