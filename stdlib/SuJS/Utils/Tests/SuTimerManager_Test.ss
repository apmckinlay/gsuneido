// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_delays()
		{
		mock = Mock(SuTimerManager)
		mock.SuTimerManager_timerId = 1
		mock.SuTimerManager_timers = Object()
		mock.SuTimerManager_delayQueue = Object()
		mock.SuTimerManager_delayToCancel = Object()
		mock.When.SetTimer([anyArgs:]).CallThrough()
		mock.When.KillTimer([anyArgs:]).CallThrough()
		mock.When.FlushDelays([anyArgs:]).CallThrough()
		mock.When.Timeout([anyArgs:]).CallThrough()
		mock.When.recordAction([anyArgs:]).Do({ |@unused| })
		cb = { |@unused| mock.SetTimer(2000, false) }

		Assert(mock.SetTimer(0, cb) is: 1)
		Assert(mock.SuTimerManager_delayQueue is: Object(Object(id: 1, :cb)))

		Assert(mock.SetTimer(1000, cb) is: 2)
		Assert(mock.SuTimerManager_delayQueue
			is: Object(
				Object(id: 1, :cb),
				Object(id: 2, action: #(false, #SuSetTimer, (2, 2, 1000, false)))))
		Assert(mock.SuTimerManager_timers is: Object(2: cb))
		Assert(mock.SuTimerManager_delayToCancel is: #())

		Assert(mock.SetTimer(0, cb) is: 3)
		Assert(mock.SetTimer(1001, cb) is: 4)
		Assert(mock.SuTimerManager_delayQueue
			is: Object(
				Object(id: 1, :cb),
				Object(id: 2, action: #(false, #SuSetTimer, (2, 2, 1000, false))),
				Object(id: 3, :cb),
				Object(id: 4, action: #(false, #SuSetTimer, (4, 4, 1001, false)))))
		Assert(mock.SuTimerManager_timers is: Object(2: cb, 4: cb))
		Assert(mock.SuTimerManager_delayToCancel is: #())

		Assert(mock.KillTimer(1))
		Assert(mock.SuTimerManager_delayToCancel is: #(1:))

		Assert(mock.KillTimer(4))
		Assert(mock.SuTimerManager_delayToCancel is: #(1:))
		Assert(mock.SuTimerManager_delayQueue
			is: Object(
				Object(id: 1, :cb),
				Object(id: 2, action: #(false, #SuSetTimer, (2, 2, 1000, false))),
				Object(id: 3, :cb),
				Object(id: 4, action: #(false, #SuSetTimer, (4, 4, 1001, false))),
				Object(id: 4, action: #(false, #SuKillTimer, (4, 4)))))
		Assert(mock.SuTimerManager_timers is: Object(2: cb))

		mock.FlushDelays()
		// delay 2
		mock.Verify.Times(1).recordAction([false, #SuSetTimer, #(2, 2, 1000, false)])
		// delay 4
		mock.Verify.Times(1).recordAction([false, #SuSetTimer, #(4, 4, 1001, false)])
		// kill delay 4
		mock.Verify.Times(1).recordAction([false, #SuKillTimer, #(4, 4)])
		// delay 5 created by delay 1
		mock.Verify.Times(1).recordAction([false, #SuSetTimer, #(5, 5, 2000, false)])
		mock.Verify.Times(4).recordAction([anyArgs:])

		Assert(mock.SuTimerManager_delayQueue is: #())
		Assert(mock.SuTimerManager_timers is: Object(2: cb, 5: false))
		Assert(mock.SuTimerManager_delayToCancel is: #())

		mock.Timeout(2)
		Assert(mock.SuTimerManager_timers is: Object(2: cb, 5: false, 6: false))
		}
	}