// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	threadState: false
	Init()
		{
		if .threadState is false
			.threadState = [lastIdle: Date.End(), destroyed?: false]
		}

	threadBusy: false
	IdleAfterChange()
		{
		.threadState.lastIdle = Timestamp()
		if .PreThread() is false
			return
		if .threadBusy is true
			return
		Thread(.runThread, 'ScintillaAddonForThreadTasks-' $ .AddonName $ 'Thread')
		}

	runThread()
		{
		.threadBusy = true
		startTime = false
		Finally({
			while not .threadState.destroyed? and .IsOutdatedRecord(startTime)
				if .Destroyed?()
					.Destroy()
				else
					{
					startTime = .threadState.lastIdle
					_checkStop? = { .checkStop(.threadState, startTime) }
					.ThreadFn(startTime)
					}
			},
			{ .threadBusy = false })
		}

	checkStop(state, startTime)
		{
		// If the reference to state is lost, return true as the reference was most
		// likely destroyed, (returning true will stop the calling code)
		try
			return state.destroyed? or .IsOutdatedRecord(startTime)
		return true
		}

	IsOutdatedRecord(startTime)
		{
		return startTime < .threadState.lastIdle
		}

	Destroy()
		{
		.threadState.destroyed? = true
		.threadState.lastIdle = Date.End()
		}
	}
