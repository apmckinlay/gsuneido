// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(.UniqueId, msg)
		{
		SuRender().Register(.UniqueId, this)
		SuRender().Confirm(msg, false, .cancel)
		}

	cancel?: false
	cancel()
		{
		if .finished?
			return

		.cancel? = true
		SuRender().Overlay.Show(id: #taskWithPause, msg: 'Cancelling...', level: 99)
		}

	Pause()
		{
		SuRender().Event(.UniqueId, 'DonePause', args: [.cancel?])
		}

	finished?: false
	Finish()
		{
		.finished? = true
		SuRender().Overlay.Close(id: #taskWithPause)
		SuRender().Overlay.Close(id: #confirm)
		SuRender().UnRegister(.UniqueId)
		}
	}
