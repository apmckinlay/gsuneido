// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
QueryFormat
	{
	CallClass(query, title = "", excludeFields = #())
		{
		if title is ""
			title = query
		query = QueryStripSort(query)
		fields = QueryColumns(query)
		sf = SelectFields(fields, excludeFields)

		return Object('Params',
			Object(this, Sf: sf, Basequery: query),
			name: "Global_Print - " $ title,
			validField: "globalreport_valid",
			SetParams: Object(:sf),
			Params: Object(.paramsCtrl, title, sf)
			)
		}

	Header()
		{
		font = #(name: 'Arial', size: 6.5)
		return Object('Vert',
			Object('PageHead', .Params.title),
			Object('Text', 'Sort By: ' $ .Params.sort, :font)
			'Vskip'
			super.Header())
		}

	Query()
		{
		.fields = .Sf.PromptsToFields(.Params.fields)
		if .fields.Empty?()
			.fields = QueryColumns(.Basequery)
		query = .Basequery $ .Sf.Joins(.Params.fields $ ',' $ .Params.sort)
		sortfields = .Sf.PromptsToFields(.Params.sort)
		if not sortfields.Empty?()
			query $= ' sort ' $ sortfields.Join(',')
		return query
		}

	Output()
		{
		return .format
		}

	getter_format()
		{
		return .format = Object("Row").Add(@.fields) // once only
		}

	Total()
		{
		return .Sf.PromptsToFields(.Params.total_fields)
		}

	After(data)
		{
		format = Object('_output')
		total_fields = .Sf.PromptsToFields(.Params.total_fields.Split(','))
		// wrap field format in Total format
		for field in total_fields
			{
			field_fmt = Datadict(field).Format.Copy()
			field_fmt.data = data['total_' $ field]
			f = Object("Total", field_fmt)
			format.Add(f, at: field)
			}
		return format
		}

	paramsCtrl: PassthruController
		{
		Xstretch: 1
		Ystretch: 1
		New(title, sf)
			{
			super(.controls(sf))
			.title = title
			}

		controls(sf)
			{
			fields = sf.Prompts().Sort!()
			return Object('Horz',
				Object('Vert'
					Object('Pair',
						Object('Static', 'Title')
						Object('Field', name: "title"))
					Object('Pair',
						Object('Static', 'Columns')
						Object('ChooseTwoList', fields, title: 'Columns', name: "fields"))
					Object('Pair',
						Object('Static', 'Sort By')
						Object('ChooseTwoList', fields, title: 'Sort By', name: "sort"))
					Object('Pair',
						Object('Static', 'Total Fields')
						Object('ChooseMany', fields, name: "total_fields"))
					)
				'Skip'
				#(Vert
					(Button 'Save' xstretch: 0)
					(Button 'Load...' xstretch: 0)
					(Button 'Clear' xstretch: 0)
					Fill
					)
				)
			}

		On_Save()
			{
			data = .Send('Get')
			if data.title is ""
				{
				Alert("Please enter a Title")
				return
				}
			save_as = "Global_Print - " $ .title $ " - " $ data.title
			QueryDo("delete params where report is " $ Display(save_as))

			params = data.Project(@.save_fields)
			QueryOutput("params", [user: Suneido.User, report: save_as, :params])
			}

		save_fields: (title fields sort total_fields)
		On_Load()
			{
			if false is rpt = ToolDialog(.Window.Hwnd, Object(.load_ctrl, .title))
				return
			x = Query1("params", report: rpt)
			data = .Send("Get")
			for field in .save_fields
				data[field] = x.params[field]
			.Send("Set", data)
			}

		load_ctrl: Controller
			{
			Title: "Load Report"
			New(title)
				{
				.list = .Vert.ListBox
				.prefix = "Global_Print - " $ title $ " - "
				QueryApply("params where report =~ '^(?q)" $ .prefix $ "'")
					{|x|
					.list.AddItem(x.report[.prefix.Size() ..])
					}
				}

			Controls:
				(Vert
					(ListBox xmin: 200)
					(Skip 5)
					(Horz
						Fill
						(Button Load)
						Skip
						(Button Cancel)
						Fill
					)
				)

			ListBoxDoubleClick(sel /*unused*/)
				{
				.On_Load()
				}

			On_Load()
				{
				sel = .list.GetCurSel()
				if sel is -1
					{
					Alert("Please select a report.")
					return
					}
				.Window.Result(.prefix $ .list.GetText(sel))
				}
			}

		On_Clear()
			{
			data = .Send("Get")
			for field in .save_fields
				data[field] = ""
			.Send("Set", data)
			}
		}
	}