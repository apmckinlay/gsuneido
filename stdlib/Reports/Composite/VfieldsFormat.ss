// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// WARNING: must pass data to New
GridFormat
	{
	New(@args)
		{
		super(@.format(args))
		}
	format(args)
		{
		data = args.data
		rows = Object()
		access = args.GetDefault('access', false)
		grid_access = Object()
		maxWidth = 0
		for (i = 0; args.Member?(i); ++i)
			{
			rows[i] = Object()
			arg = args[i]

			field = Object?(arg) ? arg.field : arg
			dd = Datadict(field)
			heading = Heading(field) $ ": "

			// handle info fields
			info_value = ''
			if dd.Base?(Field_info)
				{
				heading = data[field].BeforeFirst(':').Trim()
				if heading isnt ''
					heading $= ': '
				info_value = StripInfoLabel(data[field]).AfterFirst(':').Trim()
				}
			row_access = Object()

			prompt = Object("Text", heading $ " ", justify: 'right')
			if (args.Member?('promptfont'))
				prompt.font = args.promptfont
			else if (args.Member?('font'))
				prompt.font = args.font
			rows[i].Add(prompt)
			row_access.Add(false)

			value = Object?(dd.Format) ? dd.Format.Copy() : Object(dd.Format)
			if Object?(arg) and arg.Member?('wrap')
				value[0] = 'Wrap'
			if Object?(arg) and arg.Member?('width')
				value.width = arg.width
			if value.Member?('width')
				maxWidth = Max(maxWidth, value.width)

			value.data = info_value isnt '' ? info_value : data[field]
			if args.Member?('justify')
				value.justify = args.justify
			if (args.Member?('font'))
				value.font = args.font
			rows[i].Add(value)
			row_access.Add(access isnt false and access.Size() > i ? access[i] : false)

			grid_access.Add(row_access)
			}
		.applyMaxWidth(maxWidth, rows)
		return Object(rows, top:, access: grid_access)
		}
	applyMaxWidth(maxWidth, rows)
		{
		if maxWidth is 0
			return
		for row in rows
			if row[1].Member?('width')
				row[1].width = maxWidth
		}
	}
