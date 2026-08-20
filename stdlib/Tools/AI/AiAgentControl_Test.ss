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
Test_percentUsed()
		{
		percentUsed = AiAgentControl.AiAgentControl_percentUsed
		Assert(percentUsed(0, "1M") is: "0% / 1M")
		Assert(percentUsed(1, "1M") is: "0% / 1M")
		Assert(percentUsed(500000, "1M") is: "50% / 1M")
		Assert(percentUsed(999999, "1M") is: "99% / 1M")
		Assert(percentUsed(1000000, "1M") is: "100% / 1M")
		Assert(percentUsed(12345678, "1M") is: "1234% / 1M")
		Assert(percentUsed(128000, "256K") is: "50% / 256K")
		// 129280/256000*100 == 50.5% -> Int() truncates to 50 (round would be 51)
		Assert(percentUsed(129280, "256K") is: "50% / 256K")
		Assert(percentUsed(256000, "256K") is: "100% / 256K")
		}
	}
