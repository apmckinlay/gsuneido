// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{ }

	LaunchOptions(options)
		{
		options.Add('Desktop IDE')
		}

	MenuButtons(menuBar /*unused*/, tools)
		{
		tools.Add('Code Recovery')
		}

	On_Launch_Desktop_IDE()
		{
		CSDevServerWindow.StartClient('CSDevServerWindow.StartIDE()')
		}

	On_Tools_Code_Recovery()
		{
		CodeRecoveryControl()
		}
	}