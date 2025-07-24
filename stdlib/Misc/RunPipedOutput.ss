// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(cmd, input = "")
		{
		return .WithExitValue(cmd, input).output
		}

	WithExitValue(cmd, input = "")
		{
		RunPiped(cmd)
			{|p|
			p.Write(input)
			p.CloseWrite() // otherwise some programs don't exit
			if false is output = p.Read()
				output = ""
			exitValue = p.ExitValue()
			return Object(:output, :exitValue)
			}
		}
	}