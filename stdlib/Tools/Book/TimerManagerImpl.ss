// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	timer: false
	New()
		{
		.setTimer()
		}
	timerFunc()
		{
		Plugins().ForeachContribution('Timer', 'timerfunction', showErrors:)
			{|c|
			LogErrors('TimerManager')
				{
				(c.func)()
				}
			}
		.setTimer()
		return 0
		}
	setTimer()
		{
		.timer = Delay(60.SecondsInMs(), .timerFunc)
		}
	Stop()
		{
		if .timer isnt false
			.timer.Kill()
		}
	}