// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "QueryView"

	CallClass()
		{
		GotoPersistentWindow('QueryViewControl', QueryViewControl)
		}
	New()
		{
		.edit = .FindControl('code')
		.Redir('On_Find')
		.Redir('On_Find_Next')
		.Redir('On_Find_Previous')
		.Redir('On_Replace')

		.results = .FindControl('results')
		.q_or_c = .FindControl('q_or_c')
		.asof = .FindControl('asof')
		.status = .Vert.Status
		.Redir('On_Select', this)
		.Defer(.edit.SetFocus)
		}
	Menu:
		(
		("&File",
			"&Close")
		("&Edit",
			"&Undo", "&Redo", "", "Cu&t", "&Copy", "&Paste", "&Delete" ""
			"&Find...", "Find &Next", "Find &Previous", "R&eplace...")
		("&Query",
			"&Run", "&Strategy", "&Profile", "&Keys", "Schema", "&Count", "&Hash",
			"&Format", "Show Format")
		)
	Controls()
		{
		.menu = #(Copy, 'Copy Record', 'Copy Field:Value', 'Inspect Value',
			'Inspect Record', 'Print Record', '',	'Print', 'Export')
		return ['Vert',
			['Horz'
				#(Toolbar, Split, "",
					Cut, Copy, Paste, "", Undo, Redo, "", New_Query_View, "",
					Keys, Strategy, Schema, Count, Hash, Format,
					Benchmark, Profile, Run, ">"),
				.extraControls()
				],
			.vertLayout,
			'Statusbar']
		}

	extraControls()
		{
		extra = Object(#Vert, #(EtchedLine, before: 0 after: 0), horz = Object(#Horz))
		horz.Add(#(RadioButtons Auto Query Cursor name: q_or_c horz:), #Skip)
		horz.Add(#Skip, #(StaticText 'As of'), #Skip, #(ChooseDateTime, name: asof),
			#(EnhancedButton command: 'ClearAsof', image: 'cross.emf',
				imageColor: 0x737373, mouseOverImageColor: 0x0000ff,
				imagePadding: 0.15, tip: 'clear As of'))
		return extra
		}
	On_ClearAsof()
		{
		if .asof isnt false
			.asof.Set("")
		}

	listLayout()
		{
		return ['VirtualList', headerSelectPrompt: 'no_prompts', enableMultiSelect:,
			menu: .menu, name: 'list_browse', preventCustomExpand?:,
			headerMenu: #('Copy', 'Go To Definition', 'Show References', ''),
			xmin: 100, ymin: 100]
		}

	browse: false
	switchToList()
		{
		.switchTo(.listLayout())
		.browse = .FindControl('list_browse')
		}
	switchToText(text)
		{
		.switchTo([#ScintillaIDE, readonly:, set: text, xmin: 100, ymin: 100])
		}
	switchTo(control)
		{
		.browse = false
		.results.Remove(0)
		.results.Append(control)
		}
	Commands:
		(
		(Find,			"Ctrl+F",	"Find text")
		(Find_Next,		"F3",		"Find the next occurrence")
		(Find_Previous,	"Shift+F3",	"Find the previous occurrence")
		(Replace,		"Ctrl+H",	"Find and replace text in the current item")
		(Select_Line,	"F3",		"Extend select up one line")
		(New_Query_View, "", 		"Open a new QueryView", 'hsplit')
		(Run,			"F9",		"Execute the current selection", '!')
		(Strategy, 		"Shift+F9",	"Show how the query will be executed", '?')
		(Profile,       "",         "Run QueryProfiler", 'P')
		(Schema, 		"Ctrl+S",	"Show Schema View or query columns", 'S')
		(Keys,			"",			"Show the keys for the selected query", 'K')
		(Users_Manual,	"F1")
		(Select, 		"Alt+S",	"Open Select Filter")
		(Count,			"Alt+C",	"summarize count", 'C')
		(Hash,			"Alt+H",	"hash the query results", 'H')
		(Format,		"Alt+F",	"format the query (Shift to show)", 'F')
		(Show_Format,	"Shift+Alt+F",	"show the formatted query")
		(Benchmark,		"Ctrl+B",	"Benchmark QueryApply", 'B')
		(Split,			"",			"toggle vertical/horizontal split", '+')
		)
	getRunText(bare = false)
		{
		s = GetRunText(.edit)
		if s is ""
			.setstatus_invalid("please select something")
		else if not bare
			s = QuerySuppress(s)
		return s
		}

	split: 'vert'
	On_Split()
		{
		text = .edit.Get()
		.Vert.Remove(1)
		layout = .split is 'vert' ? .horzLayout : .vertLayout
		.split = #(vert: horz, horz: vert)[.split]
		.Vert.Insert(1, layout)
		.edit = .FindControl('code')
		.results = .FindControl('results')
		.edit.Set(text)
		}
	vertLayout: (VertSplit,
		(Horz // to get rid of Top from ScintillaAddonsControl
			(QueryCode, name: code, xmin: 100, ymin: 100)),
		('Vert', (Fill, xmin: 100, ymin: 100), name: 'results')
		name: 'split')
	horzLayout: (HorzSplit,
		(Vert // to get rid of Top from ScintillaAddonsControl
			(QueryCode, name: code, xmin: 100, ymin: 100)),
		('Horz', (Fill, xmin: 100, ymin: 100), name: 'results')
		name: 'split')

	On_Run()
		{
		if "" is s = .getRunText()
			return
		try
			switch (first = .firstToken(s))
				{
			case 'create', 'alter', 'ensure', 'rename', 'drop', 'destroy',
				'view', 'sview' :
				Database(s)
				.setstatus_valid("SUCCEEDED: " $ s)
			case 'insert', 'update', 'delete' :
				n = QueryDo(s)
				action = first.Replace("e$", "") $ "ed "
				.setstatus_valid("SUCCEEDED: " $ action $ n $
					" record" $ (n isnt 1 ? 's' : ''))
			default :
				.run(s)
				}
		catch (err)
			.setstatus_invalid(err)
		}
	warnings(query)
		{
		s = FormatQuery(query)
		w = ""
		if s.Has?("project /*NOT UNIQUE*/") and
			not query.Has?("CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE")
			w $= ", PROJECT NOT UNIQUE"
		if s.Has?("union /*NOT DISJOINT*/") and
			not query.Has?("CHECKQUERY SUPPRESS: UNION NOT DISJOINT")
			w $= ", UNION NOT DISJOINT"
		if s.Has?("join /*MANY TO MANY*/") and
			not query.Has?("CHECKQUERY SUPPRESS: JOIN MANY TO MANY")
			w $= ", JOIN MANY TO MANY"
		return w[2..]
		}
	run(s)
		{
		warnings = .warnings(.getRunText(bare:))
		.query = s
		t = Timer() { .setQuery('/*SLOWQUERY SUPPRESS*/\n' $ s) }
		mode = .browse.UsingCursor?() ? 'Cursor' : 'Query'
		if warnings isnt ""
			.setstatus_invalid(warnings)
		else
			.setstatus_valid(
				mode $ " took " $ t.Format('###,###.###') $
				' seconds to load the first screen-full')
		}

	setQuery(s)
		{
		.switchToList()
		useQuery = #(Query: true, Cursor: false, Auto: auto)[.q_or_c.Get()]
		try
			.browse.SetQuery(s, asof: .asof(), :useQuery)
		catch(err)
			{
			.switchToList()
			throw err
			}
		}

	asof()
		{
		if .asof is false
			return false
		return Date?(asof = .asof.Get()) ? asof : false
		}

	firstToken(s)
		{
		scanner = Scanner(s)
		for token in scanner
			if scanner.Type() not in (#COMMENT, #WHITESPACE, #NEWLINE)
				return token
		return ''
		}
	setstatus_invalid(msg)
		{
		.status.SetBkColor(CLR.ErrorColor)
		.status.Set(' ' $ msg)
		}
	setstatus_valid(msg = "")
		{
		.status.SetBkColor(false)
		.status.Set(' ' $ msg)
		}
	On_New_Query_View()
		{
		PersistentWindow(QueryViewControl)
		}
	On_Strategy()
		{
		if "" is s = .getRunText()
			return
		q_or_c = .q_or_c.Get()
		.switchTo([QueryStrategyViewer, s, :q_or_c])
		}
	On_Profile()
		{
		if "" is s = .getRunText(bare:)
			return
		QueryProfiler(s)
		}
	On_Benchmark()
		{
		if "" is s = .getRunText(bare:)
			return

		q_or_c = .q_or_c.Get()
		if q_or_c is 'Cursor'
			block =
				{
				Cursor(s)
					{|c|
					Transaction(read:)
						{|t|
						while false isnt c.Next(t)
							{}
						}
					}
				}
		else
			block = { QueryApply(s, {|unused| }) }
		.setstatus_valid("using " $ (q_or_c is "Cursor" ? "Cursor" : "Query"))
		Working("Benchmark")
			{ b = Bench(block) }
		.switchToText(b)
		}
	On_Select()
		{
		if .browse isnt false
			.browse.On_Select()
		}
	On_Keys()
		{
		if "" is s = .getRunText()
			return
		try
			s = QueryKeys(s).Map!({ it.Replace(',', ', ') }).Join('\r\n')
		catch (e)
			s = e
		.switchToText("Keys:\n" $ s)
		}

	On_Schema()
		{
		s = .edit.GetCurrentWord()
		if s !~ '\s'
			{
			SchemaView.Goto(s)
			return
			}

		try
			s $= "\n\n" $ QueryColumns(s).Sort!().Join(', ').
				WrapLines(50 /*= width*/).Join('\r\n')
		catch (e)
			s $= "\n\n" $ e
		ScintillaIDEControl(set: s)
		}
	On_Count()
		{
		if "" is s = .getRunText()
			return
		.run('(' $ QueryStripSort(s) $ ') summarize count')
		}
	On_Hash()
		{
		if "" is s = .getRunText()
			return
		a = QueryHash(s, details:)
		b = QueryAltHash(s, details:)
		.switchToText(
			(a is b ? "SAME\n" : "DIFFERENT\n") $
			"   QueryHash: " $ a $
			"QueryAltHash: " $ b)
		}
	On_Format(show = false)
		{
		query = GetRunText(.edit)
		try
			query = Global(#FormatQuery)(query)
		catch (e)
			query = e
		if show or KeyPressed?(VK.SHIFT)
			.switchToText(query)
		else
			.edit.Paste(query)
		}
	On_Show_Format()
		{
		.On_Format(show:)
		}

	GetState()
		{
		return Object(text: .edit.Get())
		}
	SetState(statedata)
		{
		if statedata.Member?('text')
			{
			.edit.Set(statedata.text)
			SendMessage(.edit.Hwnd, EM.SETSEL, 0, -1)	// Select all the text
			}
		}

	AddQuery(q)
		{
		if not String?(q)
			return
		oldValue = Opt(.edit.Get(), '\r\n\r\n')
		.edit.Set(oldValue $ q)
		.edit.SetSelect(oldValue.Size(), q.Size())
		}

	VirtualList_DoubleClick(rec, col /*unused*/)
		{
		.On_Context_Inspect_Record(rec)
		return true
		}

	On_Context_Copy(rec, col)
		{
		if false isnt value = .copyValue(col, rec)
			ClipboardWriteString(value)
		}
	copyValue(col, rec)
		{
		if col is false
			return false
		value = rec is false ? col : rec[col]
		if not String?(value)
			value = Display(value)
		if value.Capitalized?()
			value = value[0].Upper() $ value[1..]
		return value
		}
	On_Context_Copy_Record(rec, col /*unused*/)
		{
		ClipboardWriteString(String(.cleanupRec(rec.Copy())))
		}
	cleanupRec(rec)
		{
		rec.Delete('vl_full_display').Copy()
		}
	On_Context_Copy_FieldValue(rec, col)
		{
		if false isnt value = .copyValue(col, rec)
			ClipboardWriteString(col $ ': ' $ value)
		}
	On_Context_Go_To_Definition(rec /*unused*/, col)
		{
		GotoLibView(col)
		}
	On_Context_Show_References(rec /*unused*/, col)
		{
		FindReferencesControl(col)
		}
	On_Context_Inspect_Value(rec, col)
		{
		if rec isnt false or col is false
			Inspect.Window(rec[col], 'Inspect Value')
		}
	On_Context_Inspect_Record(rec)
		{
		if rec isnt false
			Inspect.Window(.cleanupRec(rec), 'Inspect Record')
		}
	On_Context_Export()
		{
		GlobalExportControl(.query, includeInternal:)
		}
	On_Context_Print_Record(rec)
		{
		if rec isnt false
			CurrentPrint(rec, .Window.Hwnd, .query, TruncateKey(.query))
		}
	}
