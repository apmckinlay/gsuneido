// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	title: "Auto Format and Send"
	CallClass(libview)
		{
		libview.Save()
		lib = libview.CurrentTable()
		name = libview.CurrentName()
		svcTable = SvcTable(lib)

		if OptContribution(#CheckUpdateBuildTime, function() { return false })()
			.alert("Sending code during potential update building time\r\n")

		if lib is "" or name is "" or false is rec = svcTable.Get(name)
			return .alert("No current record to format")

		if "" isnt why = .unsendable(rec)
			return .alert(name $ ' ' $ why)

		if false is formatted = .format(name, rec.text)
			return false

		if formatted is rec.text
			return .alert(name $ " is already formatted")

		if not AstFmtEquals?(rec.text, formatted)
			return .alert(
				"Formatting would change the code, not just the layout.\n" $
					"Nothing was sent - please report this.", warn:)

		if false is settings = SvcSettings()
			return .alert("Invalid version control settings")

		result = .run(
			Object(:settings, :svcTable, :lib, :name, :formatted, oldText: rec.text),
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
			return .apply(config.svcTable, config.name, config.formatted)

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
			return "Offline from version control - you can still apply locally with OK"

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
			return "is a new record - send it with Version Control first"
		if rec.lib_modified isnt ""
			return "has unsent changes - send them with Version Control first," $
				" so that only formatting is committed here"
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
		Title: "Auto Format and Send"
		New(config)
			{
			super(.layout(config))
			.config = config
			.sendBtn = .FindControl(#Send_Changes)
			.warn = .FindControl(#warning)
			.setBlocked(config.blocked)
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
				[#Horz, #Fill,
					[#Button, #OK, tip: "Update " $ name $ " locally without sending"],
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
