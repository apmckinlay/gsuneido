// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_updateLibView?()
		{
		m = Addon_check_code.Addon_check_code_updateLibView?

		warnings = Object(
			qc: curQc = Object(rating: false, warningText: ''),
			lines: curLines = Object())
		prevWarnings = Object(
			qc: prevQc = Object(rating: false, warningText: ''),
			lines: prevLines = Object())
		Assert(m(warnings, prevWarnings) is: false)

		curQc.rating = 5
		Assert(m(warnings, prevWarnings))

		prevQc.rating = 5
		Assert(m(warnings, prevWarnings) is: false)

		prevQc.warningText = 'TEST'
		Assert(m(warnings, prevWarnings))

		curQc.warningText = 'TEST'
		Assert(m(warnings, prevWarnings) is: false)

		curLines.Add(#('fakeLine 1'), #('fakeLine 2'))
		prevLines.Add(#('fakeLine 2'))
		Assert(m(warnings, prevWarnings))

		prevLines.Add(#('fakeLine 1'))
		Assert(m(warnings, prevWarnings) is: false)
		}

	Test_createCheckCodeText()
		{
		m = Addon_check_code.Addon_check_code_createCheckCodeText

		warnings = Object()
		Assert(m(warnings, 'TestName', 'TestLib') is: '')

		warnings.Add([line: 5, msg: 'WARNING: initialized but not used: var1'])
		Assert(m(warnings, 'TestName', 'TestLib')
			is: '\nTestLib:TestName:6 WARNING: initialized but not used: var1')

		warnings.Add([line: 10, msg: 'SHOULD NOT SEE', noOutput:])
		warnings.Add([line: 12, msg: 'ERROR: used but not initialized: var2'])
		warnings.Add([line: 12, msg: 'ERROR: useless expression'])
		Assert(m(warnings, 'TestName', 'TestLib')
			is: '\nTestLib:TestName:6 WARNING: initialized but not used: var1' $
				'\nTestLib:TestName:13 ERROR: used but not initialized: var2' $
				'\nTestLib:TestName:13 ERROR: useless expression')
		}
	}
