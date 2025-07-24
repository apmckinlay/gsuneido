// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
function (libs = false)
	{
	if libs is false
		libs = Libraries()
	t = Timer()
		{
		s = CheckLibraries(libs)
		}
	pre = "Check Libraries (" $ String(t) $ "s) - "
	return pre $ (s is "" ? "OKAY\n\n" : "FAILURES:\n" $ s $ "\n")
	}