// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	interval: 5 // secs

	CallClass(msg, block)
		{
		task = .createTask(msg)
		Finally({
			while not task.Canceled?() and block()
				task.PauseIfNeeded()
			}, {
			task.Finish()
			})
		return not task.Canceled?()
		}

	createTask(msg)
		{
		if not Sys.SuneidoJs?()
			return .dummyTask

		task = Suneido.GetInit(#TaskWithPause, { new this(msg) })
		task.Ref()
		return task
		}

	// for server
	dummyTask: class
		{
		PauseIfNeeded() {}
		Ref() {}
		Finish() {}
		Canceled?() { return false }
		}

	now()
		{
		return Date()
		}

	New(msg)
		{
		.prev = .now()
		.setup(msg)
		}

	setup(msg)
		{
		.UniqueId = SuRenderBackend().NextId()
		SuRenderBackend().Register(.UniqueId, this)
		SuRenderBackend().RecordAction(false, 'DoTaskWithPauseClient',
			args: [.UniqueId, msg])
		}

	ref: 0
	Ref()
		{
		.ref++
		}

	PauseIfNeeded()
		{
		cur = .now()
		if cur.MinusSeconds(.prev) >= .interval
			{
			.prev = cur
			.pause()
			}
		}

	cancel?: false
	pause()
		{
		SuRenderBackend().RecordAction(.UniqueId, 'Pause', args: #())
		JsWebSocketServer.MessageLoop()
		}

	DonePause(.cancel?)
		{
		throw WebSocketHandler.QUITLOOP
		}

	Finish()
		{
		if --.ref is 0
			.cleanup()
		}

	cleanup()
		{
		Suneido.Delete(#TaskWithPause)
		SuRenderBackend().RecordAction(.UniqueId, 'Finish', args: #())
		SuRenderBackend().UnRegister(.UniqueId)
		}

	Canceled?()
		{
		return .cancel?
		}
	}
