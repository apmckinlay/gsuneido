// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Main()
		{
		.styles = .instance()
		suneidolog = .SpyOn(SuneidoLog).Return('')

		addons = .mockAddons()
		.styles.DefineStyles(addons)

		.test_defineStyles()
		.test_errorHandling(suneidolog.CallLogs())
		.test_markerColors()
		}

	instance()
		{
		.sci = Mock(ScintillaAddonsControl)
		.sci.When.BaseStyling([anyArgs:]).Return([
				[level: 10,
					marker: [SC.MARK_ROUNDRECT, back: CLR.CYAN]
					indicator: [INDIC.DASH, fore: CLR.LIGHTRED]],
				[level: 5, marker: [SC.MARK_ROUNDRECT, back: CLR.RED]]])
		.sci.When.DefineMarker([anyArgs:]).Do({ })
		.sci.When.DefineXPMMarker([anyArgs:]).Do({ })
		.sci.When.DefineIndicator([anyArgs:]).Do({ })
		.sci.When.IndicatorAllOnFor([anyArgs:]).Return(0)

		inst = ScintillaAddonsLineStyles(.sci)
		inst.MAXMARKERS = 8
		return inst
		}

	mockAddons()
		{
		mock = Mock(AddonManager)
		mock.When.Collect([anyArgs:]).Return([
			[
				[level: 7,
					marker: [SC.MARK_ROUNDRECT, back: CLR.GREEN],
					indicator: [INDIC.TT, fore: CLR.BLUE]],
				[level: 8,
					marker: [SC.MARK_ROUNDRECT, back: CLR.BLACK],
					indicator: [INDIC.TT, fore: CLR.BLACK]]],
			[
				[level: 8,  // Duplicate
					marker: [SC.MARK_ROUNDRECT, back: CLR.BLUE],
					indicator: [INDIC.TT, fore: CLR.BLACK]],
				[level: 6, marker: [SC.MARK_ROUNDRECT, back: CLR.DARKRED]],
				[level: 9, marker: [SC.MARK_ROUNDRECT, back: CLR.DARKGREEN]],
				[level: 10,
					marker: [SC.MARK_ROUNDRECT, back: CLR.WHITE, type: #special]],
				[level: 11,
					marker: [SC.MARK_ROUNDRECT, back: CLR.YELLOW, type: #special]]],
			[
				[level: 1, marker: [SC.MARK_ROUNDRECT, back: CLR.NONE]],
				[level: 8, marker: [SC.MARK_ROUNDRECT, back: CLR.GRAY, type: #special]],
				]])
		mock.When.Send([anyArgs:]).Do({ })
		return mock
		}

	test_defineStyles()
		{
		Assert(.styles.MarkerTypes equalsSet: #(default, special))
		markers = .styles.ScintillaAddonsLineStyles_markers
		indicators = .styles.ScintillaAddonsLineStyles_indicators

		.verifyStyling(indicators, #(7, 8, 10), #indicator)
		Assert(.styles.IndicatorIdx(7) is: 8) // First indicator is at 8
		Assert(.styles.IndicatorIdx(8) is: 9)
		Assert(.styles.IndicatorIdx(10) is: 10)

		.verifyStyling(markers.default, #(6, 7, 8, 9, 10), #marker)
		Assert(.styles.MarkerIdx(6, #default) is: 0)

		.verifyStyling(markers.special, #(8, 10, 11), #marker)
		Assert(.styles.MarkerIdx(11, #special) is: 7)
		}

	verifyStyling(styleOb, expected, styleType)
		{
		Assert(styleOb isSize: expected.Size())
		for level in expected
			Assert(styleOb.Any?({ it.level is level }),
				msg: 'failed to find ' $ styleType $ ' with level: ' $ level)
		}

	test_errorHandling(logs)
		{
		duplicates = logs[0]
		Assert(duplicates.message
			is: 'ERROR: (CAUGHT) Two markers share the same type and level. ' $
				'Skipping marker')
		Assert(duplicates.params.level is: 8)

		duplicates = logs[1]
		Assert(duplicates.message is:
			'ERROR: (CAUGHT) Two indicators share the same type and level. ' $
				'Skipping indicator')
		Assert(duplicates.params.level is: 8)

		excessMarkers = logs[2]
		Assert(excessMarkers.message is: 'ERROR: (CAUGHT) Too many markers defined: 10')
		.verifyStyling(excessMarkers.params, #(1, 5), #marker)
		}

	test_markerColors()
		{
		.verifyColor(6, #default, CLR.DARKRED)
		.verifyColor(7, #default, CLR.GREEN)
		// Not verifying [level 8, type: default], as this may vary due to the
		// duplicates and the order the addons are processed.
		// As it is technically an error to have duplicates we will not assert its color
		.verifyColor(9, #default, CLR.DARKGREEN)
		.verifyColor(10, #default, CLR.CYAN)

		.verifyColor(8, #special, CLR.GRAY)
		.verifyColor(10, #special, CLR.WHITE)
		.verifyColor(11, #special, CLR.YELLOW)
		}

	verifyColor(level, type, expected)
		{
		idx = .styles.MarkerIdx(level, type)
		Assert(.styles.MarkerColor(idx) is: expected,
			msg: 'type: ' $ type $ ', level: ' $ level)
		}
	}