// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Valid?()
		{
		ctrl = Mock()
		ctrl.Field = Mock()
		ctrl.Field.When.Get().Return('')
		Assert(ctrl.Eval(ChooseDatesControl.Valid?))

		ctrl.Field.When.Get().Return('2017-01-01')
		Assert(ctrl.Eval(ChooseDatesControl.Valid?))

		ctrl.Field.When.Get().Return('2017-01-01,20160101')
		Assert(ctrl.Eval(ChooseDatesControl.Valid?))

		ctrl.Field.When.Get().Return('2017-01-01,20160101,20198801')
		Assert(ctrl.Eval(ChooseDatesControl.Valid?) is: false)

		ctrl.Field.When.Get().Return('2017-01-01,20160101,sssss')
		Assert(ctrl.Eval(ChooseDatesControl.Valid?) is: false)
		}
	}