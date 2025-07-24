// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(.UniqueId)
		{
		SuRender().Register(.UniqueId, this)
		SuRender().Overlay.Show(id: #webtest, msg: 'Running Web Test...', level: 99)
		}

	RunTests(tests)
		{
		observer = TestObserverString(quiet:)
		SuTestRunner.RunList(tests, observer)
		SuRender().Event(.UniqueId, 'TestsResults', args: [observer.Result])
		}

	finished?: false
	Finish()
		{
		.finished? = true
		SuRender().Overlay.Close(id: #webtest)
		SuRender().UnRegister(.UniqueId)
		}
	}