// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// based on contribution from Helmut Enck-Radana <her@paradigma-software.de>
function (alsoServer = false, exitCode = 0)
	{
	if Sys.Server?()
		{
		// give client time to exit
		Thread({ Thread.Sleep(1.SecondsInMs()); Exit(exitCode) })
		}
	else
		{
		if alsoServer is true and Sys.Client?()
			ServerEval('Shutdown', :exitCode)
		PersistentWindow.CloseSet()
		ExitClient(exitCode) // no timer on standalone or client
		}
	}