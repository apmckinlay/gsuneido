// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{
		.timerId = 1
		.timers = Object()
		.delayQueue = Object()
		.delayToCancel = Object()
		}

	ReserveId()
		{
		return .timerId++
		}

	// NOTO: this is not thread safe.
	// Not sure if we should have multi-thread per connection on server side
	SetTimer(ms, cb, id = false, _forceOnBrowser = false)
		{
		if id is false
			id = .timerId++
		if ms is 0 and forceOnBrowser isnt true
			.delayQueue.Add(Object(:id, :cb))
		else
			{
			.timers[id] = cb
			args = Object(id, id, ms, false)
			// if callback is DelayBase,
			// let the browser side cancel the timer automatically
			// to avoid extra unnecessary events
			if Instance?(cb) and cb.Base?(DelayBase)
				args.once? = true
			.delayQueue.Add(Object(:id, action: Object(false, #SuSetTimer, args)))
			}
		return id
		}

	KillTimer(id)
		{
		if .timers.Member?(id)
			{
			.timers.Delete(id)
			.delayQueue.Add(Object(:id,
				action: Object(false, #SuKillTimer, Object(id, id))))
			}
		else
			.delayToCancel[id] = true
		return true
		}

	Timeout(id)
		{
		if .timers.Member?(id)
			(.timers[id])(0, 0, id, 0)
		}

	FlushDelays()
		{
		while not Same?(.delayQueue, item = .delayQueue.PopFirst())
			{
			if .delayToCancel.Member?(item.id)
				continue
			if item.Member?(#cb)
				(item.cb)(0, 0, item.id, 0)
			else
				.recordAction(item.action)
			}
		.delayToCancel = Object()
		}

	recordAction(action)
		{
		SuRenderBackend().RecordAction(@action)
		}
	}