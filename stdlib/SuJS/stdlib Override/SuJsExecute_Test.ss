// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_open()
		{
		spy = .SpyOn(SuUI.Open).Return('')

		SuJsExecute('open', 'http://test.com')
		Assert(spy.CallLogs().Last().url is: 'http://test.com')

		SuJsExecute('open', 'mailto:test@test.com')
		Assert(spy.CallLogs().Last().url is: 'mailto:test@test.com')

		SuJsExecute('open', 'axis:t:xxxxxxx')
		Assert(spy.CallLogs().Last().url is: 'axis:t:xxxxxxx')

		SuJsExecute('open', 'test.com')
		Assert(spy.CallLogs().Last().url is: 'https://test.com')
		}
	}