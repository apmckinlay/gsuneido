// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Title:      "Suneido Type Checker"
	binaryPath: ""
	orderedSrc: ()
	tabImages: (ok: 0, warn: 1, error: 2)
	New(.libview)
		{
		super(.buildLayout())
		if false isnt tabs = .FindControl(#Tabs)
			tabs.SetImageList(
				[.tabImageSpec("checkmark.emf", CLR.ButtonGreen),
					.tabImageSpec("triangle-warning.emf", CLR.WarnColor),
					.tabImageSpec("cross.emf", CLR.ErrorColor)])
		.Defer(.annotate) // so users dont need to click Check the first time around
		.sub = PubSub.Subscribe(#LibraryRecordChange, .refreshIfStale)
		}

	tabImageSpec(name, color)
		{
		if Sys.SuneidoJs?()
			{
			code = IconFont().MapToCharCode(name)
			return Object(char: code.char, font: code.font, :color)
			}
		return [ImageResource(name), color]
		}

	noRecordOpenView()
		{
		return [#Vert, #Skip,
			[#Horz, #Skip,
				[#Static, "Please open a record to type check.", size: "+4",
					weight: #bold], #Skip], #Skip]
		}

	buildLayout()
		{
		if false is .libview.CurrentName()
			return .noRecordOpenView()

		.lib = .libview.CurrentTable()
		.rec = .libview.CurrentName()

		.orderedSrc = TypeCheckHelper.OrderedSrc(.libview.CurrentName(),
			skipLineageOrLibName: .isRecordUnloadedOrAFunction())

		.buildTabs(tabs = [#Tabs, close_button: false])

		.binaryPath = TypeCheckHelper.BinaryPath()
		return [#Vert, tabs,
			[#Horz,
				[#Record,
					[#Vert,
						[#Pair, [#Static, "TypeChecker Binary"],
							TypeCheckerBinaryPicker(.binaryPath)], #Skip]],
				[#Button, #Check], [#Skip, xstretch: 1], [#Button, #Policy], [#Skip]],
			[#TodoOutput name: #diagnostics, readonly:],
			[#Horz, [#Skip, medium:], [#Static, "", name: #timeElapsed, xstretch: 1]]]
		}

	buildTabs(tabs)
		{
		i = .orderedSrc.Size()
		while i > 0
			{
			i -= 1
			x = .orderedSrc[i]
			tabs.Add(
				[#CodeView data: [text: x.src, name: x.name, table: .lib], Tab: x.name,
					readonly:])
			}
		}

	On_Policy()
		{
		if false isnt policy2 = TypeCheckerPolicyDialog(TypeCheckHelper.Policy())
			{
			TypeCheckHelper.SetPolicy(policy2)
			.annotate() // re-run with new policy applied
			}
		}

	On_Check()
		{
		if false isnt browse = .FindControl(#TypeCheckerBinary)
			{
			.binaryPath = browse.Get()
			TypeCheckHelper.SetBinaryPath(.binaryPath)
			}
		.annotate()
		}

	Activate()
		{
		.refreshIfStale()
		}

	refreshIfStale()
		{
		if .orderedSrc.Empty?()
			return
		prev = .orderedSrc
		leaf = prev[prev.Size() - 1].name
		try
			.orderedSrc = TypeCheckHelper.OrderedSrc(leaf)
		catch
			return // class missing or otherwise unavailable
		if .orderedSrcChanged?(prev, .orderedSrc)
			.annotate()
		}

	orderedSrcChanged?(a, b)
		{
		if a.Size() isnt b.Size()
			return true
		for (i = 0; i < a.Size(); i += 1)
			if a[i].name isnt b[i].name or a[i].src isnt b[i].src
				return true
		return false
		}

	isRecordUnloadedOrAFunction()
		{
		skipLineageOrLibName = false
		src = Query1Cached(.lib, name: .rec, group: -1).text
		if not Libraries().Has?(.lib) or Function?(src.Compile())
			skipLineageOrLibName = .lib

		return skipLineageOrLibName
		}

	annotate()
		{
		if false is tctrl = .FindControl(#Tabs)
			return
		.ensureConstructed(tctrl)

		if not TypeCheckHelper.BinaryExists?()
			{
			.AlertError("Type Checker",
				"Binary not found at:\n" $ TypeCheckHelper.BinaryPath())
			return
			}

		// response is in the same order as the request: base->child->...->grandchild.
		// tabs are in the opposite order (leaf-first), so reverse-index when splicing.
		response = ""
		method = TypeCheckerMethods.Annotate
		try
			{
			elapsed = Timer()
				{
				response = TypeCheckHelper.Run(.libview.CurrentName(), method,
					policy: TypeCheckHelper.Policy(),
					skipLineageOrLibName: .isRecordUnloadedOrAFunction())
				}
			.updateStatus(elapsed, response.GetDefault(#version, Date.Begin()))
			diagnostics = response.GetDefault(#diagnostics, false)
			.updateTabSources(tctrl, response.GetDefault(#result, false))
			.updateTabStatus(tctrl, diagnostics)
			.showDiagnostics(diagnostics)
			}
		catch (e)
			{
			names = .orderedSrc.Map({ it.name }).Join(", ")
			msg_limit = 200
			AlertError(
				"suneidotypes: failed to decode response\n\n" $ "Exception:\n" $
					String(e) $ "\n\n" $ "Request method: " $ method $ "\n" $
					"Request sources: " $ names $ "\n\n" $ "Response:\n" $
					String(response[..msg_limit]))
			}
		}

	ensureConstructed(ctrl)
		{
		n = ctrl.GetAllTabCount()
		for idx in ..n
			ctrl.ConstructTab(idx)
		// realize every tab so its class-outline tree gets an initial Resize
		selected = ctrl.GetSelected()
		for idx in ..n
			ctrl.Select(idx, keepFocus:)
		ctrl.Select(selected, keepFocus:)
		}

	updateStatus(elapsed, version)
		{
		if false is sctrl = .FindControl(#timeElapsed)
			return
		elapsed = elapsed.Round(4)
		sctrl.Set(
			"Type checker as of " $ String(version) $ " took " $ String(elapsed) $ 's')
		}

	updateTabSources(tabs, annotatedSrc)
		{
		if annotatedSrc is false
			return
		n = tabs.GetAllTabCount()
		if annotatedSrc.Size() isnt n or not Number?(n)
			return
		for (i = 0; i < n; i += 1)
			{
			ed = tabs.GetControl(i).Editor
			state = ed.GetState()
			ed.Set(annotatedSrc[n-1-i])
			ed.SetState(state)
			}
		}

	updateTabStatus(tabs, diagnostics)
		{
		status = Object()
		if diagnostics isnt false
			{
			for w in diagnostics.GetDefault(#warnings, #())
				if status.GetDefault(String(w.class), false) isnt #error
					status[String(w.class)] = #warn
			for err in diagnostics.GetDefault(#errors, #())
				status[String(err.class)] = #error
			}
		n = tabs.GetAllTabCount()
		for (i = 0; i < n; i += 1)
			{
			name = .orderedSrc[n-1-i].name
			which = status.GetDefault(name, #ok)
			tabs.SetImage(i, .tabImages[which])
			}
		}

	showDiagnostics(diagnostics)
		{
		if false is dctrl = .FindControl(#diagnostics)
			return

		errors, warnings = TypeCheckHelper.FormatDiagnostics(diagnostics)
		dctrl.Set(Opt(errors.Join("\n"), "\n") $ warnings.Join("\n"))
		}

	Scintilla_DoubleClick(source)
		{
		if false is dctrl = .FindControl(#diagnostics)
			return
		if source isnt dctrl // ignore double-clicks in code-view editors
			return
		if false is loc = .parseDiagnosticLine(source.GetLine())
			return
		.gotoLocation(loc.class, loc.line)
		}

	parseDiagnosticLine(text)
		{
		// formatDiagnostic produces: "KIND: Class.Method:Line Msg"
		rest = text.AfterFirst(": ") // "Class.Method:Line Msg"
		classname = rest.BeforeFirst('.') // "Class"
		lineStr = rest.AfterFirst(':').BeforeFirst(' ') // "Line"
		try
			return Object(class: classname, line: Number(lineStr))
		catch
			return false
		}

	gotoLocation(classname, line)
		{
		if false is tabs = .FindControl(#Tabs)
			return
		n = tabs.GetAllTabCount()
		for (j = 0; j < .orderedSrc.Size(); j += 1)
			if .orderedSrc[j].name is classname
				{
				tab_idx = n - 1 - j
				tabs.Select(tab_idx)
				if false isnt cv = tabs.GetControl(tab_idx)
					cv.Editor.GotoLine(line - 1) // Scintilla is 0-indexed
				return
				}
		}

	Destroy()
		{
		TypeCheckHelper.StopServer()
		if .Member?(#sub)
			.sub.Unsubscribe()
		super.Destroy()
		}
	}
