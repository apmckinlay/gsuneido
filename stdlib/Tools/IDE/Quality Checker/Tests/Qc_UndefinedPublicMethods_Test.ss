// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_calculatePubMethodCalls()
		{
		testCode = "//Suneido 3017
			class
				{
				PubMeth1()
					{
					.CallOne()
					.
					CallTwo()
					x = .CallThree()
					x = .
					CallFour()
					ClassHelp.Methods(.CallFive())
					ClassHelp.Methods(.
						CallSix())
					}
				}"
		calculatedPubMethCalls = Qc_UndefinedPublicMethods.
			Qc_UndefinedPublicMethods_calculatePubMethodCalls(testCode)
		expectedPubMethCalls = #(
			[method: "CallOne", line: 6],
			[method: "CallTwo", line: 8],
			[method: "CallThree", line: 9],
			[method: "CallFour", line: 11],
			[method: "CallFive", line: 12],
			[method: "CallSix", line: 14])
		Assert(calculatedPubMethCalls is: expectedPubMethCalls)
		}

	Test_()
		{
		child = "//Suneido 3017
			Parent
				{
				PubMeth1()
					{
					.PubMeth2()
					.PubMeth3()
					.Undef1()
					.
					Undef2()
					.PubMeth4()
					.PubMeth5()
					cl = class { }
					cl.NonDef1()
					cl.
					NonDef2()
					fn = function() {.Undef3()}
					fn = function() {.
						Undef5()}
					fn(`test`).
					NonDef3()
					}
				PubMeth2(){ }
				PubMeth3()
					{
					.Undef4()
					.Member?()
					}
				}"
		.MakeLibraryRecord([name: "Child", text: child])

		parent = "//Suneido 3017
			GrandParent
				{
				PubMeth4(){ }
				}"
		.MakeLibraryRecord([name: "Parent", text: parent])

		grandParent = "//Suneido 3017
			class
				{
				PubMeth5()
					{
					.UndefGrandParent(){ }
					}
				}"
		.MakeLibraryRecord([name: "GrandParent", text: grandParent])

		recordData = Record()
		recordData.lib = "Test_lib"
		recordData.code = child
		recordData.recordName = "Child"

		calculatedResult = Qc_UndefinedPublicMethods(recordData, true)
		expectedResult = #(warnings:
			([name: "Test_lib:Child:8 - .Undef1 is undefined"],
			[name: "Test_lib:Child:10 - .Undef2 is undefined"],
			[name: "Test_lib:Child:17 - .Undef3 is undefined"],
			[name: "Test_lib:Child:19 - .Undef5 is undefined"],
			[name: "Test_lib:Child:26 - .Undef4 is undefined"]),
			desc: "Undefined public methods were found")
		calculatedLineWarnings = calculatedResult.Extract('lineWarnings')
		Assert(calculatedResult is: expectedResult)

		expectedLineWarnings = #((8), (10), (17), (19), (26))
		Assert(calculatedLineWarnings is: expectedLineWarnings)

		result = Qc_UndefinedPublicMethods(recordData, false)
		expectedResult = #(warnings:
			([name: ".Undef1 is undefined"],
			[name: ".Undef2 is undefined"],
			[name: ".Undef3 is undefined"],
			[name: ".Undef5 is undefined"],
			[name: ".Undef4 is undefined"]),
			desc: "Undefined public methods were found")
		calculatedLineWarnings = result.Extract('lineWarnings')
		Assert(result is: expectedResult)
		Assert(calculatedLineWarnings isSize: 0)
		}

	Test_Two()
		{
		base = "//Suneido 3017
			class
				{
				PubMeth1()
					{
					.PubMeth2()
					.PubMeth3()
					.Undef1()
					.
					Undef2()
					.PubMeth4()
					.PubMeth5()
					}
				PubMeth2(){ }
				PubMeth3()
					{
					.Undef3()
					}
				Default(){}
				}"
		.MakeLibraryRecord([name: "Base", text: base])

		recordData = Record()
		recordData.lib = "Test_lib"
		recordData.code = base
		recordData.recordName = "Base"
		calculatedResult = Qc_UndefinedPublicMethods(recordData, true)
		expectedResult = #(warnings: (), desc: "", lineWarnings: #())
		Assert(calculatedResult is: expectedResult)

		allResults = Qc_UndefinedPublicMethods(recordData, false)
		expectedResults = #(warnings: (), desc: "No undefined public methods were found",
			lineWarnings: #())
		Assert(allResults is: expectedResults)

		recordData.lib = "FAKELIB"
		calculatedResult = Qc_UndefinedPublicMethods(recordData, true)
		Assert(calculatedResult is: expectedResult)

		allResults = Qc_UndefinedPublicMethods(recordData, false)
		expected = #(warnings: (), desc: "Undefined public method checking aborted",
			lineWarnings: #())
		Assert(allResults is: expected)
		}

	Test_noSocketServerBuiltIn()
		{
		text = "//Suneido 3017
			SocketServer
				{
				New()
					{
					.Writeline()
					.Read()
					.Readline()
					.Write()
					.RemoteUser()
					.NotDefined()
					}
				}"
		.MakeLibraryRecord([name: 'TestSocket', :text])
		recordData = Record(lib: 'Test_lib', code: text, recordName: 'TestSocket')
		calculatedResult = Qc_UndefinedPublicMethods(recordData, true)
		//Should ignore SocketServer built-in methods, but still catch other undefined
		expectedResult = #(warnings:
			([name: "Test_lib:TestSocket:11 - .NotDefined is undefined"]),
				desc: 'Undefined public methods were found', lineWarnings: ((11)))
		Assert(calculatedResult is: expectedResult)

		result = Qc_UndefinedPublicMethods(recordData, false)
		expectedResult = #(warnings: ([name: ".NotDefined is undefined"]),
			desc: 'Undefined public methods were found', lineWarnings: ())
		Assert(result is: expectedResult)
		}

	Test_noBuiltInString()
		{
		text = "//Suneido 3017
			class
				{
				New()
					{
					s1 = `short text, shouldn't have an issue`.Tr('\r\n')
					s2 = `loooooooong text,
						with line breaks,
						shouldn't have an issue`.Tr('\r\n')
					s3 = `another loooooooong text,
						with line breaks,
						shouldn't have an issue`.Split('\r\n')
					s4 = `another loooooooong text,
						with line breaks,
						and a fake method,
						still shouldn't have an issue`.NotRealFn('\r\n')
					}
				}"
		.MakeLibraryRecord([name: 'TestStringMethods', :text])
		recordData = Record(lib: 'Test_lib', code: text, recordName: 'TestStringMethods')
		calculatedResult = Qc_UndefinedPublicMethods(recordData, true)
		expectedResult = #(warnings: (), desc: "", lineWarnings: #())
		Assert(calculatedResult is: expectedResult)

		result = Qc_UndefinedPublicMethods(recordData, false)
		expectedResult = #(warnings: (), desc: "No undefined public methods were found",
			lineWarnings: #())
		Assert(result is: expectedResult)
		}
	}