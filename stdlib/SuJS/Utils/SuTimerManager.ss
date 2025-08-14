// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{
		.timerId = 1
		.mutex = Mutex()
		.timers = Object()
		.delayQueue = Object()
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
			.addDelay(:id, :cb)
		else
			{
			.timers[id] = cb
			args = Object(id, id, ms, false)
			// if callback is DelayBase,
			// let the browser side cancel the timer automatically
			// to avoid extra unnecessary events
			if Instance?(cb) and cb.Base?(DelayBase)
				args.once? = true
			.addDelay(:id, action: Object(false, #SuSetTimer, args))
			}
		return id
		}

	KillTimer(id)
		{
		removed? = .removeDelayById(id)
		if .timers.Member?(id)
			{
			.timers.Delete(id)
			if not removed?
				.addDelay(:id, action: Object(false, #SuKillTimer, Object(id, id)))
			}
		return true
		}

	addDelay(@delay)
		{
		.mutex.Do({ .delayQueue.Add(delay) })
		}

	removeDelayById(id)
		{
		removed? = false
		.mutex.Do()
			{
			.delayQueue.RemoveIf()
				{
				find? = it.id is id
				if find?
					removed? = true
				find?
				}
			}
		return removed?
		}

	popDelay()
		{
		return .mutex.Do({ .delayQueue.PopFirst() })
		}

	Timeout(id)
		{
		if .timers.Member?(id)
			(.timers[id])(0, 0, id, 0)
		}

	FlushDelays()
		{
		while not Same?(.delayQueue, item = .popDelay())
			{
			if item.Member?(#cb)
				(item.cb)(0, 0, item.id, 0)
			else
				.recordAction(item.action)
			}
		}

	recordAction(action)
		{
		SuRenderBackend().RecordAction(@action)
		}
	}