// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_HasFieldCondition?()
		{
		has? = ChooseFiltersControl.HasFieldCondition?
		Assert(has?(#(), '') is: false)
		Assert(has?(#(), 'field') is: false)
		Assert(has?(#([]), 'field') is: false)
		Assert(has?(#([condition_field: 'field',
			field: #(value: 'val', operation: 'equals')]), 'field'))
		Assert(has?(#([condition_field: 'field',
			field: #(value: '', operation: '')]), 'field') is: false)

		filters = #(
			[condition_field: 'testField', testField: #(operation: 'equals', value: '2')],
			[condition_field: 'otherField', otherField: #(operation: '', value: '')])
		Assert(has?(filters, 'field') is: false)
		Assert(has?(filters, 'otherField') is: false)
		Assert(has?(filters, 'testField'))
		}

	Test_resetFilter()
		{
		fakeSource = class
			{
			Get()
				{ return _data }
			}

		reset = ChooseFiltersControl.ChooseFiltersControl_resetFilter
		_data = []

		reset(fakeSource)
		Assert(_data is: [])

		_data = [condition_field: 'fakeField']
		reset(fakeSource)
		Assert(_data is: [condition_field: 'fakeField',
			fakeField: #(operation: '', value: '', value2: '')])

		_data = [condition_field: 'fakeField',
			fakeField2: #(operation: 'not empty', value: '', value2: '')]
		reset(fakeSource)
		Assert(_data is: [condition_field: 'fakeField',
			fakeField: #(operation: '', value: '', value2: '')])

		_data = [fakeField2: #(operation: 'not empty', value: '', value2: '')]
		reset(fakeSource)
		Assert(_data is: [fakeField2: #(operation: 'not empty', value: '', value2: '')])

		// ensure we don't create empty record members
		_data = [condition_field: '']
		reset(fakeSource)
		Assert(_data is: [condition_field: ''])
		}

	Test_BuildWhereFromFilter()
		{
		cl = ChooseFiltersControl
			{ ChooseFiltersControl_logInvalidConditions(unused) { } }
		m = cl.BuildWhereFromFilter
		conditionFields = []
		Assert(m(Object([condition_field: 'fakeField',
			fakeField: #(operation: '', value: '', value2: '')], :conditionFields))
			is: '')
		Assert(conditionFields isSize: 0)

		Assert(m(Object([condition_field: 'fakeField',
			fakeField: #(operation: 'equals', value: 'test', value2: '')],
			:conditionFields))
			is: ' where fakeField is "test"')
		Assert(conditionFields isSize: 0)

		Assert(m(Object([condition_field: 'fakeField',
			check: false
			fakeField: #(operation: 'equals', value: 'test', value2: '')]),
			useCheckBoxes:, :conditionFields)
			is: '')
		Assert(conditionFields isSize: 0)

		Assert(m(Object([condition_field: 'fakeField'
			fakeField: #(operation: 'equals', value: 'test', value2: '')]),
			useCheckBoxes:, :conditionFields)
			is: '')
		Assert(conditionFields isSize: 0)

		Assert(m(Object([condition_field: 'fakeField',
			check: true
			fakeField: #(operation: 'equals', value: 'test', value2: '')]),
			useCheckBoxes:, :conditionFields)
			is: ' where fakeField is "test"')
		Assert(conditionFields is: #(fakeField))

		Assert(m(Object([condition_field: 'fakeField', check: true
				fakeField: #(operation: 'greater than', value: 'test', value2: '')],
			[condition_field: 'fakeField', check: true
				fakeField: #(operation: 'greater than', value: 'test1', value2: '')],
			[condition_field: 'fakeField2', check: true
				fakeField2: #(operation: 'equals', value: 'test2', value2: '')],
			),
			useCheckBoxes:, :conditionFields)
			is: ' where fakeField > "test" where fakeField > "test1"' $
				' where fakeField2 is "test2"')
		Assert(conditionFields is: #(fakeField, fakeField2))
		}

	Test_ForceValid()
		{
		cls = ChooseFiltersControl
			{
			ChooseFiltersControl_useCheckBoxes: true
			GetRows() { return _getRows }
			}

		_getRows = #()
		Assert(cls.ForceValid() is: '')

		_emptyField = Object(Field: FakeObject(Get: ''))
		_getRows = Object(class
			{
			Get() { return [] }
			Valid() { return true }
			FindControl(@unused) { _emptyField }
			})
		Assert(cls.ForceValid() is: '')

		emptyRow = class
			{
			Get() { return [check:] }
			Valid() { return false }
			FindControl(@unused) { _emptyField }
			}
		_getRows = Object(emptyRow)
		Assert(cls.ForceValid() is: 'Please select a field\n')

		invalidRow = class
			{
			Get() { return [check:] }
			Valid() { return false }
			FindControl(@unused) { Object(Field: FakeObject(Get: 'wrong field')) }
			}
		_getRows = Object(invalidRow)
		Assert(cls.ForceValid() is: 'Invalid: wrong field')

		row = class
			{
			Get() { return [check:, condition_field: 'test'] }
			Valid() { return false }
			FindControl(@unused) { return Object(Field: FakeObject(Get: 'Test')) }
			}
		_getRows = Object(row, row)
		Assert(cls.ForceValid() is: 'Invalid: Test')

		row2 = class
			{
			Get() { return [check:, condition_field: 'test2'] }
			Valid() { return true }
			FindControl(@unused) { return Object(Field: FakeObject(Get: 'Test2')) }
			}
		row3 = class
			{
			Get() { return [check:, condition_field: 'test3'] }
			Valid() { return false }
			FindControl(@unused) { return Object(Field: FakeObject(Get: 'Test3')) }
			}
		_getRows = Object(row, row, row2, row3, emptyRow, emptyRow)
		Assert(cls.ForceValid() is: 'Please select a field\nInvalid: Test, Test3')
		}
	}
