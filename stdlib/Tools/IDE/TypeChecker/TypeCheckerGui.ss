// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "TypeCheckRunner"
	binaryPath: ""
	symPass: ''
	symFail: 'X'
	symWarn: '!'
	symSkip: '?'
	CallClass()
		{
		GotoPersistentWindow('TypeCheckerGui', TypeCheckerGui)
		}
	maxThreads: 16
	workerTimeoutSecs: 86400 // WaitGroup returns the instant all workers finish
	New(lib = "")
		{
		super(.controls())
		.libs = .FindControl('ChooseList')
		.list = .FindControl('List')
		.list.SetMultiSelect(true)
		.ovbar = .FindControl('ovbar')
		.ovbar.SetTopMargin(GetSystemMetrics(SM.CXHSCROLL))
		.statusbar = .Vert.Status
		.stop = .FindControl('stop_on_error')
		.stop.Set(false)
		.time = .FindControl('total_time')
		.nerrors = .FindControl('total_nerrors')
		.nwarnings = .FindControl('total_nwarnings')
		.threads = .FindControl('thread_count')
		.threads.Set(String(.defaultThreads()))
		.okColor = CLR.ButtonGreen
		.warnColor = CLR.WarnColor
		.errColor = CLR.ErrorColor
		.skipColor = CLR.Inactive
		if lib isnt ""
			{
			.libs.Set(lib)
			.set_lib(lib)
			}
		}

	// ----- layout -----------------------------------------------------

	columns: (typecheckrunner_result, typecheckrunner_lib,
		typecheckrunner_name, typecheckrunner_time,
		typecheckrunner_nerrors, typecheckrunner_nwarnings)

	controls()
		{
		libs = Libraries().Add("(All)", at: 0)
		.binaryPath = TypeCheckHelper.BinaryPath()
		return Object('Vert',
			Object('Horz', 'Skip', .top(libs)),
			Object('Horz',
				Object('ListStretch', columns: .columns,
					stretchColumn: 'typecheckrunner_name',
					columnsSaveName: 'TypeCheckerGui'),
				Object('OverviewBar', name: 'ovbar',
					priorityColor: CLR.ErrorColor)),
			.totalsRow(),
			"Statusbar")
		}

	top(libs)
		{
		row1 = Object('Horz', Object('ChooseList', libs), 'Skip',
			#(Button, 'Run All'), 'Skip',
			#(Button, 'Continue', tip: 'Run from selected'), 'Skip',
			#(Button, 'Run Failed'), 'Skip',
			#(Button, 'Policy', tip: 'Configure strictness levels'), 'Skip',
			#(Pair (Static 'Threads')
				(Field name: thread_count, width: 3,
				tip: 'Concurrent checks (1 = sequential)')), 'Skip',
			#(CheckBox, 'Stop on Error', name: stop_on_error), 'Skip')
		row2 = Object('Horz',
			#(Static, 'Binary '),
			TypeCheckerBinaryPicker(.binaryPath), 'Skip')
		return Object('Flow', Object(row1, row2), skip: 0)
		}

	totalsRow()
		{
		return Object('Horz',
			#(Skip 4),
			#(Static "Totals:", weight: "bold"),
			'Skip',
			#(Pair (Static "Time")
				(Field width: 8, readonly:, name: total_time)),
			'Skip',
			#(Pair (Static "Errors")
				(Field width: 5, readonly:, name: total_nerrors)),
			'Skip',
			#(Pair (Static "Warnings")
				(Field width: 5, readonly:, name: total_nwarnings)),
			'Skip')
		}

	// ----- commands / menu --------------------------------------------

	Commands: (
		(Run_Selected,		"F9")
		(Go_To_Definition,	"F12")
		(Close,				"",			"Close this window")
		(Refresh,			"F5")
		)
	Menu:
		(
		('&File',
			'&Refresh', '',
			'&Close')
		)

	// ----- library selection / list population ------------------------

	NewValue(value, source)
		{
		if not .Member?(#libs) or source isnt .libs or not source.Valid?()
			return
		.libs.Dirty?(false)
		.Defer({ .set_lib(value) })
		}

	On_Refresh()
		{
		.Defer({ .set_lib(.libs.Get()) })
		}

	set_lib(lib)
		{
		.reset_statusbar()
		if lib is ""
			{
			.list.Set(Object())
			return
			}
		recs = Object()
		addRec = {|name|
			recs.Add(Object(typecheckrunner_lib: lib,
				typecheckrunner_name: name))
			}
		if lib is '(All)'
			for one in Libraries()
				.foreachRecord(one, addRec)
		else
			.foreachRecord(lib, addRec)
		.list.Set(recs)
		.reset(clearMarks?:)
		}

	foreachRecord(lib, block)
		{
		QueryApply(lib $ " where name =~ '^[A-Z]' and group is -1 sort name")
			{|x|
			try
				{
				if x.GetDefault(#text, "") is ""
					continue
				if not Class?(Global(x.name))
					continue
				block(x.name)
				}
			catch
				{ }
			}
		}

	// ----- reset helpers ----------------------------------------------

	reset(clearMarks? = false)
		{
		.time.Set("")
		.nerrors.Set("")
		.nwarnings.Set("")
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

	reset_row(i)
		{
		.list.ClearHighlight(i)
		.ovbar.RemoveMark(i)
		rec = .list.Get()[i]
		rec.typecheckrunner_result = ""
		rec.typecheckrunner_time = ""
		rec.typecheckrunner_nerrors = ""
		rec.typecheckrunner_nwarnings = ""
		}

	// ----- run buttons ------------------------------------------------

	On_Run_All()
		{
		if not .syncBinaryPath()
			return
		if .libs.Get() is ""
			.libs.Set('(All)')
		// set_lib + run_list must run together so run_list sees the
		// freshly-populated list, hence one Defer block, not two.
		.Defer({
			.set_lib(.libs.Get())
			.run_list(.list.Get().Members())
			})
		}

	On_Continue()
		{
		selected = .list.GetSelection()
		if selected is #()
			return
		first = selected.Sort!()[0]
		while first < .list.GetNumRows() and
			(.list.GetRow(first).typecheckrunner_result is .symFail or
			 .list.GetRow(first).typecheckrunner_result is .symSkip)
			first++
		runlist = Object()
		for (i = first; i < .list.GetNumRows(); i++)
			runlist.Add(i)
		.run_list(runlist)
		}

	On_Run_Failed()
		{
		runlist = Object()
		data = .list.Get()
		for i in data.Members()
			{
			r = data[i].GetDefault(#typecheckrunner_result, "")
			if r is .symFail or r is .symSkip
				runlist.Add(i)
			}
		if runlist isnt #()
			.run_list(runlist)
		}

	On_Run_Selected()
		{
		selected = .list.GetSelection()
		if selected is #()
			return
		.run_list(selected.Sort!())
		}

	On_Policy()
		{
		if false isnt next = TypeCheckerPolicyDialog(TypeCheckHelper.Policy())
			TypeCheckHelper.SetPolicy(next)
		}

	// ----- run kick-off -----------------------------------------------

	defaultThreads()
		{
		n = 4
		try
			{
			env = Getenv('NUMBER_OF_PROCESSORS')   // Windows; "" elsewhere
			if env isnt ''
				n = Number(env) - 1   // leave a core for the UI thread
			}
		catch
			{ }
		return Min(Max(n, 1), .maxThreads)
		}

	readThreadCount()
		{
		n = 1
		try
			n = Number(.threads.Get())
		catch
			{ n = 1 }
		if not Number?(n)
			n = 1
		return Min(Max(n.Int(), 1), .maxThreads)
		}

	syncBinaryPath()
		{
		if false isnt browse = .FindControl(#TypeCheckerBinary)
			{
			TypeCheckHelper.SetBinaryPath(browse.Get())
			}
		if not TypeCheckHelper.BinaryExists?()
			{
			.AlertError("Type Checker",
				"Binary not found at:\n" $ TypeCheckHelper.BinaryPath())
			return false
			}
		return true
		}

	run_list(runlist)
		{
		if .libs.Get() is "" or not .syncBinaryPath()
			return
		.reset()
		for i in runlist
			.reset_row(i)
		.list.Repaint()
		this.SetEnabled(false)
		.runWorker(runlist, .list.Get(), .stop.Get() is true ? 1 : .readThreadCount())
		}

	startupDelayMs: 100
	runWorker(runlist, data, nThreads)
		{
		ui = Object(
			ovbarMarks: Object(),
			repaintRows: Object(),
			state: 'starting',
			startTime: Date(),
			startupError: false,
			stopOnError: .stop.Get() is true,
			stopRequested: false)
		Suneido.typecheckUi = ui
		.Delay(.startupDelayMs, .updateUi, uniqueID: 'updateUi')
		.spawnRunner(runlist, data, nThreads, ui)
		}

	spawnRunner(runlist, data, nThreads, ui)
		{
		Thread(name: 'type check runner')
			{
			try
				.runOnWorker(runlist, data, nThreads, ui)
			catch (e)
				{
				ui.startupError = String(e)
				ui.state = 'failed'
				}
			}
		}

	// Coordinator thread. Warms the shared server and the global class table
	// (single-threaded, so concurrent Global() loads can't contend), then fans
	// the rows out across nThreads workers. The UI thread is untouched it just
	// polls Suneido.typecheckUi as before.
	runOnWorker(runlist, data, nThreads, ui)
		{
		try
			TypeCheckHelper.Server()
		catch (e)
			{
			ui.startupError = String(e)
			ui.state = 'failed'
			return
			}
		.prefetchLineage(data, runlist)
		ui.state = 'running'
		if nThreads <= 1
			.runStride(data, runlist, 0, 1, ui)   // proven sequential path
		else
			.runParallel(data, runlist, nThreads, ui)
		ui.state = 'finished'
		}

	prefetchLineage(data, runlist)
		{
		names = Object()
		for i in runlist
			names.Add(data[i].typecheckrunner_name)
		try
			TypeCheckerLineage(names)   // warms the global class table for the workers
		catch
			{ }
		}

	runParallel(data, runlist, nThreads, ui)
		{
		n = runlist.Size()                    // instantiates the sequence
		jobs = Channel(Max(n, 1))
		for (j = 0; j < n; j++)
			jobs.Send(runlist[j])             // index access, same as runStride
		jobs.Close()
		wg = WaitGroup()
		for (k = 0; k < nThreads; k++)
			{
			wg.Add(1)
			Thread(name: 'typecheck worker ' $ k)
				{ .runFromChannel(data, jobs, wg, ui) }
			}
		wg.Wait(.workerTimeoutSecs)
		}

	runFromChannel(data, jobs, wg, ui)
		{
		forever
			{
			if ui.stopRequested
				break
			i = jobs.Recv()
			if i is jobs   // channel returns itself when closed + drained
				break
			try
				.runOne(data, i, ui)
			catch
				try .markSkipped(data, i, 0, ui)
			}
		wg.Done()
		}

	// round-robin: worker k handles runlist positions k, k+stride, k+2*stride, ...
	// no shared job queue, so zero dispatch contention.
	runStride(data, runlist, start, stride, ui)
		{
		for (j = start; j < runlist.Size(); j += stride)
			{
			if ui.stopRequested
				break
			i = runlist[j]
			try
				.runOne(data, i, ui)
			catch
				try .markSkipped(data, i, 0, ui)
			}
		}

	// ----- binary lifecycle -------------------------------------------

	// the shared -serve process is owned by TypeCheckHelper and stays warm
	// across runs; we just hand it our pre-resolved lineage chain.
	invokeBinary(className)
		{
		method = TypeCheckerMethods.Infer
		return TypeCheckHelper.Run(className, method, policy: TypeCheckHelper.Policy(),
			restartOnError?: false)
		}

	// ----- per-row execution ------------------------------------------

	runOne(data, i, ui)
		{
		rec = data[i]
		response = false
		err = false
		elapsed = Timer()
			{
			try
				response = .invokeBinary(rec.typecheckrunner_name)
			catch (e)
				err = String(e)
			}
		if err isnt false or response is false
			.markSkipped(data, i, elapsed, ui)
		else
			.applyResponse(data, i, response, elapsed, ui)
		}

	// Each row checks one class with its full base chain prepended, so
	// the binary may emit diagnostics tagged to bases too. Filter to
	// the leaf class so bases only count when their own row runs.
	applyResponse(data, i, response, elapsed, ui)
		{
		rec = data[i]
		counts = .countLeafDiagnostics(response, rec.typecheckrunner_name)
		verdict = .classifyResult(counts.nerr, counts.nwarn)
		rec.typecheckrunner_time = elapsed.Round(3)
		rec.typecheckrunner_nerrors = counts.nerr
		rec.typecheckrunner_nwarnings = counts.nwarn
		rec.typecheckrunner_result = verdict.result
		.stageRowUpdate(i, verdict, ui)
		}

	// reached when the class can't be loaded for lineage, or the binary/HTTP
	// call throws neither is a real type error, so don't count toward errs
	// and don't trip stop-on-error
	markSkipped(data, i, elapsed, ui)
		{
		rec = data[i]
		rec.typecheckrunner_time = elapsed is 0 ? "" : elapsed.Round(3)
		rec.typecheckrunner_result = .symSkip
		rec.typecheckrunner_nerrors = ""
		rec.typecheckrunner_nwarnings = ""
		.stageRowUpdate(i, Object(color: .skipColor, highlight: .skipColor,
			failed?: false), ui)
		}

	stageRowUpdate(i, verdict, ui)
		{
		ui.ovbarMarks.Add(Object(i, verdict.color))
		ui.repaintRows.Add(Object(i, verdict.highlight))
		if verdict.failed? and ui.stopOnError
			ui.stopRequested = true
		}

	countLeafDiagnostics(response, leaf)
		{
		diags = response.GetDefault(#diagnostics, false)
		if not Object?(diags)
			return Object(nerr: 0, nwarn: 0)
		return Object(
			nerr: .countForLeaf(diags.GetDefault(#errors, #()), leaf),
			nwarn: .countForLeaf(diags.GetDefault(#warnings, #()), leaf))
		}

	countForLeaf(list, leaf)
		{
		n = 0
		for d in list
			if String(d.GetDefault(#class, '')) is leaf
				n++
		return n
		}

	classifyResult(nerr, nwarn)
		{
		if nerr > 0
			return Object(result: .symFail, color: .errColor,
				highlight: .errColor, failed?:)
		if nwarn > 0
			return Object(result: .symWarn, color: .warnColor,
				highlight: false, failed?: false)
		return Object(result: .symPass, color: .okColor,
			highlight: false, failed?: false)
		}

	// ----- UI updates -------------------------------------------------

	timerDelay: 300 // ms
	updateUi()
		{
		ui = Suneido.typecheckUi
		touchedOv = .drainOvBar(ui)
		frontier = .drainList(ui)
		if touchedOv
			.ovbar.Update()
		if frontier isnt false
			{
			.list.Update()
			.list.ScrollRowToView(frontier)   // follow checking down the list
			}
		if ui.state is 'starting' or ui.state is 'running'
			{
			.Delay(.timerDelay, .updateUi, uniqueID: 'updateUi')
			return
			}
		.finishRun()
		}

	drainOvBar(ui)
		{
		touched = false
		while false isnt item = ui.ovbarMarks.Extract(0, false)
			{
			.ovbar.AddMark(@item)
			touched = true
			}
		return touched
		}

	drainList(ui)
		{
		frontier = false
		while false isnt item = ui.repaintRows.Extract(0, false)
			{
			.addRowHighlight(@item)
			if frontier is false or item[0] > frontier
				frontier = item[0]
			}
		return frontier
		}

	addRowHighlight(i, highlight = false)
		{
		if highlight isnt false
			.list.AddHighlight(i, highlight)
		.list.RepaintRow(i)   // queues invalidation; the single .list.Update() flushes
		}

	finishRun()
		{
		ui = Suneido.typecheckUi
		if ui.GetDefault(#startupError, false) isnt false
			{
			.AlertError("TypeChecker",
				"failed to start binary in serve mode:\n" $ ui.startupError)
			this.SetEnabled(true)
			return
			}
		elapsed = Date().MinusSeconds(ui.startTime)
		.time.Set(elapsed.RoundToPrecision(2) $ " sec")
		totals = .tallyTotals()
		.nerrors.Set(totals.errs)
		.nwarnings.Set(totals.warns)
		summary = .summarize(totals.errs, totals.warns, totals.nskip)
		.updateStatusBar(summary.color, summary.msg)
		this.SetEnabled(true)
		}

	tallyTotals()
		{
		errs = warns = nskip = 0
		for rec in .list.Get()
			{
			if rec.GetDefault(#typecheckrunner_result, "") is .symSkip
				nskip++
			ne = rec.GetDefault(#typecheckrunner_nerrors, 0)
			nw = rec.GetDefault(#typecheckrunner_nwarnings, 0)
			if Number?(ne)
				errs += ne
			if Number?(nw)
				warns += nw
			}
		return Object(:errs, :warns, :nskip)
		}

	summarize(errs, warns, nskip)
		{
		skipTail = nskip > 0 ? ", " $ nskip $ " unchecked" : ""
		if errs > 0
			return Object(
				color: .errColor,
				msg: " F A I L U R E - " $ errs $ " error(s), " $
					warns $ " warning(s)" $ skipTail)
		if warns > 0
			return Object(
				color: .warnColor,
				msg: " W A R N I N G S - 0 errors, " $ warns $
					" warning(s)" $ skipTail)
		if nskip > 0
			return Object(
				color: .skipColor,
				msg: " " $ nskip $ " unchecked, 0 errors, 0 warnings")
		return Object(
			color: .okColor,
			msg: " S U C C E S S - 0 errors, 0 warnings")
		}

	updateStatusBar(color, message)
		{
		.statusbar.SetBkColor(color)
		if message isnt ""
			.statusbar.Set(message)
		.statusbar.Update()
		}

	// ----- list interactions ------------------------------------------

	List_AllowCellEdit(@unused)
		{ return false }
	List_EditFieldReadonly(@unused)
		{ return true }
	List_WantNewRow()
		{ return false }

	contextMenu: ("Run\tF9", "Go To Definition\tF12", "Copy Name")
	List_ContextMenu(x, y)
		{
		if .list.Get().Empty?()
			return
		ContextMenu(.contextMenu).ShowCall(this, x, y)
		}
	List_DoubleClick(row, col)
		{
		if row is false or .list.GetCol(col) isnt 'typecheckrunner_name'
			return 0
		GotoLibView(.list.GetRow(row).typecheckrunner_name)
		return false
		}
	On_Context_Run()
		{ .On_Run_Selected() }
	On_Context_Go_To_Definition()
		{ .On_Go_To_Definition() }
	On_Go_To_Definition()
		{
		selected = .list.GetSelection()
		if selected.Size() isnt 1
			return
		GotoLibView(.list.Get()[selected[0]].typecheckrunner_name)
		}
	On_Context_Copy_Name()
		{
		selected = .list.GetSelection()
		if selected.Empty?()
			return
		ClipboardWriteString(.list.Get()[selected[0]].typecheckrunner_name)
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
			r = data[i].GetDefault("typecheckrunner_result", false)
			if r is false or r is ""
				continue
			color = .warnColor
			if r is .symFail
				color = .errColor
			else if r is .symSkip
				color = .skipColor
			.ovbar.AddMark(i, color)
			}
		.ovbar.Update()
		}

	Destroy()
		{
		if Suneido.Member?(#typecheckUi)
			{
			Suneido.typecheckUi.stopRequested = true // let live workers exit
			Suneido.Delete(#typecheckUi) // dont pollute the global object
			}
		TypeCheckHelper.StopServer()
		super.Destroy()
		}
	}
