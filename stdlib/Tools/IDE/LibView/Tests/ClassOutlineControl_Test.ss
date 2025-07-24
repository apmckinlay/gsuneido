// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_loadHierarchy()
		{
		name = .TempName().Capitalize()
		child = name $ '2'
		.MakeLibraryRecord([name: 'ClassOutlineControl',
			text: `_ClassOutlineControl { Abc() {} }`])
		.MakeLibraryRecord([:name, text: `ClassOutlineControl { Abc() {} }`])
		.MakeLibraryRecord([name: child, text: name $ ` { Abc() {} }`])
		_items = Object()
		mock = .setupMock()
		mock.Eval(ClassOutlineControl.ClassOutlineControl_loadHierarchy,
			'Test_lib', child)
		Assert(_items.Count([parent: TVI.ROOT, name: 'Test_lib:' $ name]) is: 1)
		Assert(_items has: [parent: 'Test_lib:' $ name, name: 'Abc'])
		Assert(_items.Count([parent: TVI.ROOT, name: 'Test_lib:ClassOutlineControl'])
			is: 1)
		Assert(_items has: [parent: 'Test_lib:ClassOutlineControl', name: 'Abc'])
		Assert(_items.Count([parent: TVI.ROOT, name: 'stdlib:ClassOutlineControl'])
			is: 1)
		Assert(_items.Count([parent: TVI.ROOT, name: 'stdlib:Controller']) is: 1)
		Assert(_items.Count([parent: TVI.ROOT, name: 'stdlib:Container']) is: 1)
		Assert(_items.Count([parent: TVI.ROOT, name: 'stdlib:Control']) is: 1)

		.MakeLibraryRecord([name: 'TestAbcServer', text: `SocketServer { Abc() {} }`])
		_items = Object()
		mock = .setupMock()
		result = mock.Eval(ClassOutlineControl.ClassOutlineControl_loadHierarchy,
			'Test_lib', 'AbcServer')
		Assert(result)
		Assert(_items is #())
		}

	setupMock()
		{
		mock = Mock(ClassOutlineControl)
		mock.ClassOutlineControl_folder = 'folder'
		mock.ClassOutlineControl_document = 'document'
		mock.ClassOutlineControl_tree = class
			{
			AddItem(parent, name,
				image /*unused*/= 0 , container? /*unused*/ = false, param /*unused*/ = 0)
				{
				_items.Add([:parent, :name])
				return name
				}
			}
		mock.ClassOutlineControl_sort = class { Get() { return true } }
		mock.ClassOutlineControl_data = class { Get() { return true } }
		mock.ClassOutlineControl_meth = class { Get() { return true } }
		mock.ClassOutlineControl_public = class { Get() { return true } }
		mock.ClassOutlineControl_private = class { Get() { return true } }
		mock.ClassOutlineControl_maxInheritances = ClassOutlineControl.
			ClassOutlineControl_maxInheritances
		mock.When.ClassOutlineControl_addTreeNode([anyArgs:]).CallThrough()
		return mock
		}
	}