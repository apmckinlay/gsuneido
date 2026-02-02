// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "TestRunner"

	CallClass()
		{
		GotoPersistentWindow('TestRunnerGui', TestRunnerGui)
		}

	New(lib = "")
		{
		super(.controls())
		.libs = .FindControl('ChooseList')
		.list = .FindControl('List')
		.list.SetMultiSelect(true)
		.ovbar = .FindControl('ovbar')
		.ovbar.SetTopMargin(GetSystemMetrics(SM.CXHSCROLL))
		// assuming header control is same height as scroll bar
		.statusbar = .Vert.Status
		.stop = .FindControl('stop_on_error')
		.stop.Set(true)
		.check_tables = .FindControl('check_tables')
		.runOnServer = .FindControl('run_on_server')
		.debugOptions = .FindControl('debugOptions')
		.snapshotTables = .FindControl('snapshotTables')
		.snapshotTables.SetReadOnly(true)
		.errors = .FindControl('errors')
		.errors.SetReadOnly(true)
		.time = .FindControl('total_time')
		.space = .FindControl('total_space')
		.memory = .FindControl('total_memory')
		.nwarnings = .FindControl('total_nwarnings')
		if lib isnt ""
			{
			.libs.Set(lib)
			.set_lib(lib)
			}
		.successColor = CLR.GREEN
		for c in GetContributions('TestRunnerButtons')
			.Redir('On_' $ ToIdentifier(c[0]), c[1])
		}
	columns: (testrunnergui_result, testrunnergui_lib,
		testrunnergui_name, testrunnergui_time, testrunnergui_nwarnings,
		testrunnergui_dbspace, testrunnergui_arena, testrunnergui_errors)
	controls()
		{
		libs = Libraries().Add("(All)", at: 0)
		return Object('Vert',
			Object('Horz', 'Skip', .top(libs)),
			Object('Horz'
				Object('ListStretch' columns: .columns,
					stretchColumn: 'testrunnergui_errors',
					columnsSaveName: 'TestRunnerGui')
				Object('OverviewBar', name: 'ovbar', priorityColor: CLR.ErrorColor)),
			Object('Horz'
				#(Skip 1),
				#(Editor height: 2, ystretch: 0, name: errors),
					// should be readonly but then you can't make it red
				['Vert',
					#(Skip 6)
					#(Horz
						Skip
						(Static "Totals:", weight: "bold")
						Skip
						(Pair (Static "Time")
							(Field width: 8, readonly:, name: total_time))
						Skip
						(Pair (Static "Warnings")
							(Field width: 5, readonly:, name: total_nwarnings))
						Skip
						(Pair (Static "DB Space")
							(Field width: 8, readonly:, name: total_space))
						Skip
						(Pair (Static "Memory")
							(Field width: 8, readonly:, name: total_memory))
						Skip),
					#(Skip 6),
					['Horz', #(Skip 4), .bottom_buttons()],
					#(Skip 6)
					]),
			"Statusbar"
			)
		}
	top(libs)
		{
		row1 = Object('Horz', Object('ChooseList', libs), 'Skip')
		row1.Append(.top_buttons())
		row1.Add(#(Button, 'Run All'), 'Skip',
			#(Button, 'Continue', tip: 'Run from selected'), 'Skip',
			#(Button, 'Run Failed'), 'Skip',
			#(CheckBox, 'Stop on Error', name: stop_on_error), 'Skip',
			#(CheckBox, 'Check per Test', name: check_tables,
				tip: 'Check file and table differences after each test'), 'Skip',
			#(CheckBox, 'Run Tests on Server', name: run_on_server), 'Skip')
		row2 = Object('Horz',
			#(Static, 'Debugging Options '),
			Object('ChooseManyAsObject', idField: 'Option', displayField: 'Option',
				cols: #(Option, Description), list: .debugOptionsList(), height: 1,
				delimiter: ',\r\n', saveColName: 'TestRunner_Debugging_Options',
				name: 'debugOptions'), 'Skip',
			#(Static, 'Snapshot Tables '),
			Object('ChooseMany', list: SnapshotTables.Candidates(),
				name: 'snapshotTables'),
			#(Button, 'Review', tip: 'Review Snapshot Table differences'))
		return Object('Flow', Object(row1, row2), skip: 0)
		}
	top_buttons()
		{
		ob = Object()
		for c in GetContributions('TestRunnerButtons')
			if c.GetDefault('front', false) is true
				ob.Add(Object('Button', c[0])).Add('Skip')
		return ob
		}
	bottom_buttons()
		{
		ob = Object('Horz', #(Fill fill: .001))
		for c in GetContributions('TestRunnerButtons')
			if c.GetDefault('front', false) isnt true
				ob.Add(Object('Button', c[0])).Add('Skip')
		ob.Add(#(Fill fill: .001))
		return ob
		}

	debugOptionsList()
		{
		return #(
			(Option: 'Snapshots: Output',
				Description: 'Will snapshot all the tables in: Snapshot Tables')
			(Option: 'Snapshots: Compare',
				Description: 'Will compare the specified snapshot tables against ' $
					'their live counterparts')
			(Option: 'Snapshots: Persist',
				Description: 'Will prevent Snapshots > Compare from dropping the ' $
					'snapshots post comparison')
			)
		}

	Commands: (
		(Run_Selected,		"F9")
		(Debug_Selected,	"Ctrl+K")
		(Go_To_Definition,	"F12")
		(Close,				"",			"Close this window")
		(Refresh, 			"F5")
		)
	Menu:
		(
		('&File',
			'&Refresh', '',
			'&Close')
		)
	NewValue(value, source)
		{
		if source is .debugOptions
			{
			snapsUsed? = value.Has?('Snapshots: Compare') or
				value.Has?('Snapshots: Output')
			.snapshotTables.SetReadOnly(not snapsUsed?)
			}
		if (source isnt .libs or not source.Valid?())
			return
		.set_lib(value)
		.libs.Dirty?(false)
		}
	On_Refresh()
		{
		.set_lib(.libs.Get())
		}
	set_lib(lib)
		{
		.reset_statusbar()
		if (lib is "")
			{
			.list.Set(Object())
			return
			}
		tests = Object()
		block = {|name| tests.Add(Object(testrunnergui_lib: lib,
			testrunnergui_name: name)) }
		if (lib is '(All)')
			for (lib in Libraries())
				TestRunner.ForeachTest(lib, block)
		else
			TestRunner.ForeachTest(lib, block)
		.list.Set(tests)
		.reset(clearMarks?:)
		}
	reset(clearMarks? = false)
		{
		.errors.Set("")
		.errors.SetReadOnly(true)
		.time.Set("")
		.nwarnings.Set("")
		.space.Set("")
		.memory.Set("")
		.resize_ovbar()
		if clearMarks?
			.ovbar.ClearMarks()
		.reset_statusbar()
		.Window.Update()
		}
	resize_ovbar()
		{
		.ovbar.SetNumRows(.list.GetNumRows())
		.ovbar.SetMaxRowHeight(.list, #GetRowHeight)
		}
	reset_statusbar()
		{
		.statusbar.Set("")
		.statusbar.SetBkColor(false)
		}

	runAll: false
	clearDate: false
	On_Run_All()
		{
		if not .runTests?(libs = .libs.Get())
			return
		// stdlib transaction test will fail if outstanding transactions
		if #('' '(All)' 'stdlib').Has?(libs) and not Database.Transactions().Empty?()
			{
			Alert("There are outstanding transactions.\n" $
				"Make sure all transactions are closed before running all tests.",
				"Run All", .Window.Hwnd, MB.ICONERROR)
			return
			}

		if libs is ""
			.libs.Set('(All)')
		.On_Refresh() // refresh list to handle switching libraries, new or removed tests
		.runAll = true
		.clearDate = Date()
		.run_testlist(.list.Get().Members())
		}

	runTests?(libs)
		{
		libs = libs is '' or libs is '(All)' ? Libraries() : Object(libs)
		return .checkLibraries(libs)
		}

	checkLibraries(libs)
		{
		invalids = Object()
		for lib in libs
			invalids.Merge(CodeState.InvalidRecs(lib))
		return .continueRun?(invalids)
		}

	continueRun?(invalids)
		{
		invalidStr = .flaggedStr('Found records with syntax errors', invalids)
		if invalidStr is ''
			return true
		sep = '\r\n\r\n'
		return YesNo(
			Opt(invalidStr, 'Running tests now may result in a false positive.', sep) $
			'Continue?',
			.Title, .Window.Hwnd, MB.ICONWARNING)
		}

	flaggedStr(prefix, flagged)
		{
		return flagged.NotEmpty?()
			? prefix $ ':\r\n\t- ' $ flagged.Sort!().Join('\r\n\t- ') $ '\r\n'
			: ''
		}

	On_Review()
		{
		title = 'Snapshot Tables Review'
		differences = SnapshotTables.Review()
		if differences.Empty?()
			.AlertInfo(title, 'No snapshot differences to review')
		else
			Inspect.Window(differences, :title, hwnd: .Window.Hwnd)
		}

	On_Continue()
		{
		selected = .list.GetSelection()
		if (selected is #())
			return
		runlist = Object()
		first = selected.Sort!()[0]
		while first < .list.GetNumRows() and
			.list.GetRow(first).testrunnergui_result is 'X'
			first++
		for (i = first; i < .list.GetNumRows(); i++)
			runlist.Add(i)
		.run_testlist(runlist)
		}
	On_Run_Failed()
		{
		list = Object()
		data = .list.Get()
		for i in data.Members()
			if data[i].testrunnergui_errors isnt "" and
				not data[i].testrunnergui_errors.Prefix?("WARNING")
				list.Add(i)
		if list isnt #()
			.run_testlist(list)
		}
	On_Run_Selected()
		{
		list = .list.GetSelection()
		if (list is #())
			return
		list.Sort!()
		.run_testlist(list)
		}
	run_testlist(testlist)
		{
		if (.libs.Get() is '')
			return
		origValidateQueryAny1? = .getSuneidoVariable('ValidateQueryAny1?')
		.reset()
		.beforeRunList()
		for i in testlist
			.reset_row(i)
		.list.Repaint()
		// prevent users from trying to "modify" the data
		// causes refresh issues.
		// also stop them from running the tests multiple times at once
		this.SetEnabled(false)
		data = .list.Get()
		.testRunnerServerOrClient(testlist, data, origValidateQueryAny1?)
		.setSuneidoVariable('ValidateQueryAny1?', origValidateQueryAny1?)
		}
	setSuneidoVariable(member, val)
		{
		if .runOnServer.Get() is true
			ServerSuneido.Set(member, val)
		else
			Suneido[member] = val
		}
	getSuneidoVariable(member)
		{
		return .runOnServer.Get() is true
			? ServerSuneido.Get(member, false)
			: Suneido.GetDefault(member, false)
		}
	prevConditions: false
	beforeRunList()
		{
		.prevConditions = Object()
		Plugins().ForeachContribution('TestRunner', 'beforeRun')
			{
			.prevConditions.MergeNew((it.func)())
			}
		}
	reset_row(i)
		{
		.list.ClearHighlight(i)
		.ovbar.RemoveMark(i)
		test = .list.Get()[i]
		test.testrunnergui_result = ""
		test.testrunnergui_time = ""
		test.testrunnergui_nwarnings = ""
		test.testrunnergui_dbspace = ""
		test.testrunnergui_arena = ""
		test.testrunnergui_errors = ""
		}

	testRunnerServerOrClient(testlist, data, origValidateQueryAny1?)
		{
		observer = .runOnServer.Get() is true
			? TestObserverOnServer() : TestObserverGui()
		runTestRunner = .runOnServer.Get() is true
			? TestRunner.Run1OnServer : TestRunner.Run1
		TestRunner.RunTests(
			{ .runTests(testlist, data, observer, runTestRunner,origValidateQueryAny1?) })
		}

	runTests(testlist, data, observer, runTestRunner, origValidateQueryAny1?)
		{
		// used to track ui updates in the thread.
		Suneido.testUi = Object(ovbarMarks: Object(),
			repaintRows: Object(),
			statusBar: Object(color: "", msg: ""),
			state: 'running',
			:observer,
			totalTime: 0,
			:origValidateQueryAny1?)
		debuggingOptions = .processDebuggingOptions()
		Thread(name: 'test runner')
			{
			if debuggingOptions.snapshot?
				.snapshotCompare(debuggingOptions)
			_systemChanges_excludeTables = SnapshotTables.Snaps()
			.timeTestEvent('Collecting system state')
				{
				before = SystemChanges.GetState()
				}
			guiTests = Object()
			Suneido.testUi.totalTime += Timer()
				{
				.forEachTest(testlist)
					{ |i|
					.handleTestFn(data, i, observer, runTestRunner, guiTests)
					}
				}
			if debuggingOptions.compareSnaps?
				.snapshotCompare(debuggingOptions, compare?:)
			Defer({ .guiTests(:data, :guiTests, :observer, :runTestRunner, :before) })
			}
		.Delay(100, .updateUi, uniqueID: 'updateUi') /*= 1/10 sec */
		return
		}

	processDebuggingOptions()
		{
		debugOptions = .debugOptions.Get()
		snapshot? = debugOptions.Has?('Snapshots: Output')
		compareSnaps? = debugOptions.Has?('Snapshots: Compare')
		persistSnaps? = debugOptions.Has?('Snapshots: Persist')
		snapshotTables = compareSnaps? or snapshot?
			? .snapshotTables.Get().Split(',')
			: #()
		return [:persistSnaps?, :snapshot?, :compareSnaps?, :snapshotTables]
		}

	snapshotCompare(debuggingOptions, compare? = false)
		{
		if not compare?
			.timeTestEvent('Snapshotting tables')
				{
				SnapshotTables.Ensure()
				SnapshotTables(debuggingOptions.snapshotTables)
				}
		else
			.timeTestEvent('Comparing snapshots against live tables')
				{
				tableDifferences = SnapshotTables.Compare(debuggingOptions.snapshotTables)
				SnapshotTables.Log(tableDifferences, debuggingOptions.persistSnaps?)
				}
		}

	timeTestEvent(msg, block)
		{
		Defer({ .displayTestEventInfo(msg, '\r\n') })
		t = Timer(block)
		Defer({ .displayTestEventInfo('duration: ' $ t $ ' (secs)', ', ') })
		}

	displayTestEventInfo(msg, separator = '')
		{
		if not .Destroyed?()
			.errors.Set(Opt(.errors.Get(), separator) $ msg)
		}

	forEachTest(testlist, block)
		{
		for i in testlist
			if false is block(i)
				break
		}

	handleTestFn(data, i, observer, runTestRunner, guiTests)
		{
		lib = data[i].testrunnergui_lib
		name = data[i].testrunnergui_name
		x = Query1(lib, :name, group: -1)
		// need to handle any tests that are not thread safe
		// make sure they run on the main thread
		if CodeTags.ExtractTags(x.text).Has?('win32')
			{
			guiTests.Add(i)
			return true
			}
		try
			result = .runSingleTest(data, i, observer, runTestRunner)
		catch (err)
			{
			.handleTestRunnerError(data, i, observer, err)
			result = false
			}
		return result
		}

	handleTestRunnerError(data, i, observer, err)
		{
		observer.BeforeTest('Test Error')
		error = 'ERROR: ' $ err
		observer.Error('Test Error: ', error)
		data[i].testrunnergui_errors = error
		Suneido.testUi.repaintRows.Add(Object(i, true))
		observer.AfterTest('Test Error', 0, 0, 0, debug:)
		Suneido.testUi.state = 'failed'
		}

	runSingleTest(data, i, observer, runTestRunner)
		{
		// this is needed for tests running on the main thread
		// to ensure that the test thread waits untill this test finishes before
		// continuing on
		.setSuneidoVariable('ValidateQueryAny1?', true)
		result = .runtest(data[i], i, observer, runTestRunner)
		if result isnt ""
			{
			Suneido.testUi.ovbarMarks.Add(Object(i, CLR.ErrorColor))
			if observer.Totals.n_failures is 1
				Suneido.testUi.statusBar.color = CLR.ErrorColor
			if .stop.Get() is true
				{
				Suneido.testUi.state = 'failed'
				return false
				}
			}
		else
			{
			Suneido.testUi.ovbarMarks.Add(Object(i, .successColor))
			}
		return true
		}

	guiTests(data, guiTests, observer, runTestRunner, before)
		{
		if .Destroyed?()
			return
		Suneido.testUi.totalTime +=
			.runGuiTests(data, guiTests, observer, runTestRunner)
		.finishTests(observer, before)
		}

	timerDelay: 300 // Adjust this to change how "smooth" the gui updates
	updateUi()
		{
		// runs on the main thread
		// this is the part that actually updates the ui.
		// don't use for each. need to prevent object modified durring interation
		// they way this works the UI updates should be continuoisly playing "catch-up"
		// with the tests - i.e. the tests can be adding to these objects while we are
		// still looping through them. once we catch up, then restart the timer
		while Suneido.testUi.ovbarMarks.Size() > 0 or
			Suneido.testUi.repaintRows.Size() > 0
			{
			if false isnt item = Suneido.testUi.ovbarMarks.Extract(0, false)
				.updateOvBar(@item)
			if false isnt item = Suneido.testUi.repaintRows.Extract(0, false)
				.updateList(@item)
			}

		// Restart the timer if the tests haven't finished/failed yet
		if Suneido.testUi.state is 'running'
			{
			.Delay(.timerDelay, .updateUi, uniqueID: 'updateUi')
			return
			}

		// need to do this here instead of in .finishTests otherwise thie status bar
		// will update BEFORE the results of the thread safe tests have been displayed
		.displayTestResults()

		// ensure that the bottom of the list is selected after the non thread safe tests
		// have run (only if running All Tests
		if Suneido.testUi.state is 'finished' and .runAll is true
				.list.SetSelection(.list.Get().Size()-1)

		// reset runAll flag if tests finished, either succesfully, or from a failure
		.runAll = false
		}

	displayTestResults()
		{
		observer = Suneido.testUi.observer
		color = observer.Totals.n_failures is 0 ? .successColor : CLR.ErrorColor
		msg = observer.Totals.n_failures is 0
			? " S U C C E S S - " $ observer.Totals.n_tests $ " tests"
			: " F A I L U R E - " $ observer.Totals.n_failures $ " failure(s) out of " $
						observer.Totals.n_tests $ " tests"
		.updateStatusBar(color, msg)
		.time.Set(Suneido.testUi.totalTime.RoundToPrecision(2) $ " sec")
		.space.Set(ReadableSize(observer.Totals.dbgrowth))
		.memory.Set(ReadableSize(observer.Totals.memory))
		.nwarnings.Set(observer.Totals.nwarnings)
		}

	// NOTE: updateUi will still run one more time after we run this.
	// any thing in this function that updates the ui will happen BEFORE the results
	// from the thread safe tests gets displayed
	finishTests(observer, before)
		{
		// signal that we have finished running the tests
		// check if the tests haven't failed first
		if Suneido.testUi.state isnt 'failed'
			Suneido.testUi.state = 'finished'

		.checkSystemChanges(observer, before)

		// re-enable after the tests are finished
		this.SetEnabled(true)
		ResetCaches()
		.afterRunList()

		// only output the "Success" for SVC if user selected RunAll
		if .runAll is true
			{
			if observer.Totals.n_failures is 0 and .libs.Get() is '(All)'
				TestRunner.OutputSuccess(.clearDate)
			}
		.setSuneidoVariable('ValidateQueryAny1?', Suneido.testUi.origValidateQueryAny1?)
		}

	// this runs on the main thread. Runs any tests that are not thread safe (tag: win32)
	runGuiTests(data, guiTests, observer, runTestRunner)
		{
		if observer.Totals.n_failures > 0 and .stop.Get() is true
			return 0
		Timer()
			{
			.forEachTest(guiTests)
				{ |i|
				try
					result = .runSingleTest(data, i, observer, runTestRunner)
				catch (err)
					{
					.handleTestRunnerError(data, i, observer, err)
					result =  false
					}
				result
				}
			}
		}

	updateOvBar(i, markColor)
		{
		.ovbar.AddMark(i, markColor)
		.ovbar.Update()
		}

	updateStatusBar(color, message)
		{
		.statusbar.SetBkColor(color)
		if message isnt ""
			.statusbar.Set(message)
		.statusbar.Update()
		}

	updateList(i, error = false)
		{
		.list.SetSelection(i)
		.list.Update()
		if error
			.list.AddHighlight(i, CLR.ErrorColor)
		.list.RepaintRow(i)
		.list.Update()
		}

	runtest(test, i, observer, runTestRunner)
		{
		runTestRunner(test.testrunnergui_name,
			:observer, check_tables: .check_tables.Get(), lib: test.testrunnergui_lib)
		x = observer.Data
		test.testrunnergui_errors = x.errors $ x.warnings
		test.testrunnergui_result = x.errors is '' ? '' : 'X'
		test.testrunnergui_time = x.time
		test.testrunnergui_nwarnings = x.nwarnings
		test.testrunnergui_dbspace = x.dbgrowth
		test.testrunnergui_arena = x.memory

		Suneido.testUi.repaintRows.Add(Object(i, x.errors isnt ""))
		return x.errors
		}
	checkSystemChanges(observer, before)
		{
		_systemChanges_excludeTables = SnapshotTables.Snaps()
		if "" is s = SystemChanges.CompareState(before)
			return
		.errors.Set(s)
		.errors.SetReadOnly(false) // so you can change the color
		.errors.SetBgndColor(CLR.ErrorColor)
		.statusbar.SetBkColor(CLR.ErrorColor)
		observer.BeforeTest("System Changes")
		observer.Error("", s)
		observer.AfterTest("System Changes", 0, 0, 0)
		}
	afterRunList()
		{
		Plugins().ForeachContribution('TestRunner', 'afterRun')
			{
			(it.func)(.prevConditions)
			}
		.prevConditions = false
		}
	On_Debug_Selected()
		{
		orig = .getSuneidoVariable('ValidateQueryAny1?')
		.setSuneidoVariable('ValidateQueryAny1?', true)
		list = .list.GetSelection()
		if list.Size() isnt 1
			{
			.AlertInfo("Debug Selected", "Please select a single test to debug")
			return
			}
		i = list[0]
		test = .list.Get()[i]
		Suneido.TestRunner = true
		ServerSuneido.Set(#TestRunner, true)
		Construct(test.testrunnergui_name).Debug()
		Suneido.Delete(#TestRunner)
		ServerSuneido.DeleteMember(#TestRunner)
		.reset_row(i)
		.reset_statusbar()
		.ovbar.AddMark(i, CLR.WHITE)
		.setSuneidoVariable('ValidateQueryAny1?', orig)
		}
	List_AllowCellEdit(col, row /*unused*/)
		{
		return .list.GetCol(col) is 'testrunnergui_errors'
		}
	List_EditFieldReadonly(@unused)
		{
		return true
		}
	List_WantNewRow()
		{
		return false
		}
	contextMenu: ("Run\tF9", "Profile", "Debug\tCtrl+K", "Go To Definition\tF12",
		"Copy Name", "Shuffle")
	List_ContextMenu(x, y)
		{
		if (.list.Get().Empty?())
			return
		ContextMenu(.contextMenu).ShowCall(this, x, y)
		}
	List_DoubleClick(row, col)
		{
		if row is false or .list.GetCol(col) isnt 'testrunnergui_name'
			return 0
		x = .list.GetRow(row)
		GotoLibView(x.testrunnergui_name)
		return false // don't edit
		}
	On_Context_Run()
		{
		.On_Run_Selected()
		}
	On_Context_Profile()
		{
		if .libs.Get() is ''
			return

		rec = .list.GetCurrentRecord()
		if rec is false
			{
			.AlertError("Profile", "Please select one test to profile")
			return
			}

		this.SetEnabled(false)
		RunWithProfile()
			{ TestRunner.Run1(rec.testrunnergui_name, observer: TestObserver) }
		this.SetEnabled(true)
		}
	On_Context_Debug()
		{
		.On_Debug_Selected()
		}
	On_Context_Go_To_Definition()
		{
		.On_Go_To_Definition()
		}
	On_Go_To_Definition()
		{
		list = .list.GetSelection()
		if list.Size() isnt 1
			return
		data = .list.Get()
		GotoLibView(data[list[0]].testrunnergui_name)
		}
	On_Context_Copy_Name()
		{
		list = .list.GetSelection()
		i = list[0]
		test = .list.Get()[i]
		ClipboardWriteString(test.testrunnergui_name)
		}
	On_Context_Shuffle()
		{
		.list.Set(.list.Get().Shuffle!())
		.list.Repaint()
		}
	GetLibs()
		{
		libs = Object(.libs.Get())
		if libs[0] is "" or libs[0] is "(All)"
			libs = Libraries()
		libs.Remove("configlib")
		return libs
		}
	Overview_Click(row)
		{
		lines = .list.GetNumVisibleRows()
		up = Max((row - (lines / 2)).Int(), 0)
		down = Min((row + (lines / 2)).Int(), .list.Get().Size() - 1)
		.list.SetSelection(row)
		.list.ScrollRowToView(up)
		.list.ScrollRowToView(down)
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.resize_ovbar()
		}
	List_AfterSort()
		{
		.ovbar.ClearMarks()
		data = .list.Get()
		for i in data.Members()
			{
			if false is result = data[i].GetDefault("testrunnergui_result", false)
				continue
			.ovbar.AddMark(i, result is "X" ? CLR.ErrorColor : .successColor)
			}
		.ovbar.Update()
		}

	Destroy()
		{
		if .prevConditions isnt false
			.afterRunList()
		super.Destroy()
		}
	}
