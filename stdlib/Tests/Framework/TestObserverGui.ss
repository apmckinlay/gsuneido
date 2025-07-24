// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
TestObserver
	{
	New()
		{
		.Totals = Object().Set_default(0)
		}
	BeforeTest(name/*unused*/)
		{
		.Errors = .Warnings = .Nwarnings = ""
		}
	Error(method, error)
		{
		.Errors $= 'ERROR ' $ method $ ': ' $ error $ '\n'
		}
	Warning(method, warning)
		{
		.Warnings $= 'WARNING ' $ method $ ': ' $ warning $ '\n'
		++.Nwarnings
		}
	AfterTest(name, time, dbgrowth, memory)
		{
		.Data = [:name, errors: .Errors, warnings: .Warnings,
			:time, :dbgrowth, :memory, nwarnings: .Nwarnings]
		++.Totals.n_tests
		if .Errors isnt ""
			++.Totals.n_failures
		.Totals.dbgrowth += dbgrowth
		.Totals.memory += memory
		.Totals.nwarnings += .Nwarnings
		}
	}