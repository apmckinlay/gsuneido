// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_checkPermission()
		{
		fakeScreenControl1 = Controller
			{
			AccessPermission: #('TestPermission1', 'TestPermission2')
			}

		// CASE 1: permissions are specified on the control; User has one permission
		// but not the other.  User should have permission in this case
		mock = Mock(AccessGoTo)
		mock.When.getAccess('fakeScreenControl1').Return(fakeScreenControl1)
		mock.When.getPermission('TestPermission1').Return(false)
		mock.When.getPermission('TestPermission2').Return(true)

		result = mock.Eval(AccessGoTo.CheckPermission, 'fakeScreenControl1')
		Assert(result.permission)

		mock.Verify.getPermission('TestPermission1')
		mock.Verify.getPermission('TestPermission2')
		mock.Verify.Times(2).getPermission([any:])

		// CASE 2: permissions specified on the control; user does not have permission
		mock = Mock(AccessGoTo)
		mock.When.getAccess('fakeScreenControl1').Return(fakeScreenControl1)
		mock.When.getPermission('TestPermission1').Return(false)
		mock.When.getPermission('TestPermission2').Return(false)

		result = mock.Eval(AccessGoTo.CheckPermission, 'fakeScreenControl1')
		Assert(result.permission is: false)

		mock.Verify.getPermission('TestPermission1')
		mock.Verify.getPermission('TestPermission2')
		mock.Verify.Times(2).getPermission([any:])

		// CASE 3: permission is specified on control as "All".  Everyone has permission
		fakeScreenControl2 = Controller
			{
			AccessPermission: #('All')
			}

		mock = Mock(AccessGoTo)
		mock.When.getAccess('fakeScreenControl2').Return(fakeScreenControl2)
		result = mock.Eval(AccessGoTo.CheckPermission, 'fakeScreenControl2')
		Assert(result.permission)
		mock.Verify.Never().getPermission([any:])

		// CASE 4: permission is not specified on the control.
		fakeScreenControl3 = Controller
			{
			}

		mock = Mock(AccessGoTo)
		mock.When.getAccess('fakeScreenControl3').Return(fakeScreenControl3)
		mock.When.getPermission('fakeScreenControl3').Return(true)
		result = mock.Eval(AccessGoTo.CheckPermission, 'fakeScreenControl3')
		Assert(result.permission)
		}
	}
