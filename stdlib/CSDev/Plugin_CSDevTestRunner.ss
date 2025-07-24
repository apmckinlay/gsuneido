// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
#(
Contributions:
	(
	(TestRunner, beforeRun, func: function ()
		{
		// don't set current directory if server
		if CSDev?() and not Server?()
			{
			prevDir = GetCurrentDirectory()
			folder = GetAppTempPath() $ 'csdevtests'
			try
				{
				EnsureDir(folder)
				SetCurrentDirectory(folder)
				}
			catch (e)
				Print('Create temp test folder failed: ' $ Display(e))
			return Object(:prevDir)
			}
		return Object()
		})
	(TestRunner, afterRun, func: function(prevSettings)
		{
		if CSDev?() and not Server?()
			{
			SetCurrentDirectory(prevSettings.prevDir)
			folder = GetAppTempPath() $ 'csdevtests'
			if DirExists?(folder)
				try
					DeleteDir(folder)
				catch (e)
					Print('DeleteDir on temp test folder failed: ' $ Display(e))
			}
		})
	)
)
