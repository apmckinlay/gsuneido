// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 550
	Ymin: 150
	Title: "View History"
	New(fields = #(), data = false)
		{
		super(.layout(data))
		.listview = .FindControl('viewhistory')
		.load_list_records(fields, data)
		.listview.SetReadOnly(true, grayOut: false)
		}
	columns: (curhis_action, curhis_date, curhis_user, curhis_comment)
	layout(data)
		{
		.title = data isnt false ? .Title $ ' - ' $ data.transaction_type : .Title
		return Object('Vert'
			Object('ListStretch',
				columns: .columns,
				columnsSaveName: .title,
				name: "viewhistory",
				)
			)
		}
	load_list_records(fields, data)
		{
		.list_records = Object()
		.add_list_records(fields, data)
		.list_records.Sort!({ |x,y| x.timestamp < y.timestamp })
		for rec in .list_records
			.listview.AddRow(rec)
		}
	add_list_records(fields, data)
		{
		for action in fields.Members()
			{
			if fields[action].Member?('query')
				{
				.add_list_rec_from_table(action, fields[action], data)
				continue
				}

			if not Date?(data[fields[action].date])
				continue

			rec = Object()
			rec.timestamp = data[fields[action].date]
			rec.curhis_date = (data[fields[action].date]).ShortDateTime()
			rec.curhis_action = action
			user = fields[action].Member?('user')
				? .format_value(data[fields[action].user], fields[action].user)
				: ''
			rec.curhis_user = user isnt 'default' ? user : ''
			rec.curhis_comment = fields[action].Member?('comment')
				? .format_value(data[fields[action].comment], fields[action].comment)
				: ''
			.list_records.Add(rec)
			}
		}
	add_list_rec_from_table(action, ob, data)
		{
		QueryApply(data[ob.query])
			{ |x|
			if not Date?(x[ob.date])
				continue

			rec = Object()
			rec.timestamp = x[ob.date]
			rec.curhis_date = (x[ob.date]).ShortDateTime()
			rec.curhis_action = action
			user =  ob.Member?(#user) ? .format_value(x[ob.user], ob.user) : ''
			rec.curhis_user = user isnt 'default' ? user : ''
			rec.curhis_comment = ob.Member?('comment')
				? .format_value(x[ob.comment], ob.comment)
				: ''
			.list_records.Add(rec)
			}
		}
	format_value(value, field = false)
		{
		if Date?(value)
			value = value.NoTime() is value ? value.ShortDate() : value.ShortDateTime()
		else if Number?(value)
			value = value.Format(Datadict(field).Format.mask)
		else if Boolean?(value)
			value = value is true ? 'Yes' : 'No'
		return value
		}
	}
