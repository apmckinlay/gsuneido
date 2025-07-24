// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* Testing:
	- tab from double-click should not be reused
	- both double-click and next/prev should "reuse" tab from next/prev
	- next/prev to same record should not close tab and reopen it
*/
Controller
	{
	Title: "Find in Libraries"
	uid: false
	New()
		{
		super(.layout())
		.uid = Timestamp()
		.data = .FindControl('Data')
		.output = .FindControl('output')
		.output.SetWordChars(
			"_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!:")
		.recent = Object()
		.buttons = #(Reset, Clear__Find, Find, Next, Prev).Map(.FindControl)
		.stop = .FindControl(#Stop)
		.stop.SetEnabled(false)
		.stop.SetTextColor(CLR.Inactive)
		// cannot use PersistentWindow since the control is not ready
		.data.SetField('sort', UserSettings.Get('FindInLibrariesControl_sort', 'name'))
		.resetLibsAndExclude()
		.printContent = Object()
		SetTimer(.WindowHwnd(), 0, 100/*=inverval*/, .updateOutput)
		}

	ResetTheme()
		{
		.output.SetWordChars(
			"_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!:")
		}

	layout()
		{
		return Object('Vert'
			Object('Record'
				Object('Horz'
					Object('Vert',
						.controlLayout(),
						#(Skip 4),
						Object('Pair'
							#(Static 'Name'), .findRepeatControl('name')),
						#(Skip 4),
						Object('Pair'
							#(Static 'Text'), .findRepeatControl('text')),
						#(Skip 4),
						.findByTokenOption()
						#(Skip 4),
						.findByExpressionOption()
						#(Skip 4),
						Object('Horz', .showOptions(),
							#(Skip 80), #(Button Next), #Skip, #(Button Prev))
						),
					#(Skip 20),
					.buttons()
					)
				),
			#(Skip 6)
			#(WorkSpaceOutput readonly:, edge: false, name: 'output')
			)
		}
	controlLayout()
		{
		return #(Pair
			(Static 'Library')
			(Horz
				(ChooseLibraries, all:, name: "libs"),
				#Skip #(CheckBox, "Exclude"	tabover:, name: 'exclude')
				Skip
				(EnhancedButton 'Reset', buttonStyle:, mouseEffect:, pad: 30),
				Skip
				(EnhancedButton 'Clear && Find', buttonStyle:, mouseEffect:, pad: 20),
				Skip
				(EnhancedButton 'Find', defaultButton:, buttonStyle:,
					mouseEffect:, pad: 30),
				Skip
				(EnhancedButton 'Stop', buttonStyle:, mouseEffect:, pad: 30)
				)
			)
		}
	findByTokenOption()
		{
		return #(Pair
			(Static 'Token'),
			(Horz
				(FieldHistory, font: '@mono', size: '+1', width: 60,
					cue: 'must be valid tokens e.g. no mismatched quotes'
					name: 'bytoken')
				)
			)
		}
	findByExpressionOption()
		{
		return Object(#Pair
			#(Static 'Expression'),
			Object(#Horz
				#(FieldHistory, font: '@mono', size: '+1', width: 60,
					cue: 'must be a valid expression, single letters match expressions'
					name: 'byexpression')
				#Skip
				Object('EnhancedButton', command: 'Find By Expression Help',
					image: 'questionMark_black.emf', imageColor: CLR.Inactive,
					mouseOverImageColor: CLR.Highlight, mouseEffect:, imagePadding: .1)))
		}
	showOptions()
		{
		return #(Horz
			#(Pair
				(Static 'Show'),
				(RadioButtons, "lines", "records", "random record",
					horz:, name: "show"))
			#Skip, #Skip, #Skip
			#(Pair
				(Static 'Sort By'),
				(RadioButtons, "name", "library, name",
					horz:, name: "sort"))
			)
		}
	buttons()
		{
		#(Vert
			(Presets 'Find/Replace', 'Find/Replace' xstretch: 0),
			(Skip 4),
			(MenuButton Recent xstretch: 0)
			)
		}
	DefaultButton: "Find"

	findRepeatControl(name)
		{
		return Object('Repeat', 'WorkSpaceFind', name: name $ 'Repeat')
		}

	SetDefaultFocus()
		{
		if false isnt repeat = .FindControl('textRepeat')
			if false isnt ctrl = repeat.FindControl('text')
				ctrl.SetFocus()
		}

	FieldReturn()
		{
		if not .runningInThread
			.On_Find()
		}

	getData()
		{
		data = .data.Get()
		if data.GetDefault(#nameRepeat, "") is ""
			data.nameRepeat = Object()
		if data.GetDefault(#textRepeat, "") is ""
			data.textRepeat = Object()
		data.findUid = .uid
		return data
		}

	On_Clear__Find()
		{
		.ClearResults()
		.On_Find()
		}

	On_Find()
		{
		data = .getData()
		if '' isnt msg = .validate(data)
			{
			InfoWindowControl(msg, titleSize: 0)
			return
			}
		.save()
		.printHeader()
		.enableUI(false)
		.findInThread(data)
		}

	runningInThread: false
	findInThread(data)
		{
		Thread()
			{
			Thread.Name('FindInLibraries-thread')
			.runningInThread = true
			.find(data)
			.Defer({ .enableUI(true) }, uniqueID: #enableUI)
			.runningInThread = false
			}
		}

	enableUI(enable?)
		{
		.data.SetReadOnly(not enable?)
		.buttons.Each({ it.SetEnabled(enable?) })
		.stop.SetEnabled(not enable?)
		.stop.SetTextColor(not enable? ? CLR.Active : CLR.Inactive)
		}

	find(data)
		{
		.updateWorkSpaceFindStop(false)
		try
			FindInLibraries(data, .print, .printFooter)
		catch (e)
			{
			if e.Prefix?('regex:')
				info = "Invalid regex: " $ e.RemovePrefix('regex: ')
			else
				info = "Invalid expression: " $ e

			.showInfoWindow(info)
			.printFooter(0, 0)
			.print(info)
			}
		}
	// wrapper so we only need to know about .uid here
	updateWorkSpaceFindStop(value)
		{
		ServerSuneido.Add('workSpaceFindStop', value, .uid)
		}

	showInfoWindow(info)
		{
		if .runningInThread is false
			{
			InfoWindowControl(info, titleSize: 0)
			return
			}
		.Defer({ InfoWindowControl(info, titleSize: 0) }, uniqueID: 'infoWindow')
		}

	validate(data)
		{
		if data.exclude is true and data.libs is '(All)'
			return 'Please choose some libraries to find'
		return data.textRepeat.Any?({ it.text isnt '' }) or
			data.nameRepeat.Any?({ it.text isnt '' }) or
			data.bytoken isnt '' or
			data.byexpression isnt ''
			? "" : 'Please enter something to find'
		}

	printHeader()
		{
		if .output.Get() isnt ''
			.print('='.Repeat(80/*=headDecorateLength*/))
		.print(.header())
		}
	header()
		{
		data = .getData()

		h0 = data.show is 'random record' ? data.show $ ' from ' : ''

		inText = data.textRepeat
		h1 = Opt(.headerFor(inText), ' in text')

		inName = data.nameRepeat
		h2 = Opt(.headerFor(inName), ' in name')

		h3 = Opt(data.bytoken, ' by token')

		h4 = Opt(data.byexpression, ' by expression')

		return 'Find ' $ h0 $ Join(' with ', h1, h2, h3, h4) $
			(data.exclude is true ? ' excluding ' : ' in ') $ data.libs
		}
	headerFor(repeat)
		{
		return repeat.
			Filter({ not it.text.Blank?() }).
			Map({ (it.exclude is true ? 'not ' : '') $ Display(it.text) }).
			Join(' and ')
		}

	printFooter(recs, lines)
		{
		.print(recs is ''
			? "No matches"
			: "Found " $
				(lines is 0 ? '' : (lines $ ' lines in ')) $
				recs $ " records ")
		}

	print(s)
		{
		.printContent.Add(s $ '\n')

		if .runningInThread is false
			{
			.updateOutput()
			return
			}
		}

	updateOutput(@unused)
		{
		content = ''
		while .printContent isnt s = .printContent.PopFirst()
			content $= s
		if content > ''
			{
			.output.AppendText(content)
			.output.Update()
			}
		}

	On_Next()
		{
		pos = .output.GetSelectionStart()
		line = .output.LineFromPosition(pos)
		if false is line = .findNext(line, 1)
			return
		.output.GotoLine(line)
		.goto(recycleTab:)
		}

	findNext(line, dir)
		{
		cur = line
		if 1 >= nlines = .output.GetLineCount()
			return false

		do
			{
			line = (line + dir + nlines) % nlines
			s = .output.GetLine(line)
			if s =~ '^[[:alpha:]]+:[[:alpha:]][_[:alnum:]]*[?!]?(:[[:digit:]]+:|\s)'
				return line
			}
		while line isnt cur

		return false
		}

	On_Prev()
		{
		pos = .output.GetSelectionStart()
		line = .output.LineFromPosition(pos)
		if false is line = .findNext(line, -1)
			return
		.output.GotoLine(line)
		.goto(recycleTab:)
		}

	ClearResults() // called by LibViewControl
		{
		.output.Set("")
		}

	On_Reset()
		{
		.FindControl('Presets').On_Presets('New')
		.resetLibsAndExclude()
		}

	resetLibsAndExclude()
		{
		if TableExists?('Contrib')
			{
			.data.SetField('libs', 'Contrib')
			.data.SetField('exclude', true)
			}
		else
			{
			.data.SetField('libs', '(All)')
			.data.SetField('exclude', false)
			}
		}

	On_Stop()
		{
		.updateWorkSpaceFindStop(true)
		}

	max_recent: 10
	save()
		{
		heading = .header().Replace('^Find ').Trim()
		.recent.RemoveIf({ it.heading is heading })
		.recent.Add([:heading, data: .getData().Copy()] at: 0)
		.recent.Delete(.max_recent + 1)
		}

	MenuButton_Recent()
		{
		return .recent.Map({ it.heading })
		}

	On_Recent(heading)
		{
		if false isnt i = .recent.FindIf({ it.heading is heading })
			.data.Set(.recent[i].data)
		else
			Beep()
		}

	Scintilla_DoubleClick(source /*unused*/)
		{
		.goto()
		}

	goto(recycleTab = false)
		{
		line = .output.GetLine()
		if line !~ '^.+?:[[:alpha:]][_[:alnum:]]*[?!]?\>'
			return false
		.output.Home()
		.output.CharRight()
		if recycleTab
			.recycleTab(line.BeforeFirst(':'), line.AfterFirst(':').BeforeFirst(':'))
		.output.Recv(#On_Go_To_Definition)
		return true
		}

	nextPrevTab: false
	recycleTab(lib, name)
		{
		if .nextPrevTab is false
			.nextPrevTab = [:lib, :name]
		if .nextPrevTab?(lib, name)
			return
		libView = .libView()
		if libView.Editor is false // No tabs are opened, no tabs to close/check
			return
		curLib = libView.CurrentTable()
		curName = libView.CurrentName()
		if .nextPrevTab?(curLib, curName) and not .modifiedRec?(curLib, curName)
			libView.Explorer.CloseActiveTab()
		.nextPrevTab = [:lib, :name]
		}

	nextPrevTab?(lib, name)
		{
		return lib is .nextPrevTab.lib and name is .nextPrevTab.name
		}

	libView()
		{
		return GotoPersistentWindow('LibViewControl', LibViewControl)
		}

	modifiedRec?(lib, name)
		{
		rec = Query1(lib, :name, group: -1)
		return rec is false ? false : rec.lib_modified isnt ''
		}

	On_Find_By_Expression_Help()
		{
		OpenBook('suneidoc', path: '/Getting Started/WorkSpace#FindByExpression')
		}

	Destroy()
		{
		UserSettings.Put('FindInLibrariesControl_sort', .data.GetField('sort'))
		super.Destroy()
		}
	}
