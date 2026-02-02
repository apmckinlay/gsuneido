// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_match()
		{
		m = PatternControl.PatternControl_match
		Assert(m('', '', '') is: '')
		Assert(m('', 'A#A #A#', '^[a-zA-Z][0-9][a-zA-Z][ ]?[0-9][a-zA-Z][0-9]$')
			is: false)
		Assert(m('abc', 'A#A #A#','^[a-zA-Z][0-9][a-zA-Z][ ]?[0-9][a-zA-Z][0-9]$')
			is: false)
		Assert(m('v6b 3y3', 'A#A #A#','^[a-zA-Z][0-9][a-zA-Z][ ]?[0-9][a-zA-Z][0-9]$')
			is: 'V6B 3Y3')
		Assert(m('s0l3j0', 'A#A #A#','^[a-zA-Z][0-9][a-zA-Z][ ]?[0-9][a-zA-Z][0-9]$')
			is: 'S0L 3J0')
		phoneRegex = '[0-9][0-9][0-9][-]?[0-9][0-9][0-9][0-9][x]?[0-9]'
		Assert(m('', '###-####x#', phoneRegex) is: false)
		Assert(m('abc', '###-####x#', phoneRegex) is: false)
		Assert(m('555-1212', '###-####x#', phoneRegex) is: false)
		Assert(m('555-1212ext9', '###-####x#', phoneRegex) is: false)
		Assert(m('555-1212x9', '###-####x#', phoneRegex) is: '555-1212x9')
		}

	Test_ValidData?()
		{
		Assert(PatternControl.ValidData?('a', '#') is: false)
		Assert(PatternControl.ValidData?('1', '#'))
		Assert(PatternControl.ValidData?('', '#'))
		Assert(PatternControl.ValidData?('', '#', mandatory:) is: false)
		}
	}