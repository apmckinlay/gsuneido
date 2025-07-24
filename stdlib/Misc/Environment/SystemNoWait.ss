// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function (cmd)
	{
	if Sys.Linux?()
		cmd $= ' &'
	else // escape ">>" to avoid incorrectly redirecting output of start command
		cmd = 'start ' $ cmd.Replace('>>', '^>^>')
	return System(cmd)
	}