// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xstretch: 1
	Ystretch: 1
	Title: "PluginsView"
	New()
		{
		data = .get()
		.list = .Vert.List
		.list.Set(data)
		.list.SetColWidth(.columns.Find(#From), width: 250)
		}
	columns: (Plugin, 'Extension Point', From, Contribution)
	Controls()
		{
		return Object('Vert'
			#(Form
				'Skip', (Static "Plugin > contains > ")
					(Field, name: "pluginWhere", group: 1) (Button "Find") 'nl'
				'Skip', (Static "Extension Point > contains > ")
					(Field, name: "epWhere", group: 1)
				)
			Object('ListStretch', .columns, defWidth: 160,
				headerSelectPrompt: 'no_prompts')
			#(Skip 3)
			#(Horz
				Fill
				(Button 'Go To' defaultButton:, tip: '(same as double click)' pad: 30)
				Fill)
			#(Skip 3)
			)
		}
	List_WantNewRow()
		{ return false }
	List_WantEditField(col)
		{
		return .columns[col] is 'Contribution'
			? Object('Editor', xmin: .list.GetColWidth(2), height: 19)
			: false
		}
	List_DoubleClick(row, col)
		{
		if .columns[col] is 'Contribution' or row is false
			return 0
		.On_Go_To()
		return true
		}
	get(pluginWhere = "", epWhere = "")
		{
		ob = Object().Set_default(Object())
		Plugins().ForeachContrib()
			{ |c|
			plugin = c[0]
			ep = c[1]
			if pluginWhere isnt "" and plugin !~ "(?i)" $ pluginWhere
				plugin = false
			if epWhere isnt "" and ep !~ "(?i)" $ epWhere
				ep = false
			from = c.from
			c = c.Copy()
			c.Delete(1, 0, #from)
			if plugin isnt false and ep isnt false
				ob.Add([Plugin: plugin, 'Extension Point': ep, From: from,
					Contribution: c])
			}
		return ob.Sort!(By('Plugin', 'Extension Point', 'From'))
		}
	On_Find()
		{
		pluginWhere = .FindControl("pluginWhere").Get()
		epWhere = .FindControl("epWhere").Get()
		data = .get(pluginWhere, epWhere)
		.list.Set(data)
		}
	On_Go_To()
		{
		if false isnt x = .get_selected()
			GotoLibView(x.From)
		}
	get_selected()
		{
		sel = .list.GetSelection()
		if sel.Size() isnt 1
			{
			.AlertInfo(.Title, "No row selected")
			return false
			}
		data = .list.Get()
		return data[sel[0]]
		}
	}