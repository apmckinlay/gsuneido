// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// Using Mock because this is a control,
	// do not have a way to construct the control in a test
	Test_List()
		{
		mock = Mock(ChooseManyAsObjectControl)
		mock.ChooseManyAsObjectControl_listField = 'bedrock_residents'
		mock.When.Send('GetField', 'bedrock_residents').Return(#('Fred', 'Barney'))
		mock.When.GetList().CallThrough()

		Assert(mock.Eval(ChooseManyAsObjectControl.List) is: #('Fred', 'Barney'))
		mock.Verify.Times(1).GetList()
		mock.Verify.Times(1).Send('GetField', 'bedrock_residents')

		mock.ChooseManyAsObjectControl_listarg = #('Wilma', 'Betty')
		Assert(mock.Eval(ChooseManyAsObjectControl.List) is: #('Wilma', 'Betty'))
		// GetList should not get called the second time
		mock.Verify.Times(1).GetList()
		mock.Verify.Times(1).Send('GetField', 'bedrock_residents')
		}

	Test_NoData()
		{
		cmoData = ChooseManyAsObjectControl
			{
			List()
				{
				return #('Pebbles', 'Bambam')
				}
			}
		cmoNoData = ChooseManyAsObjectControl
			{
			List()
				{
				return #()
				}
			}
		Assert(cmoData.NoData?() is: false)
		Assert(cmoNoData.NoData?())
		}

	Test_GetListItems()
		{
		cmo = ChooseManyAsObjectControl
			{
			ChooseManyAsObjectControl_idField: 'name'
			}
		list = Object(
			Record(name: 'fred', choosemany_select: false),
			Record(name: 'wilma', choosemany_select: false),
			Record(name: 'barney', choosemany_select: false),
			Record(name: 'betty', choosemany_select: false))
		selected = #('fred', 'betty')
		Assert(cmo.GetListItems(list, selected) is: Object(
			Record(name: 'fred', choosemany_select: true),
			Record(name: 'wilma', choosemany_select: false),
			Record(name: 'barney', choosemany_select: false),
			Record(name: 'betty', choosemany_select: true)))

		}

	Test_setListOb()
		{
		setListOb = ChooseManyAsObjectControl.ChooseManyAsObjectControl_setListOb
		result = Object(
			Record(name: 'fred', choosemany_select: true),
			Record(name: 'wilma', choosemany_select: false),
			Record(name: 'barney', choosemany_select: false),
			Record(name: 'betty', choosemany_select: true))
		Assert(setListOb(result, 'name') is: #('fred', 'betty'))
		}

	Test_Set()
		{
		cmo = ChooseManyAsObjectControl
			{
			ChooseManyAsObjectControl_delimiter: ', '
			ChooseManyAsObjectControl_idField: 'name'
			ChooseManyAsObjectControl_displayField: 'desc'
			ChooseManyAsObjectControl_listarg: false
			GetList()
				{
				return #(
					[name: 'fred', desc: 'Fred Flintstone', choosemany_select: true],
					[name: 'wilma', desc: 'Wilma Flintstone', choosemany_select: false],
					[name: 'barney', desc: 'Barney Rubble', choosemany_select: false],
					[name: 'betty', desc: 'Betty Rubble', choosemany_select: true])
				}
			GetOtherList()
				{
				return #(
					[name: 'pebbles', desc: 'Pebbles Flintstone'])
				}
			}

		set = cmo.ChooseManyAsObjectControl_setField
		field = MockObject(#(
			(Set ""),
			(Set "Fred Flintstone, Betty Rubble, Pebbles Flintstone"),
			(Assert)))

		Assert(set(field, '') is: #())
		Assert(set(field, #('fred','betty','pebbles','george'))
			is: #('fred', 'betty', 'pebbles'))
		field.Assert() // ensures all calls were done
		}

	Test_ellipseIfNeeded?()
		{
		cmo = ChooseManyAsObjectControl
			{
			ChooseManyAsObjectControl_delimiter: ', '
			}
		ellipse? = cmo.ChooseManyAsObjectControl_ellipseIfNeeded?
		Assert(ellipse?('') is: '')
		Assert(ellipse?('1') is: '1')
		// limit is 1021 'bobby, ' 7 chars, should get 145 full 'bobby, '(s)
		Assert(ellipse?('bobby, '.Repeat(200)) is: 'bobby, '.Repeat(145) $ ' ...')
		}
	}