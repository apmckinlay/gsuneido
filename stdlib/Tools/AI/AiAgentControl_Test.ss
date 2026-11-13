// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_formatTime()
		{
		formatTime = AiAgentControl.AiAgentControl_formatTime
		Assert(formatTime(0) is: "")
		Assert(formatTime(1) is: "1s")
		Assert(formatTime(59) is: "59s")
		Assert(formatTime(60) is: "1m")
		Assert(formatTime(61) is: "1m 1s")
		Assert(formatTime(3599) is: "59m 59s")
		Assert(formatTime(3600) is: "60m")
		Assert(formatTime(3661) is: "61m 1s")

		Assert(formatTime(0.2) is: "")
		Assert(formatTime(0.5) is: "1s")
		Assert(formatTime(59.4) is: "59s")
		Assert(formatTime(59.5) is: "1m")
		Assert(formatTime(60.4) is: "1m")
		Assert(formatTime(60.5) is: "1m 1s")
		}
	}
