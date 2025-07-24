// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{

	indicators: false
	Test_setIndicator()
		{
		fn = Addon_highlight_customFields.Addon_highlight_customFields_setIndicator
		matches = Object().Set_default(Object())
		.indicators = Object().Set_default(Object())

		for check in .checks.Members()
			{
			.indicators = Object().Set_default(Object())
			fn(check, .text, matches, .setIndic, 'cust')
			Assert(.indicators is: .checks[check])
			}

		Assert(matches.Members() equalsSet:
			Object('variable two', 'variable1', 'variable3'))
		Assert(matches['variable1'] is: #(#(0, 9)))
		Assert(matches['variable two'] is: #(#(22, 12)))
		Assert(matches['variable3'] is: #(#(45, 9), #(65, 9), #(84, 9)))
		Assert(matches.Member?('variableNotFound') is: false)
		}

	setIndic(indic, pos, len)
		{
		.indicators.Add(Object(:indic, :pos, :len))
		}

	text: 'variable1 lorem ipsum variable two sit dolar variable3 et cetera variable3 ' $
		'bill bob variable3'
	checks: #(
		variable1: #(#(indic: 'cust', pos: 0, len: 9)),
		'variable two': #(#(indic: 'cust', pos: 22, len: 12)),
		variable3: #(
			#(indic: 'cust', pos: 45, len: 9),
			#(indic: 'cust', pos: 65, len: 9),
			#(indic: 'cust', pos: 84, len: 9)),
		variableNotFound: #())
	}