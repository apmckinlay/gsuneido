// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_NM_CLICK()
		{
		mock = Mock(EditorStatusbarControl)
		mock.When.Send([anyArgs:]).Do({ })
		mock.When.NM_CLICK().CallThrough()

		mock.EditorStatusbarControl_error = ''
		mock.NM_CLICK()
		mock.Verify.Never().Send([anyArgs:])

		mock.EditorStatusbarControl_error = 'No numbers'
		mock.NM_CLICK()
		mock.Verify.Never().Send([anyArgs:])

		mock.EditorStatusbarControl_error = 'Syntax Error invalid num 29a9'
		mock.NM_CLICK()
		mock.Verify.Never().Send([anyArgs:])

		mock.EditorStatusbarControl_error = 'Syntax Error invald a45 num'
		mock.NM_CLICK()
		mock.Verify.Never().Send([anyArgs:])

		mock.EditorStatusbarControl_error = 'Syntax Error end of line 105'
		mock.NM_CLICK()
		mock.Verify.Send('SetFirstVisibleLine', 95)

		mock.EditorStatusbarControl_error = 'Syntax Error middle 299 of line'
		mock.NM_CLICK()
		mock.Verify.Send('SetFirstVisibleLine', 289)

		mock.EditorStatusbarControl_error = 'Syntax Error near rec 5 start'
		mock.NM_CLICK()
		mock.Verify.Send('SetFirstVisibleLine', 0)
		}

	Test_Status()
		{
		mock = Mock(EditorStatusbarControl)
		mock.When.set([anyArgs:]).Do({ })
		mock.When.Status([anyArgs:]).CallThrough()

		date1 = Date()
		date2 = Date()
		recs = [
			[name: 'Rec0', group: -1, lib_modified: '', 	lib_committed: ''],
			[name: 'Rec1', group: -1, lib_modified: date1, 	lib_committed: ''],
			[name: 'Rec2', group: -1, lib_modified: '', 	lib_committed: date2],
			[name: 'Rec3', group: -1, lib_modified: date1, 	lib_committed: date2]
			]
		mock.When.rec([anyArgs:]).Return(recs[0], recs[1], recs[2], recs[3], false)

		date1 = '   Modified: '  $ date1.ShortDateTime()
		date2 = '   Committed: ' $ date2.ShortDateTime()

		mock.Status('\terror text\tright text')
		mock.Verify.set('', 'error text', 'right text', false)

		mock.Status('\terror text\tright text', invalid:)
		mock.Verify.set(date1, 'error text', 'right text', CLR.ErrorColor)

		mock.Status('\terror text\tright text', valid:)
		mock.Verify.set(date2, 'error text', 'right text', CLR.GREEN)

		mock.Status('\terror text\tright text', normal:)
		mock.Verify.set(date1 $ date2, 'error text', 'right text', false)

		mock.Status('left text\terror text\tright text')
		mock.Verify.set('left text', 'error text', 'right text', false)
		}

	Test_getColor()
		{
		mock  = Mock(EditorStatusbarControl)
		mock.When.getColor([anyArgs:]).CallThrough()
		mock.EditorStatusbarControl_error = ''

		// No specific setting, return gray (false)
		Assert(mock.getColor(invalid: false, valid: false) is: false)

		// Valid, return green
		Assert(mock.getColor(invalid: false, valid:) is: CLR.GREEN)

		// Invalid, return red
		Assert(mock.getColor(invalid:, valid: false) is: CLR.ErrorColor)

		// Invalid and Valid, prioritize red
		Assert(mock.getColor(invalid:, valid:) is: CLR.ErrorColor)

		// No specific setting, previous color was red, no error present, return gray
		mock.EditorStatusbarControl_color = CLR.ErrorColor
		Assert(mock.getColor(invalid: false, valid: false) is: false)

		// No specific setting, previous color was red, error present, return red
		mock.EditorStatusbarControl_error = 'error present, return red'
		Assert(mock.getColor(invalid: false, valid: false) is: CLR.ErrorColor)
		}

	Test_statusStr()
		{
		fn = EditorStatusbarControl.EditorStatusbarControl_dateStr
		Assert(fn('','') is: '')

		modified = #20190627.1200
		modifiedStr = #20190627.1200.ShortDateTime()
		Assert(fn(modified, '') is: '   Modified: ' $ modifiedStr)

		committed = #20190626.1200
		committedStr = #20190626.1200.ShortDateTime()
		Assert(fn(modified, committed)
			is: '   Modified: ' $ modifiedStr $ '   Committed: ' $ committedStr)

		Assert(fn('', committed) is: '   Committed: ' $ committedStr)
		}
	}