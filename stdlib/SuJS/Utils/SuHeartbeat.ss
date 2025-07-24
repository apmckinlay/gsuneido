// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	timer: false
	startTime: false
	currentEvent: false
	sum: 0
	i: 0
	n: 60
	New()
		{
		.last60 = Object()
		}

	Start()
		{
		if .timer isnt false
			return

		.timer = SuSetTimer(0, 0, 1.SecondsInMs(), .heartbeat)
		}

	Stop()
		{
		if .timer is false
			return

		SuKillTimer(0, .timer)
		.timer = false
		}

	heartbeat(@unused)
		{
		// if there are outstanding events, just send heartbeat without calculating latency
		// because the server is busy
		if SuRender().HasOutstandingEvents?()
			SuRender().Heartbeat()
		else
			{
			.startTime = SuUI.GetCurrentWindow().performance.now()
			.currentEvent = SuRender().Heartbeat()
			}
		}

	Event(eventId)
		{
		if .currentEvent isnt eventId
			return

		.add()
		if .logLatency?()
			{
			params = [average: (.sum / .n).Round(0),
				period: (.startTime - .last60[.i].startTime).Round(0),
				max: .last60.MaxWith({ it.latency }).latency]
			SuRender().Event(false, 'BookLog',
				Object('Detect average latency > 1 sec in last 60 samples',
					:params, systemLog:))
			.last60 = Object()
			.i = .sum = 0
			}
		.startTime = .currentEvent = false
		}

	add()
		{
		endTime = SuUI.GetCurrentWindow().performance.now()
		latency = endTime - .startTime
		if .last60.Size() < .n
			.sum += latency
		else
			.sum += latency - .last60[.i].latency
		.last60[.i] = [:latency, startTime: .startTime]
		.i = ++.i % .n
		}

	logLatency?()
		{
		return .last60.Size() is .n and .sum > 1000 * .n/*=1 sec*.n*/
		}
	}
