// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
CommandParent
	{
	Title: '' // used to store settings, should not contain data like Order#
	HorzLeftCtrls: ''
	HorzRightCtrls: ''
	BottomButtons: ''
	Columns: ()
	ProtectField: false // ready only by default
	Table: ''
	NumField: ''
	StretchColumn: true
	AccessScreen: false
	AccessKeyField: false
	DefaultColumns: false
	CheckAboveSortLimit?: true

	CallClass(@args)
		{
		args.Add(this at: 0)
		ModalWindow(args, title: .Title, onDestroy: args.GetDefault('onDestroy', false))
		}

	New(where = "", enableMultiSelect = false,
		.title = '', // this title can contain data like Order#
		select = #())
		{
		super(.Layout(:enableMultiSelect, :select))
		.List = .FindControl('List')
		.Data.AddObserver(.Record_changed)
		.Rec = .Data.Get()
		.Where = where
		.Load(where)
		}

	ContextMenu: (Copy, Zoom)
	ListQuery: false
	DefaultSelect: #()
	HeaderSelectPrompt: false
	Layout(enableMultiSelect = false, select = #())
		{
		if select is #()
			select = .Val_or_func('DefaultSelect')
		menu = .ContextMenu.Copy()
		if .AccessKeyField isnt false
			menu.Add('Access', at: 0)
		return Object('Record',
			Object('Vert',
				Object('VirtualList', .ListQuery,
					columns: .Val_or_func('Columns'), protectField: .ProtectField,
					columnsSaveName: .Title,
					:menu, startLast:,
					title: .title is '' ? .Title : .title,
					:enableMultiSelect, name: 'List',
					filtersOnTop:, :select,
					headerSelectPrompt: .HeaderSelectPrompt,
					defaultColumns: .DefaultColumns,
					disableCheckSortLimit?: not .CheckAboveSortLimit?,
					stretchColumn: .StretchColumn),
				.BottomButtons
				)
			)
		}

	ToggleFilter()
		{
		.List.ToggleFilter()
		}

	Record_changed(member /*unused*/) { }

	Load(where)
		{
		if '' isnt query = .Query(where)
			.List.SetQuery(query, .Val_or_func('Columns'))
		}
	SetColumns(cols)
		{
		.Columns = cols
		}

	Query(where /*unused*/ = "")
		{ return ''}

	On_Refresh()
		{ .List.Refresh() }

	CheckImport()
		{
		if .Destroyed?()
			return false
		if .NeedImport?() and
			YesNo("Messages have not been imported since " $
					.LastImport().LongDateTime() $ "\n\n" $
				"The automatic importing may not be working.\n\n" $
				"Import now?",
				.Title, .Window.Hwnd, MB.ICONWARNING)
			.Import()
		return true
		}

	NeedImport?()
		{
		lastimport = .LastImport()
		return Date?(lastimport) and lastimport.Plus(minutes: 10) < Timestamp()
		}

	LastImport()
		{
		return .Table is '' or false is (x = Query1(.Table)) ? "" : x[.NumField]
		}

	Import()
		{
		SuneidoLog("WARNING: manual import - scheduler not working?")
		.ImportFunc()
		.On_Refresh()
		}

	ImportFunc() { }

	On_Context_Copy(rec, col)
		{
		if rec is false or col is false
			return
		value = rec[col]
		if not String?(value)
			value = Display(value)
		ClipboardWriteString(value)
		}

	On_Context_Zoom(rec, col)
		{
		if rec isnt false and col isnt false and
			false isnt ctrlClass = GetControlClass.FromField(col)
			ctrlClass.ZoomReadonly(FormatValue(rec[col], col))
		}

	Commands: #(#("Select",	"Alt+S"))
	On_Select()
		{
		.List.On_Select()
		}
	AccessGoTo_CurrentBookOption() // need to make the help button work when accessing
		{
		book_option = 0
		if false is currentbook = Suneido.GetDefault('CurrentBook', false)
			return book_option

		option = String(this).RemoveSuffix(`()`)
		if false isnt rec = QueryFirstBookOption(currentbook, option)
			book_option = rec.path $ '/' $ rec.name
		return book_option
		}

	VirtualList_SaveSort?()
		{
		return false
		}

	VirtualList_DoubleClick(rec, col /*unused*/)
		{
		if rec is false
			return false

		if .AccessKeyField is false
			return 0

		AccessGoTo(.AccessScreen, .AccessKeyField, rec[.AccessKeyField], .Window.Hwnd,
			onDestroy: .List.Refresh)
		return false
		}

	On_Context_Access()
		{
		.VirtualList_DoubleClick(.List.GetSelectedRecord(), '')
		}
	}