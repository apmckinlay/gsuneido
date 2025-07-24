// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_matchPrefix()
		{
		cl = ChooseListControl
			{
			Field: class { Get() { return '' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false, Record(), '') is: false)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'test' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(true, Record(), '') is: false)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'test' } }
			}
		Assert(
			cl.ChooseListControl_matchPrefix(false, Record('test1', 'test', 'test2'), '')
			is: 1)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'test' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('abc', 'testing1', 'testing2'), '') is: 1)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'test' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('abc - first', 'testing1 - second', 'test - third'), ' - ') is: 2)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'test' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('abc - first', 'abc1 - second', 'abc2 - third'), ' - ') is: false)

		cl = ChooseListControl
			{
			Field: class { Get() { return '2008' } }
			}
		Assert(
			cl.ChooseListControl_matchPrefix(false, Record(2006, 2007, 2008, 2009), ' - ')
			is: 2)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'Rate' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('first', 'RATE', 'Rate', 'other'), '') is: 2)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'Rat' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('first', 'RATE', 'Rate', 'other'), '') is: 2)

		cl = ChooseListControl
			{
			Field: class { Get() { return 'RAT' } }
			}
		Assert(cl.ChooseListControl_matchPrefix(false,
			Record('first', 'RATE', 'Rate', 'other'), '') is: 1)
		}

	Test_ValidData()
		{
		validData = ChooseListControl.ValidData?
		args = Object('', list: #('one', 'two'))
		Assert(validData(@args))

		args.mandatory = false
		Assert(validData(@args))

		args.mandatory = true
		Assert(validData(@args) is: false)

		args[0] = 'one'
		Assert(validData(@args))

		args[0] = 'FRED'
		Assert(validData(@args) is: false)

		args.allowOther = true
		Assert(validData(@args))

		args = Object('three', #('one', 'two'), mandatory:)
		Assert(validData(@args) is: false)
		args[0] = 'two'
		Assert(validData(@args))

		args = Object('one', list: #('one', 'two'), listSeparator: '')
		Assert(validData(@args))
		args[0] = 'FRED'
		Assert(validData(@args) is: false)

		args = Object('four', listField: 'theList', record: #(theList: #('one', 'two')))
		Assert(validData(@args) is: false)
		args[0] = 'two'
		Assert(validData(@args))

		args = Object('four', listField: 'theList', record: #(theList: #('one', 'two')),
			otherListOptions: #('five', 'six'))
		Assert(validData(@args) is: false)
		args[0] = 'five'
		Assert(validData(@args))
		args[0] = 'six'
		Assert(validData(@args))
		args[0] = 'seven'
		Assert(validData(@args) is: false)

		args = Object('four', listField: 'theList', splitValue: '|',
			record: #(theList: 'one|two|three'))
		Assert(validData(@args) is: false)
		args[0] = 'two'
		Assert(validData(@args))

		args = Object('one', list: #('one', 'two'), otherListOptions: #(three four))
		Assert(validData(@args))
		args[0] = 'three'
		Assert(validData(@args))
		args[0] = 'five'
		Assert(validData(@args) is: false)
		}
	}
