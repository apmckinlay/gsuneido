// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
TestObserverOnServer
	{
	AfterTest(name, time, dbgrowth, memory)
		{
		origData = .Data.Copy()
		super.AfterTest(name, time, dbgrowth, memory)

		.Data = origData
		.Data[name] = [errors: .Errors, warnings: .Warnings,
			:time, :dbgrowth, :memory, nwarnings: .Nwarnings]
		.Totals.time += time
		}

	HasError?()
		{
		return .Totals.n_failures > 0
		}

	DisplayValue(member)
		{
		onlyExpectedTests = .Data.Copy().RemoveIf({ it[member] is '' })
		resultString = ''
		for expectedTest in onlyExpectedTests.Members()
			{
			valueLines = onlyExpectedTests[expectedTest][member].Replace('\n', '\r\n\t')
			resultString $= expectedTest $ '\r\n\t' $ valueLines $ '\t' $
				'Duration of Test: ' $  onlyExpectedTests[expectedTest].time $ ' sec\r\n'
			}
		return resultString
		}

	FinalResults()
		{
		time = .Totals.time.Round(1)
		return .Totals.n_failures isnt 0
			? .Totals.n_failures $ ' tests FAILED ' $ time $ ' sec\r\n'
			: .Totals.n_tests $ ' tests SUCCEEDED ' $ time $ ' sec\r\n'
		}
	}