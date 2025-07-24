// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
// tracks the last 10 tests run
// and shows them after an error or slow test
// to help track order of execution issues
// 20170726 apm - disabled showing history, can enable if we run into test order issues
TestObserverString
	{
	New(.filename = 'test.log', .quiet = false)
		{
		super(quiet)
		.testLog = Object()
		.testTimeLog = Object()
		}
	BeforeTest(name)
		{
		.start = Date()
		super.BeforeTest(name)
		}
	AfterTest(name, time, dbgrowth, memory)
		{
		str = Display(.start) $ '\t ' $ name
		.logTestHistory(str)

		str $= "    time: " $ time.RoundToPrecision(2) $ " sec" $
			"    memory: " $ ReadableSize(memory) $
			"    db growth: " $ ReadableSize(dbgrowth)
		.logTestTimes(str)

		super.AfterTest(name, time, dbgrowth, memory)
		}
	Error(method, error)
		{
		super.Error(method, error)
		if method isnt ''
			.ShowTestHistory()
		}
	Warning(method, warning)
		{
		if warning.Has?('SLOWTEST')
			{
			.Output(Join(".", .Name, method))
			.Output('        INFO: ' $ warning)
			.ShowTestHistory()
			}
		super.Warning(method, warning)
		}
	ShowTestHistory()
		{
		if .testLog.Empty?()
			return
		prev = Date(.testLog.Last().BeforeFirst('\t'))
		duration = Date().MinusSeconds(prev)
		.Output('\t\tDuration of Test: ' $ duration $ ' sec')
//		.Output('\t\tLast Run Tests Were: ')
//		.Output('\t\t' $ .testLog.Join('\r\n\t\t'))
		}
	After(time, dbgrowth, memory)
		{
		super.After(time, dbgrowth, memory)
		s = 'System Tests as of ' $ Date().StdShortDateTime() $ '\r\n' $ .Result
		AddFile(.filename, s)
		PutFile('times' $ .filename, .testTimeLog.Join('\r\n'))
		return s
		}
	logTestHistory(msg)
		{
		if .testLog.Size() >= 10 /*= only store the last 10, delete oldest */
			.testLog.Delete(0)
		.testLog.Add(msg)
		}
	logTestTimes(msg)
		{
		.testTimeLog.Add(msg)
		}
	}
