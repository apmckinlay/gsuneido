// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(ms, .f)
		{
		// Assign .timerId first to handle the case where
		// .Kill() is called before the return of .SetTimer
		.timerId = SuRenderBackend().TimerManager.ReserveId()
		SuRenderBackend().TimerManager.SetTimer(ms, this, id: .timerId)
		}

	Call(@unused)
		{
		.Kill()
		(.f)()
		}

	Kill()
		{
		SuRenderBackend().TimerManager.KillTimer(.timerId)
		}
	}