// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// TODO: if column is capitalized fetch definition from libraries
class
	{
	New()
		{
		.tables = Object()
		QueryApply('views')
			{ |x|
			.tables.Add([name: x.view_name, group: false])
			}
		QueryApply('tables')
			{|x|
			.tables.Add([name: x.table, group: true])
			}
		.tables.Sort!(By(#name))
		i = 0
		for x in .tables
			x.num = ++i
		}
	colstart: 10000
	Children(parent)
		{
		if parent is 0
			return .tables
		if parent >= .colstart
			return []
		x = .tables[parent - 1]
		if not x.group
			return []
		items = Object()
		cols = .columns(x.name)
		for i in cols.Members()
			items.Add([name: cols[i], num: parent * .colstart + i, group: false])
		return items
		}
	Get(i)
		{
		return i >= .colstart
			? .getColumn(i)
			: .tables[--i].group
				? .getTable(i)
				: .getView(i)
		}
	getView(i)
		{
		view = .tables[i].name
		if false is x = Query1("views", view_name: view)
			return [name: view]
		query = x.view_definition
		text = FormatQuery(query) $
			'\n\nStrategy =========\n\n' $
			QueryStrategyAndWarnings(view)
		return Object(name: view, :text, group: false, table: view)
		}
	getTable(i)
		{
		table = .tables[i].name
		return Object(text: Schema(table) $ .get_foreignkeys(table) $
			.get_trigger(table) $ .getTableClasses(table),
			group: true, :table, name: table)
		}
	getColumn(i)
		{
		table = .tables[(i / .colstart).Int() - 1].name
		col = .columns(table)[i % .colstart]
		return Object(text: .get_dd(col) $ .get_rule(col), group: false,
			name: col, :table)
		}
	get_dd(col)
		{
		b = dd = Datadict(col)
		if dd is Field_string
			return ""
		s = ""
		do
			s $= Display(b) $ '\n'
			while Class?(b = b.Base())
		s $= '\n'
		fields = #(Prompt, SelectPrompt, Heading, Control, Format)
		promptWidth = 12
		for m in fields
			try s $= m.LeftFill(promptWidth) $ ': ' $
				String(dd[m]).Replace('\n', '\\\\n') $ "\r\n"
		return s
		}
	get_foreignkeys(table)
		{
		cascadeMode = 3
		fkeys = ""
		QueryApply('indexes where fktable is ' $ Display(table))
			{ |x|
			fkeys $= '\t' $ x.table $ ' index(' $ x.columns $ ') in ' $
				x.fktable $ (x.fkmode is cascadeMode ? ' cascade' : "")
			if (x.columns isnt x.fkcolumns)
				fkeys $= '(' $ x.fkcolumns $ ')'
			fkeys $= '\r\n'
			}
		return fkeys is "" ? "" : '\r\nForeign Keys\r\n' $ fkeys
		}
	get_trigger(table)
		{
		s = "\n"
		trigger = 'Trigger_' $ table
		for lib in Libraries()
			if not QueryEmpty?(lib, name: trigger)
				s $= lib $ ':' $ trigger $ '\n'
		return s
		}
	getTableClasses(table)
		{
		s = '\n'
		name = 'Table_' $ table
		for lib in Libraries()
			QueryApply(lib $ ' where name in (' $ Display(name) $ ', ' $
				Display(lib.Capitalize() $ '_' $ name) $ ') and group is -1')
				{
				s $= lib $ ':' $ it.name $ '\n'
				}
		return s
		}
	get_rule(col)
		{
		rule = 'Rule_' $ col
		try
			{
			fn = Global(rule)
			return '\n' $ rule $ '\n\n' $ SourceCode(fn)
			}
		catch
			return ""
		}
	columns(table)
		{
		return QueryList('columns where table = ' $ Display(table), 'column')
		}

	Children?(parent/*unused*/)
		{ return true }
	Container?(i)
		{ return i < .colstart and .tables[i-1].group }
	Static?(x/*unused*/)
		{ return true }
	Nextnum()
		{ }
	NewItem(x/*unused*/)
		{ }
	Update(x/*unused*/)
		{ }
	DeleteItem(num/*unused*/)
		{ }
	}
