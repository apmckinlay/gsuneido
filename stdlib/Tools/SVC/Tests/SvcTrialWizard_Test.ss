// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: SvcTrialWizard
		{
		SvcTrialWizard_trialTags: (test1: 'test1 desc', test2: 'test2 desc')
		}

	Test_copyAndRestore()
		{
		fn = .cl.SvcTrialWizard_copyAndRestore

		select = Object(svc_name: 'Foo__webgui')
		tag = 'test1'
		mock = MockObject(#(
			((#Get, 'Foo__webgui', t: false), result: [parent: 0, text: '123']),
			(#Output, [name: 'Foo__webgui_test1', parent: 0,
				text: '123', lib_invalid_text: ''], t: false),
			(#Restore, 'Foo__webgui', t: false)))
		fn(select, mock, false, tag)

		select = Object(svc_name: 'Foo__test2')
		tag = 'test1'
		mock = MockObject(#(
			((#Get, 'Foo__test2', t: false), result: [parent: 0, text: '123']),
			(#Output, [name: 'Foo__test1', parent: 0,
				text: '123', lib_invalid_text: ''], t: false),
			(#Restore, 'Foo__test2', t: false)))
		fn(select, mock, false, tag)
		}

	Test_renameAndDelete()
		{
		fn = .cl.SvcTrialWizard_renameAndDelete

		select = Object(svc_name: 'Foo__webgui_test1')
		mock = MockObject(#(
			((#Get, 'Foo__webgui_test1', t: false), result: [parent: 0, text: '123']),
			((#Get, 'Foo__webgui', t: false), result: false),
			(#Rename, [parent: 0, text: '123'], 'Foo__webgui', t: false)))
		fn(select, mock, false)

		mock = MockObject(Object(
			Object(#(#Get, 'Foo__webgui_test1', t: false),
				result: [parent: 0, text: '123']),
			Object(#(#Get, 'Foo__webgui', t: false),
				result: [parent: 0, text: '789']),
			#(#Update, [parent: 0, text: '789', lib_invalid_text: ''],
				newText: '123', t: false)
			#(#StageDelete, 'Foo__webgui_test1', t: false)))
		fn(select, mock, false)
		}
	}