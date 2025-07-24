// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		paramsselect = ParamsSelectControl
			{
			ParamsSelectControl_ops_org: ("greater than", "greater than or equal to",
				"less than", "less than or equal to", "equals", "not equal to", "empty",
				"not empty", "range", "in list", "not in list")
			ParamsSelectControl_ops_trans: ("greater than trans",
				"greater than or equal to  trans",
				"less than trans", "less than or equal to trans", "equals trans",
				"not equal to trans", "empty trans",
				"not empty trans", "range trans", "in list trans", "not in list trans")
			}
		method = paramsselect.ParamsSelectControl_find_operation
		Assert(method(#(operation: "")) is: "")
		Assert(method(#(operation: "equals")) is: "equals trans")
		Assert(method(#(operation: "in list")) is: "in list trans")
		Assert(method(#(operation: "invalid")) is: "")
		}

	Test_IsStringControl?()
		{
		method = ParamsSelectControl.IsStringControl?
		Assert(method(#()))
		Assert(method(#('Field', prompt: 'My Field', width: 10)))
		Assert(method(#(1234, prompt: 'Invalid Control Name', width: 10)))
		Assert(method(#('Number', prompt: 'My Number', width: 10)) is: false)
		Assert(method(#('Dollar', prompt: 'My Dollar', width: 10)) is: false)
		Assert(method(#('CheckBox', prompt: 'My CheckBox')) is: false)
		Assert(method(#('ChooseDate', prompt: 'ChooseDate', width: 10)) is: false)
		Assert(method(#('ChooseDateTime', prompt: 'ChooseDateTime', width: 10)) is: false)
		Assert(method(#('Date', prompt: 'Date', width: 10)) is: false)
		}

	Test_validateEmptyOperation()
		{
		m = ParamsSelectControl.ParamsSelectControl_validateEmptyOperation
		ctrl = Mock()
		ctrl.ParamsSelectControl_emptyValue = 'all'
		ctrl.Name = 'test'
		ctrl.Controller = class { }
		ctrl.When.Send([anyArgs:]).Return(0)
		Assert(ctrl.Eval(m))

		ctrl.When.Send([anyArgs:]).Return(#(operation: ''))
		ctrl.ParamsSelectControl_op = class { Get() { return '' }  }
		Assert(ctrl.Eval(m))

		ctrl.When.Send([anyArgs:]).Return(#(operation: 'contains'))
		ctrl.ParamsSelectControl_op = class { Get() { return '' }  }
		Assert(ctrl.Eval(m) is: false)

		ctrl.When.Send([anyArgs:]).Return(#(operation: 'all'))
		ctrl.ParamsSelectControl_op = class { Get() { return '' }  }
		Assert(ctrl.Eval(m))

		ctrl.Controller = ListEditWindow { }
		ctrl.ParamsSelectControl_emptyValue = 'none'
		ctrl.ParamsSelectControl_op = class { Get() { return '' }  }
		Assert(ctrl.Eval(m))
		}

	Test_handleSpecialOp()
		{
		mock = Mock(ParamsSelectControl)
		mock.When.changeControl([anyArgs:]).Return(0)
		mock.When.handleSpecialOp([anyArgs:]).CallThrough()
		mock.When.handleInList([anyArgs:]).CallThrough()
		mock.When.handleFieldControlOp([anyArgs:]).CallThrough()
		mock.When.restoreControl([anyArgs:]).CallThrough()
		mock.ParamsSelectControl_field = 'test'
		mock.ParamsSelectControl_control = Object('OriginalControl')

		/*
			Field is a string and its control isn't a ChooseControl
		*/
		mock.ParamsSelectControl_string? = true
		mock.ParamsSelectControl_useFieldForCompare? = true
		// op changes to 'equals', no need to change value control
		mock.handleSpecialOp('equals')
		mock.Verify.Never().changeControl([anyArgs:])

		// op changes to 'greater than or equal to'
		// value control should be changed to FieldControl
		mock.handleSpecialOp('greater than or equal to')
		mock.Verify.changeControl(Object('Field', width: 20, mandatory:, xstretch: 1))
		mock.Verify.Times(1).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'fieldOp')

		// op changes to 'matches'
		// value control should remain as FieldControl
		mock.handleSpecialOp('matches')
		mock.Verify.Times(1).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'fieldOp')

		// op changes to 'not in list'
		// value control should be changed to ParamsChooseList
		mock.handleSpecialOp('not in list')
		mock.Verify.Times(1).changeControl(Object('ParamsChooseList', 'test',
			values: Object(), readonly: false))
		mock.Verify.Times(2).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'inList')

		// op changes to 'in list'
		// value control should remain as ParamsChooseList
		mock.handleSpecialOp('in list')
		mock.Verify.Times(2).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'inList')

		// op changes to 'empty'
		// value control should be changed to the original control
		mock.handleSpecialOp('empty')
		mock.Verify.Times(1).changeControl(Object('OriginalControl'))
		mock.Verify.Times(3).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: '')

		/*
			Field is a string and its control is a ChooseControl
		*/
		mock.ParamsSelectControl_string? = true
		mock.ParamsSelectControl_useFieldForCompare? = false
		// op changes to 'equals', no need to change value control
		mock.handleSpecialOp('equals')
		mock.Verify.Times(3).changeControl([anyArgs:])

		// op changes to 'greater than or equal to'
		// no need to change value control
		mock.handleSpecialOp('greater than or equal to')
		mock.Verify.Times(3).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: '')

		// op changes to 'matches'
		// value control should be changed to FieldControl
		mock.handleSpecialOp('matches')
		mock.Verify.Times(2).changeControl(
			Object('Field', width: 20, mandatory:, xstretch: 1))
		mock.Verify.Times(4).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'fieldOp')

		// op changes to 'not in list'
		// value control should be changed to ParamsChooseList
		mock.handleSpecialOp('not in list')
		mock.Verify.Times(2).changeControl(Object('ParamsChooseList', 'test',
			values: Object(), readonly: false))
		mock.Verify.Times(5).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'inList')

		// op changes to 'in list'
		// value control should remain as ParamsChooseList
		mock.handleSpecialOp('in list')
		mock.Verify.Times(5).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: 'inList')

		// op changes to 'empty'
		// value control should be changed to the original control
		mock.handleSpecialOp('empty')
		mock.Verify.Times(2).changeControl(Object('OriginalControl'))
		mock.Verify.Times(6).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: '')

		/* Field is not a string field */
		mock.ParamsSelectControl_string? = false
		mock.ParamsSelectControl_useFieldForCompare? = false
		// value control should not be changed to FieldControl
		mock.handleSpecialOp('greater than')
		mock.Verify.Times(6).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: '')

		mock.handleSpecialOp('starts with')
		mock.Verify.Times(6).changeControl([anyArgs:])
		Assert(mock.ParamsSelectControl_valType is: '')
		}

	Test_dateString()
		{
		fn = .dateStringTester(false)
		date = _dateForTest = Date()
		dateNoTime = date.NoTime()
		noTimeFormat = ' ' $ dateNoTime.ShortDate()
		dateWithTime = date
		timeFormat = ' ' $ dateWithTime.ShortDateTime()

		_opVal = 'equals'
		Assert(fn() is: ' ', msg: 'equals empty space')
		_date1 = 't'
		Assert(fn() is: noTimeFormat, msg: 't with no time')

		for opVal in #('empty', 'not empty', 'in list', 'all')
			{
			_opVal = opVal
			_date1 = 't'
			Assert(fn() is: '', msg: 't with ' $ opVal)
			_date1 = dateNoTime
			Assert(fn() is: '', msg: 'datenotime with ' $ opVal)
			}

		for opVal in #("equals", "greater than", "greater than or equal to",
			"less than", "less than or equal to", "equals", "not equal to")
			{
			_opVal = opVal
			_date1 = 't'
			Assert(fn() is: noTimeFormat, msg: 't with ' $ opVal)
			_date1 = dateNoTime
			Assert(fn() is: '', msg: 'datenotime with ' $ opVal)
			}

		_date1 = dateNoTime
		_date2 = 't'
		_opVal = 'range'
		Assert(fn() is: noTimeFormat $ ' to' $ noTimeFormat, msg: 't with range')

		_date1 = 'm'
		_date2 = dateNoTime
		Assert(fn() is: ' ' $
			dateNoTime.Replace(day: 1).ShortDate() $ ' to' $ noTimeFormat,
				msg: 'm with notime')

		_date1 = 'm'
		_date2 = 'h'
		Assert(fn() is: ' ' $ dateNoTime.Replace(day: 1).ShortDate() $ ' to ' $
			dateNoTime.Replace(day: 1).Plus(months: 1, days: -1).ShortDate(),
			msg: 'm with h')

		_date1 = _date2 = dateNoTime
		Assert(fn() is: '', msg: 'date1 date2 datenotime')

		fn = .dateStringTester(true)

		for opVal in #("equals", "greater than", "greater than or equal to",
			"less than", "less than or equal to", "equals", "not equal to")
			{
			_opVal = opVal
			_date1 = 't'
			timeFormat = ' ' $ date.ShortDateTime()
			Assert(fn() is: timeFormat, msg: 'timeformat with ' $ opVal)
			_date1 = dateWithTime
			Assert(fn() is: '', msg: 'datewithtime with ' $ opVal)
			}

		dateWithTime = date

		_opVal = 'range'
		_date1 = dateWithTime
		_date2 = 'pm0200'
		Assert(fn() is: timeFormat $ ' to ' $ dateNoTime.
			Replace(day: 1, hour: 02, minute: 00, second: 00, millisecond: 000).
				Plus(months: -1).ShortDateTime(), msg: 'range pm0200')

		_date1 = 'ph1530'
		_date2 = dateWithTime
		Assert(fn() is: ' ' $ dateNoTime.
			Replace(day: 1, hour: 15, minute: 30, second: 00, millisecond: 000).
				Plus(days: -1).ShortDateTime() $ ' to' $ timeFormat,
			msg: 'ph1530 datewithtime')

		_opVal = 'not in range'
		_date1 = 't1900'
		_date2 = 'h0150'
		Assert(fn() is: ' ' $ dateNoTime.
			Replace(hour: 19, minute: 00, second: 00, millisecond: 000).
				ShortDateTime() $ ' to ' $ dateNoTime.
			Replace(day: 1, hour: 01, minute: 50, second: 00, millisecond: 000).
				Plus(months: 1, days: -1).ShortDateTime(),
			msg: 'not in range')
		}

	dateStringTester(showTime)
		{
		cl = (`ParamsSelectControl
			{
			ParamsSelectControl_op: class { Get(_opVal = '') { return opVal } }
			ParamsSelectControl_val: class { Get(_date1 = '') { return date1 } }
			ParamsSelectControl_val2: class { Get(_date2 = '') { return date2 } }
			ParamsSelectControl_control: #(showTime: ` $ showTime $ `)
			Valid?() { return true }
			}`).SafeEval()
		return cl.ParamsSelectControl_dateString
		}
	}
