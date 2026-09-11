// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	title: "Formatting Tools"
	CallClass(libview)
		{
		libview.Save()
		lib = libview.CurrentTable()
		name = libview.CurrentName()
		svcTable = SvcTable(lib)

		if lib is "" or name is "" or false is rec = svcTable.Get(name)
			return .alert("No current record to format")

		if false is formatted = .format(name, rec.text)
			return false

		if formatted is rec.text
			return .alert(name $ " is already formatted")

		if not AstFmtEquals?(rec.text, formatted)
			return .alert(
				"Formatting would change the code, not just the layout.\n" $
					"Nothing was done - please report this.", warn:)

		// when nothing is committed only format locally
		if "" isnt why = .unsendable(rec)
			{
			.applyLocally(libview, name, formatted)
			return .alert(name $ ' ' $ why $ " so it was formatted locally")
			}

		if OptContribution(#CheckUpdateBuildTime, function() { return false })()
			.alert("Sending code during potential update building time\r\n")

		if false is settings = SvcSettings()
			return .alert("Invalid version control settings")

		result = .run(
			Object(:settings, :svcTable, :lib, :name, :formatted, :libview,
				oldText: rec.text),
			libview.Window.Hwnd)

		SvcSocketClient().Close()

		return result
		}

	run(config, hwnd)
		{
		config.userid = config.settings.svc_userId
		config.desc = .defaultDesc()
		config.blocked = .Blocked(config.settings, config.lib)

		if false is choice = ToolDialog(hwnd, [.confirm, config], title: .title,
			closeButton?: false)
			return false

		if .RecordChanged?(config.svcTable, config.name, config.oldText)
			return .alert(
				config.name $ " changed while the dialog was open - nothing was" $
					" done.\nPlease run Auto Format and Send again.", warn:)
		if choice.send isnt true
			return .applyLocally(config.libview, config.name, config.formatted)

		config.userid = choice.userid
		config.desc = choice.desc
		return .sendChanges(config)
		}

	sendChanges(config)
		{
		svc = .connect(config.settings)
		if String?(svc)
			return .alert("You are offline from version control:\n\n" $ svc, warn:)

		// re-check in case someone sent while the dialog was open
		if svc.Outstanding?([config.lib])
			return .alert(
				"Please get the master changes before sending.\n" $
					"Nothing was sent - use Version Control to get " $ config.lib $
					" up to date, then try again.", warn:)
		.apply(config.svcTable, config.name, config.formatted)
		return .send(svc, config.lib, config.name, config.userid, config.desc)
		}

	// returns "" when it is safe to send, otherwise why it is not
	Blocked(settings, lib)
		{
		svc = .connect(settings)
		if String?(svc)
			return "Offline from version control - you can still format locally"

		if svc.Outstanding?([lib])
			return lib $ " has changes from others that you have not got" $
				" - use the V button to get them"
		return ""
		}

	// the diff and the formatted text are only valid for the text they were built
	// from, so getting master changes invalidates them
	RecordChanged?(svcTable, name, oldText)
		{
		rec = svcTable.Get(name)
		return rec is false or rec.text isnt oldText
		}

	defaultDesc()
		{
		return #autoformatted
		}

	unsendable(rec)
		{
		if rec.lib_committed is ""
			return "is a new record,"
		if rec.lib_modified isnt ""
			return "has unsent changes,"
		return ""
		}

	format(name, text)
		{
		if not Compilable?(text)
			return .alert(name $ " does not compile so it was not formatted")

		try
			return AstFormatter(text).Replace("\r?\n", "\r\n")
		catch (e)
			{
			SuneidoLog("ERROR: AutoFormatSend: " $ e)
			return .alert("Format failed on " $ name, warn:)
			}
		}

	confirm: Controller
		{
		Title: "Formatting Tools"
		New(config)
			{
			super(.layout(config))
			.config = config
			.sendBtn = .FindControl(#Send_Changes)
			.warn = .FindControl(#warning)
			.setBlocked(config.blocked)
			.showTodos(config)
			}

		showTodos(config)
			{
			panes = .FindControl(#Diff).ListControls()
			.setTodo(#todoCurrent, config, panes[0].Get())
			.setTodo(#todoFormatted, config, panes[1].Get())
			}

		setTodo(ctrlName, config, text)
			{
			if false is ctrl = .FindControl(ctrlName)
				return
			ctrl.SetWordChars(
				"_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ?!:")
			ctrl.Set(.todos(config.lib, config.name, text).Join('\n'))
			}

		Width: 90
		TabWidth: 4
		todos(lib, recName, text)
			{
			result = Object()
			lines = text.Lines()
			for (i = 0; i < lines.Size(); ++i)
				{
				line = lines[i].Tr('\r')
				tag = lib $ ':' $ recName $ ':' $ (i + 1) $ ' '
				if line =~ `/[*/] *[A-Z][A-Z][A-Z]+`
					result.Add(tag $ line.Trim())
				else if .Width < w = line.Replace('\t', ' '.Repeat(.TabWidth)).Size()
					result.Add(tag $ "LONG LINE (" $ w $ ') ' $ line.Trim())
				}
			return result
			}

		Scintilla_DoubleClick(source)
			{
			if source isnt .FindControl(#todoCurrent) and
				source isnt .FindControl(#todoFormatted)
				return // ignore double clicks in the diff panes themselves
			if false is row = .todoRow(source.GetLine())
				return
			for pane in .FindControl(#Diff).ListControls()
				pane.GotoLine(row)
			}

		todoRow(text)
			{
			try
				return Number(text.BeforeFirst(' ').AfterLast(':')) - 1
			catch
				return false
			}

		setBlocked(blocked)
			{
			.blocked = blocked
			.sendBtn.SetEnabled(blocked is "")
			.warn.Set(blocked)
			}

		layout(config)
			{
			name = config.name
			return [#Vert,
				[#Horz,
					[#EnhancedButton, command: 'T', image: "T.emf", imagePadding: 0.05,
						mouseEffect:, ystretch: 0, tip: "Open the Test Runner"],
					#Skip,
					[#EnhancedButton, command: 'V', image: "V.emf", imagePadding: 0.05,
						mouseEffect:, ystretch: 0, tip: "Open Version Control"],
					#Skip,
					#(Static, "User:"),
					#Skip,
					[#Field, name: #user, set: config.userid, width: 12, xstretch: 0,
						cue: "user id",
						status: "the user id the formatting will be sent as"],
					#Skip,
					#(Static, "Message:"),
					#Skip,
					[#Field, name: #desc, set: config.desc, xstretch: 1,
						cue: "commit message",
						status: "the message the formatting will be sent with"],
					#Skip,
					[#Button, "Send Changes",
						tip: "Update " $ name $
							" and send the formatting to version control"]],
				#Skip,
				[#Static, "", name: #warning, textStyle: #warn, xstretch: 1],
				#Skip,
				[#Diff2, config.oldText, config.formatted, config.lib, name, #Current,
					#Formatted],
				#Skip,
				[#Horz,
					[#Vert,
						[#Static, "Todo - Current", name: #todoTitleCurrent],
						[#TodoOutput, name: #todoCurrent, readonly:, ymin: 90,
							ystretch: .25]],
					#Skip,
					[#Vert,
						[#Static, "Todo - Formatted", name: #todoTitleFormatted],
						[#TodoOutput, name: #todoFormatted, readonly:, ymin: 90,
							ystretch: .25]]],
				#Skip,
				[#Horz, #Fill,
					[#Button, "Format Locally", command: #OK,
						tip: "Update " $ name $ " locally without sending"],
					#Skip, [#Button, #Cancel, tip: "Close without changing anything"]]
				xmin: 900]
			}

		On_T()
			{
			TestRunnerGui()
			}

		On_V()
			{
			SvcControl()
			if LibViewAutoFormatSend.RecordChanged?(.config.svcTable, .config.name,
				.config.oldText)
				{
				.AlertInfo(.Title,
					.config.name $ " was updated by getting changes," $
						" so the diff is out of date.\n" $
						"Please run Auto Format and Send again.")
				.Window.Result(false)
				return
				}
			.setBlocked(LibViewAutoFormatSend.Blocked(.config.settings, .config.lib))
			}

		On_Send_Changes()
			{
			if .blocked isnt ""
				{
				.AlertInfo(.Title, .blocked)
				return
				}
			if "" is userid = .FindControl(#user).Get().Trim()
				{
				.AlertInfo(.Title, "Please enter a user id to send as")
				return
				}
			if "" is desc = .FindControl(#desc).Get().Trim()
				{
				.AlertInfo(.Title, "Please enter a commit message to send with")
				return
				}
			.Window.Result(Object(send:, :userid, :desc))
			}

		On_OK()
			{
			.Window.Result(Object(send: false))
			}

		On_Cancel()
			{
			.Window.Result(false)
			}
		}

	// paste through the editor, so that Ctrl+Z can undo the formatting
	applyLocally(libview, name, formatted)
		{
		if libview.CurrentName() isnt name
			return .apply(SvcTable(libview.CurrentTable()), name, formatted)
		editor = libview.Editor
		line = editor.LineFromPosition()
		editor.PasteOverAll(formatted)
		editor.GotoLine(line)
		libview.Save()
		return true
		}

	apply(svcTable, name, formatted)
		{
		Transaction(update:)
			{|t|
			rec = svcTable.Get(name, :t)
			svcTable.Update(rec, :t, newText: formatted)
			}
		return true
		}

	connect(settings) // returns a String (the reason) when offline
		{
		try
			{
			SvcSocketClient().RetryState() // re-checks after a get
			svc = Svc(server: settings.svc_server, local?: settings.svc_local? is true)
			if "" isnt status = svc.CheckSvcStatus()
				return status
			return svc
			}
		catch (e)
			return e
		}

	send(svc, lib, name, userid, desc)
		{
		if false is svc.SendLocalChanges([Object(type: ' ', :name, :lib)], desc, userid)
			return .alert(
				"Send failed - someone may have sent new changes.\n" $
					"Please refresh and get changes in Version Control, then try again.",
				warn:)
		return .alert(name $ " formatted and sent")
		}

	alert(msg, warn = false)
		{
		Alert(msg, .title, flags: warn ? MB.ICONWARNING : MB.ICONINFORMATION)
		return false
		}
	}
