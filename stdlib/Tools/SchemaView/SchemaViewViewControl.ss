// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
CodeViewControl
	{
	New(data)
		{
		super(:data, addons: .addons(data.type), readonly:)
		.Editor.SetWrap(true)
		}

	addons(type)
		{
		return Object(Addon_suneido_style: type is #column ? true : #(query:),
			Addon_wrap:, Addon_calltips: false,
			Addon_class_outline: false,
			Addon_folding: false,
			Addon_go_to_line: false,
			Addon_highlight_cursor_line: false,
			Addon_indent_guides: false,
			Addon_show_line_numbers: false,
			Addon_show_margin: false,
			Addon_status: false)
		}

	InitialSet(data)
		{
		data.text = data.type is #table
			? .tableText(data.name)
			: data.type is #view ? .viewText(data.name) : .columnText(data.name)
		super.InitialSet(data)
		}

	tableText(table)
		{
		return Schema(table) $ .tableForeignKeys(table) $ .tableCode(table)
		}

	tableForeignKeys(table)
		{
		cascadeMode = 3
		fkeys = ""
		QueryApply("indexes where fktable is " $ Display(table))
			{|x|
			type = x.key is true ? #key : x.key is 'u' ? #unique : #index
			fkeys $= '\t' $ x.table $ ' ' $ type $ '(' $ x.columns $ ") in " $ x.fktable $
				(x.fkmode is cascadeMode ? " cascade" : "")
			if x.columns isnt x.fkcolumns
				fkeys $= '(' $ x.fkcolumns $ ')'
			fkeys $= "\r\n"
			}
		return Opt("\r\nForeign Keys\r\n", fkeys)
		}

	tableCode(table)
		{
		trigger = "Trigger_" $ table
		triggers = Object()

		className = "Table_" $ table
		classes = Object()
		for lib in Libraries()
			{
			if not QueryEmpty?(lib, name: trigger)
				triggers.Add(lib $ ':' $ trigger)
			QueryApply(lib $ " where name in (" $ Display(className) $ ", " $
				Display(lib.Capitalize() $ '_' $ className) $ ") and group is -1")
				{
				classes.Add(lib $ ':' $ it.name)
				}
			}
		return Opt("\r\nTriggers\r\n\t", triggers.Join("\r\n\t"), "\r\n") $
			Opt("\r\nTable Definitions\r\n\t", classes.Join("\r\n\t"))
		}

	viewText(view)
		{
		if false is x = Query1(#views, view_name: view)
			return [name: view]
		query = x.view_definition
		return FormatQuery(query) $ "\n\nStrategy =========\n\n" $
			QueryStrategyAndWarnings(view)
		}

	columnText(column)
		{
		return .columnDatadict(column) $ .columnRule(column)
		}

	columnDatadict(column)
		{
		b = dd = Datadict(column)
		if dd is Field_string
			return ""
		s = ""
		do
			{
			s $= Display(b) $ '\n'
			} while Class?(b = b.Base())
		s $= '\n'
		fields = #(Prompt, SelectPrompt, Heading, Control, Format)
		promptWidth = 12
		for m in fields
			try
				s $= m.LeftFill(promptWidth) $ ": " $
					String(dd[m]).Replace('\n', "\\\\n") $ "\r\n"
		return s
		}

	columnRule(column)
		{
		rule = "Rule_" $ column
		try
			{
			fn = Global(rule)
			return '\n' $ rule $ "\n\n" $ SourceCode(fn)
			}
		catch
			return ""
		}

	Scintilla_DoubleClick()
		{
		line = .Editor.GetLine()
		if line =~ "^\w+:[[:upper:]][_[:alnum:]]*[?!]?\>"
			{
			_hwnd = .Window.Hwnd
			.Editor.Home()
			.Editor.LineEndExtend()
			.Editor.On_Context_Go_To_Definition()
			}
		}
	}
