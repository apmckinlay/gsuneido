// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()

		// Test stdlib references
		.MakeLibraryRecord(
			[name: .recA = 'AA_' $ .TempName(), text: `function () { }`],
			[name: .recB = 'BB_' $ .TempName(), text: `function () { }`],
			[name: .recC = 'CC_' $ .TempName(), text: `function () { }`],
			table: 'stdlib')
		.AddTeardown({
			QueryDo('delete stdlib where group is -1 and
				name in ("' $ .recA $ '", "' $ .recB $ '", "' $ .recC $ '")')
			})

		// Test_lib
		.MakeLibraryRecord(
			[name: .recD = 'DD_' $ .TempName(), text: `function () { }`],
			[name: .recE = 'EE_' $ .TempName(), text: `function () { }`],
			[name: .recF = 'FF_' $ .TempName(), text: `function () { }`],
			[name: .recG = 'GG_' $ .TempName(), text: `function () { }`],
			[name: .recH = 'HH_' $ .TempName(), text: `function () { }`],
			[name: .recI = 'II_' $ .TempName(), text: `function () { }`],
			[name: .recJ = 'JJ_' $ .TempName(), text: `function () { }`],
			[name: .recK = 'KK_' $ .TempName(), text: `function () { }`],
			[name: .recL = 'LL_' $ .TempName(), text: `function () { }`],
			[name: .recM = 'MM_' $ .TempName(), text: `function () { }`],
			[name: .recN = 'NN_' $ .TempName(), text: `function () { }`],
			[name: .recO = 'OO_' $ .TempName(), text: `function () { }`],
			[name: .recP = 'PP_' $ .TempName(), text: `function () { }`],
			[name: .recQ = 'QQ_' $ .TempName(), text: `function () { }`],
			[name: .recR = 'RR_' $ .TempName(), text: `function () { }`],
			[name: .recS = 'SS_' $ .TempName(), text: `function () { }`],
			[name: .recT = 'TT_' $ .TempName(), text: `function () { }`],
			[name: .recU = 'UU_' $ .TempName(), text: `function () { }`],
			[name: .recV = 'VV_' $ .TempName(), text: `function () { }`],
			[name: .recW = 'WW_' $ .TempName(), text: `function () { }`],
			[name: .recX = 'XX_' $ .TempName(), text: `function () { }`],
			[name: .recY = 'YY_' $ .TempName(), text: `function () { }`],
			[name: .recZ = 'ZZ_' $ .TempName(), text: `function () { }`])

		Qc_stdNames.ResetCache()
		Qc_whichLib.ResetCache()
		LibraryTables.ResetCache()
		}

	Test_main()
		{
		.spyOn()
		code = `class
			{
			CallClass()
				{
				` $ .recA $ `()
				` $ .recS $ `()
				` $ .recB $ `()
				` $ .recC $ `()
				.private()
				` $ .recA $ `()
				` $ .recU $ `()
				` $ .recW $ `()
				` $ .recT $ `()
				}

			private()
				{
				` $ .recV $ `()
				` $ .recZ $ `()
				` $ .recQ $ `()
				if true
					` $ .recN $ `()
				else
					` $ .recM $ `()
				` $ .recX $ `()
				` $ .recY $ `()
				}

			Public(ob)
				{
				` $ .recR $ `()
				` $ .recI $ `()
				` $ .recK $ `()
				` $ .recL $ `()
				ob.Each()
					{
					` $ .recO $ `()
					` $ .recP $ `()
					` $ .recQ $ `()
					if true
						{
						` $ .recD $ `()
						` $ .recE $ `()
						}
					else
						` $ .recH $ `()
					}
				` $ .recJ $ `()
				` $ .recF $ `()
				` $ .recJ $ `()
				` $ .recG $ `()
				` $ .recJ $ `()
				` $ .recM $ `()
				` $ .recJ $ `()
				}
			}`

		stdlibDependencies = 'stdlib: ' $ .recA $ '(2), ' $ .recB $ '(1), ' $
			.recC $ '(1)'
		testTableDependencies = 'Test_lib: ' $
			.recD $ '(1), ' $ .recE $ '(1), ' $ .recF $ '(1), ' $ .recG $ '(1), ' $
			.recH $ '(1), ' $ .recI $ '(1), ' $ .recJ $ '(4), ' $ .recK $ '(1), ' $
			.recL $ '(1), ' $ .recM $ '(2), ' $ .recN $ '(1), ' $ .recO $ '(1), ' $
			.recP $ '(1), ' $ .recQ $ '(2), ' $ .recR $ '(1), ' $ .recS $ '(1), ' $
			.recT $ '(1), ' $ .recU $ '(1), ' $ .recV $ '(1), ' $ .recW $ '(1), ' $
			.recX $ '(1), ' $ .recY $ '(1), ' $ .recZ $ '(1)'

		.assertCoupling(code, Object(
				warnings: Object([
					name: 'Depends on:\n\t' $
						stdlibDependencies $ '\n\t' $
						testTableDependencies
						]),
				desc: 'This record depends on 26 records - Attempt to limit to 20',
				rating: 0, size: 6),
			msg: 'main 1',
			minimizeOutput?:)

		.assertCoupling(code, Object(
				warnings: Object([
					name: 'Depends on:\n\t' $
						stdlibDependencies $ '\n\t' $
						testTableDependencies
						]),
				desc: 'This record depends on 26 records - Attempt to limit to 20',
				rating: 0, size: 6),
			msg: 'main 2')

		.assertCoupling(code, Object(
				warnings: Object([name: 'Depends on:\n\t' $ testTableDependencies]),
				desc: 'This record depends on 23 records - Attempt to limit to 20',
				rating: 2, size: 3),
			msg: 'main 3',
			lib: 'Test_lib')

		.assertCoupling(code, Object(
				warnings: Object([name: 'Depends on:\n\t' $ testTableDependencies]),
				desc: 'This record depends on 23 records - Attempt to limit to 20',
				rating: 2, size: 3),
			msg: 'main 4',
			lib: 'Test_lib',
			minimizeOutput?:)
		}

	spyOn()
		{
		.SpyOn(Qc_Coupling.Qc_Coupling_applicationLibraries).
			Return(GetContributions('ApplicationLibraries').Add('Test_lib'))
		}

	assertCoupling(code, expected, msg, lib = 'stdlib', minimizeOutput? = false)
		{
		result = Qc_Coupling([recordName: 'Qc_Coupling', :code, :lib], :minimizeOutput?)
		Assert(result.size is: expected.size, :msg)
		Assert(result.rating is: expected.rating, :msg)
		Assert(result.desc is: expected.desc, :msg)

		// Can't compare warnings directly as library usage can affect the final result
		expectedLines = expected.warnings.GetDefault(0, []).name.Split('\n\t')
		resultLines = result.warnings.GetDefault(0, []).name.Split('\n\t')
		Assert(resultLines, isSize: expectedLines.Size(), :msg)
		resultLines.Each({ Assert(expectedLines has: it, :msg) })
		}

	Test_unknown()
		{
		.spyOn()
		code = `function ()
			{
			` $ .recA $ `()
			` $ .recZ $ `()
			TestFunction1()
			TestFunction2()
			TestFunction3()
			` $ .recZ $ `()
			}`
		.assertCoupling(code,
			#(warnings: (), desc: '', rating: 5, size: 0),
			msg: 'unkown 1',
			minimizeOutput?:)

		.assertCoupling(code, Object(
				warnings: Object(
					[name: 'Depends on:\n\tstdlib: ' $ .recA $ '(1)\n\t' $
						'Test_lib: ' $ .recZ $ '(2)\n\t' $
						'Unknown: TestFunction1(1), TestFunction2(1), TestFunction3(1)']),
				desc: 'This record depends on 5 records - Attempt to limit to 20',
				rating: 5, size: 0),
			msg: 'unkown 2')

		.assertCoupling(code, Object(
				warnings: Object(
					[name: 'Depends on:\n\tTest_lib: ' $ .recZ $ '(2)\n\t' $
						'Unknown: TestFunction1(1), TestFunction2(1), TestFunction3(1)']),
				desc: 'This record depends on 4 records - Attempt to limit to 20'
				rating: 5, size: 0),
			msg: 'unkown 3',
			lib: 'Test_lib')
		}

	Test_empty()
		{
		code = `function ()
			{
			}`
		.assertCoupling(code,
			#(warnings: (), desc: '', rating: 5, size: 0),
			msg: 'empty 1'
			minimizeOutput?:)

		.assertCoupling(code, #(
				warnings: ([name: 'Depends on:']),
				desc: 'This record depends on 0 records - Attempt to limit to 20',
				rating: 5, size: 0),
			msg: 'empty 2')

		.assertCoupling(code, #(
				warnings: ([name: 'Depends on:']),
				desc: 'This record depends on 0 records - Attempt to limit to 20',
				rating: 5, size: 0),
			msg: 'empty 3',
			lib: 'fakelib')
		}

	Teardown()
		{
		super.Teardown()
		Qc_stdNames.ResetCache()
		Qc_whichLib.ResetCache()
		LibraryTables.ResetCache()
		}
	}