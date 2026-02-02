// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
PageBaseCheck
	{
	CallClass(book, checkTmpFiles = false, checkTransaction? = false)
		{
		msgs = Object().Set_default(Object())
		checks = Object(
			[name: 'cursor', countFn: { Database.Cursors() }],
			[name: 'pubsub', countFn: PubSub.Count])
		if not Sys.SuneidoJs?()
			{
			checks.Add([name: 'defer', countFn: { Suneido.Info("windows.nDefer") }])
			checks.Add([name: 'timer', countFn: { Suneido.Info("windows.nTimer") }])
			}
		if checkTransaction?
			checks.Add([name: 'transaction', countFn: { Database.Transactions().Size() }])
		beforeOption = Suneido.GetDefault(#CurrentBookOption, "")
		Suneido.CheckingBook = true
		Finally(
			{
			errs = .ForeachBookOption(book)
				{ |ctrl, name|
				Suneido.CurrentBookOption = name
				.checkWindow(ctrl, name, checkTmpFiles, msgs, checks)
				}
			}, { Suneido.Delete(#CheckingBook) })
		Suneido.CurrentBookOption = beforeOption
		return Opt(errs, '\r\n') $ msgs.Values().Map({ it.Join('\r\n') }).Join('\r\n')
		}

	checkWindow(ctrl, name, checkTmpFiles, msgs, checks)
		{
		start = Timestamp()
		.doWithChecks(name, checks, msgs)
			{
			if Sys.SuneidoJs?()
				{
				Dialog(0, ctrl, beforeRun: { |d|
					// let the defer through browser so that the testing controls will
					// get built on the broswer side before destroy
					_forceOnBrowser = true
					.deferClose(d)
				})
				}
			else
				Window(ctrl, show: false).Destroy()
			}

		// Only for ContinuousTests, because image temp file is shared with other windows,
		// it will not be cleaned up after new window destroyed
		if not checkTmpFiles or
			Suneido.Member?('Persistent') or not msgs.tmpFileFound.Empty?()
			return
		leftOver = Dir(GetAppTempPath() $ '*.*', details:).
			Filter({ it.date > start and it.name !~ '^(csv|cur|hsperfdata)' })
		if not leftOver.Empty?()
			{
			msgs.tmpFileFound.Add(name $ " screen leaves temp file(s): " $
				leftOver.Map(
					{
					it.name $
						"\r\n(" $
						String(GetFile(GetAppTempPath() $ it.name, limit: 100)) $
						")\r\n"
					}).Join('\r\n'))
			}
		}

	deferClose(dialog)
		{
		Defer({
			active =SujsAdapter.CallOnRenderBackend(
				#GetRegisteredControl, GetActiveWindow())
			if not Same?(active, dialog) and active.Base?(Dialog)
				{
				.deferClose(dialog) // need another defer because active.Result will throw
				active.Result(false)
				}
			dialog.Result(false)
			})
		}

	doWithChecks(name, checks, msgs, block)
		{
		for check in checks
			check.before = (check.countFn)()
		block()
		for check in checks
			if 0 isnt diff = (check.countFn)() - check.before
				msgs[check.name].Add(name $ " didn't close " $ Plural(diff, check.name) $
					', threads: ' $ Thread.List().Join(','))
		}

	CheckingBook?()
		{
		return Suneido.GetDefault(#CheckingBook, false)
		}
	}
