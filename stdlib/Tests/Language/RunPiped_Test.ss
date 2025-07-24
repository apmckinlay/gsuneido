// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		if Sys.Windows?()
			{
			result = RunPipedOutput("net")
			Assert(result hasPrefix: "The syntax of this command is:",
				msg: result)
			}
		else
			{
			result = RunPipedOutput("mv")
			Assert(result.Prefix?("usage") or result.Prefix?("mv: missing file operand"),
				msg: result)
			}
		}
	}
