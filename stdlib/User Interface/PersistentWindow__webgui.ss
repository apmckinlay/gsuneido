// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_PersistentWindow
	{
	loadFont(stateobject)
		{
		mLogName = GetCustomFontSaveName()
		stateobject.classstate.logfont =
			stateobject.classstate.GetDefault(mLogName, Suneido.logfont.Copy())
		}
	restoreLibs(@unused) { }
	resolution()
		{
		return 'web'
		}

	exit()
		{
		for w in SuRenderBackend().WindowManager.GetWindows().Copy().Reverse!()
			{
			if w isnt this
				w.DESTROY()
			}
		SuRenderBackend().RecordAction(false, #SuShutdown, #())
		}
	}