// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	LoopedAddons: #()
	baseAddons: #(AccessEditMonitor: false, AccessRefresh: false)
	New(.access, additionalAddons = #())
		{
		.LoopedAddons = Object()
		defaultTimeout = 240 // 60 (minutes) * 4 (hours) = 240 minutes
		.timeoutMin = Database.Info().GetDefault(#timeoutMin, defaultTimeout)
		for addon, options in .baseAddons.Copy().MergeNew(additionalAddons)
			{
			addonInstance = Global(addon)(.access, options)
			if addonInstance.RequirementsMet?()
				.LoopedAddons[addon] = addonInstance
			}
		}

	subs: ()
	Subscriber(enabled? = false)
		{
		if .LoopedAddons.Empty?()
			return
		.subs.Each(#Unsubscribe)
		if enabled?
			.subs = [
				PubSub.Subscribe(#WindowActivated, .forceRun)
				PubSub.Subscribe(#WindowInactivated, .forceStop)
				]
		else
			{
			.subs = #()
			.Stop()
			}
		}

	forceRun()
		{
		if .runPublishFunc?()
			{
			.timeout = Date().Plus(minutes: .timeoutMin)
			.loopFunc()
			}
		}

	forceStop()
		{
		if .runPublishFunc?()
			.Stop()
		}

	runPublishFunc?()
		{
		try
			return not .LoopedAddons.Empty?() and not .access.EditMode?() and
				.access.Window.Hwnd is GetActiveWindow()
		catch (e)
			{
			SuneidoLog('ERROR: (CAUGHT) ' $ e, calls:, caughtMsg: 'addon stopped')
			return false
			}
		}

	timer: false
	firstRun?: true
	Start()
		{
		if .LoopedAddons.Empty?()
			return
		.timeout = .key = false
		.init()
		if .firstRun?
			{
			.firstRun? = false
			.Subscriber(enabled?:)
			}
		if .timer is false
			.loop()
		}

	init() { .LoopedAddons.Each(#Init) }

	loop()
		{
		delay = 60_000 /*= 1 minute */
		.timer = .access.Delay(delay, .loopFunc, uniqueID: 'AccessLoopAddonManager')
		}
	loopFunc()
		{
		if not .timedOut?()
			{
			if .access.RecordSet?()
				.LoopedAddons.Each(#RunAddon)
			.loop()
			}
		}

	key: false
	timeout: false
	TimedOutMsg: `INFO: Access Timeout - closing idle connection `
	timedOut?()
		{
		if .access.RecordSet?() and .key isnt curKey = .access.GetLockKey()
			{
			.key = curKey
			.timeout = Date().Plus(minutes: .timeoutMin)
			}
		else if .timeout isnt false and Date() > .timeout
			{
			.Subscriber()
			// Exit(true) is required, in order to prevent affecting other windows
			// timeouts. Without it, other open windows may refresh their timeout timers
			// delaying their timeouts further
			SuneidoLog(.TimedOutMsg $ Database.SessionId())
			Delay(10000, /*= standard timeout exit delay*/) { ExitClient(true) }
			Alert('lost connection:', 'Fatal Error', 0)
			ExitClient(true)
			return true // Exit kills the client, return true stops the loop
			}
		return false
		}

	Stop()
		{
		if .timer isnt false
			{
			.timer.Kill()
			.timer = false
			}
		}
	}
