// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		LibTreeModel.Create("testlibfrom")
		LibTreeModel.Create("testlibto")
		Database("ensure testlibfrom (lib_modified, lib_committed)")
		Database("ensure testlibto (lib_modified, lib_committed)")
		}
	Test_main()
		{
		ctrl = Mock(SvcCopyRecordsControl)
		ctrl.When.print([anyArgs:]).Return(0)
		table = ctrl.SvcCopyRecordsControl_table = Mock()
		table.When.Get().Return('')
		ctrl.Eval(SvcCopyRecordsControl.OK)
		ctrl.Verify.AlertError('Copy records', 'Please select a target library')

		ctrl.SvcCopyRecordsControl_srcs = Object(Object(lib: 'testlib', name: 'TestRec'))
		table.When.Get().Return('testlib')
		ctrl.Eval(SvcCopyRecordsControl.OK)
		ctrl.Verify.AlertError('Copy records', 'Cannot copy records to same library')

		QueryOutput('testlibfrom', Record(name: 'Test1', text: 'text', num: 1, group: -1,
			lib_modified: Date()))
		QueryOutput('testlibfrom', Record(name: 'Test2', text: 'text', num: 2, group: -1,
			lib_modified: Date()))
		QueryOutput('testlibfrom', Record(name: 'Test3', text: 'text', num: 3, group: -1,
			lib_modified: Date()))

		ctrl.SvcCopyRecordsControl_srcs = Object(
			Object(lib: 'testlibfrom', name: 'Test1'),
			Object(lib: 'testlibfrom', name: 'Test2'),
			Object(lib: 'testlibfrom', name: 'Test3'))
		table.When.Get().Return('testlibto')
		folder = ctrl.SvcCopyRecordsControl_folder = Mock()
		folder.When.Get().Return('')
		overwrite = ctrl.SvcCopyRecordsControl_overwrite = Mock()
		overwrite.When.Get().Return(false)
		Assert(ctrl.Eval(SvcCopyRecordsControl.OK)
			is: 'THREE record(s) copied to testlibto')

		Assert(Query1('testlibto where name is "Test1"').num is: 1)
		Assert(Query1('testlibto where name is "Test2"').num is: 2)
		Assert(Query1('testlibto where name is "Test3"').num is: 3)
		}
	Teardown()
		{
		try Database("destroy testlibfrom")
		try Database("destroy testlibto")
		super.Teardown()
		}
	}