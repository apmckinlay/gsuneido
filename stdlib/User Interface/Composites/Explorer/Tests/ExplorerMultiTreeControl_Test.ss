// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Expanding()
		{
		_orderAdded = Object()
		cl = ExplorerMultiTreeControl
			{
			AddItem(parent /*unused*/, name, image /*unused*/ = 0 ,
				container? /*unused*/ = false, param /*unused*/ = 0)
				{
				_orderAdded.Add(name)
				return true
				}
			}
		mock = Mock(cl)
		mock.When.Expanding([anyArgs:]).CallThrough()
		mock.When.addChildren([anyArgs:]).CallThrough()
		mock.When.GetParam([anyArgs:]).Return(0, 0, 0, 1, 3, 4)
		mock.When.GetChildren([anyArgs:]).Return(#(), #(), #(1, 3, 4))

		mock.ExplorerMultiTreeControl_imageHandler = []
		mock.ExplorerMultiTreeControl_model = FakeObject(
			Children: [
				[group: false, name: 'A', num: 1],
				[group: false, name: 'B', num: 2],
				[group: false, name: 'C', num: 3],
				[group:, name: 'D', num: 4],
				[group:, name: 'E', num: 5],
				[group:, name: 'F', num: 6]],
			Get: [],
			Container?: false)

		mock.ExplorerMultiTreeControl_inorder = false
		_orderAdded.Delete(all:)
		mock.Expanding(0)
		Assert(_orderAdded is: #(D, E, F, A, B, C))

		_orderAdded.Delete(all:)
		mock.ExplorerMultiTreeControl_inorder = true
		mock.Expanding(0)
		Assert(_orderAdded is: #(A, B, C, D, E, F))

		// A, C, and D are not added as they were already "added" (faked via mock)
		_orderAdded.Delete(all:)
		mock.ExplorerMultiTreeControl_inorder = true
		mock.Expanding(0)
		Assert(_orderAdded is: #(B, E, F))
		}

	Test_addItem()
		{
		img = new ExplorerMultiImageHandler
			{
			ExplorerMultiImageHandler_init() { }
			}
		mock = Mock(ExplorerMultiTreeControl)
		mock.ExplorerMultiTreeControl_imageHandler = img
		mock.When.addItem([anyArgs:]).CallThrough()
		mock.When.AddItem([anyArgs:]).Do({ })
		mock.ExplorerMultiTreeControl_model = class
			{
			Modified?(data) // Based off LibTreeModel
				{ return data.lib_modified isnt '' or data.lib_committed is '' }
			}

		mock.addItem(0, #name, #num, false)
		mock.Verify.AddItem(0, #name, img.ModifiedTheme, false, #num)

		mock.addItem(0, #name, #num, false, '', Date())
		mock.Verify.AddItem(0, #name, img.DocumentTheme, false, #num)

		mock.addItem(0, #name, #num, false, Date())
		mock.Verify.Times(2).AddItem(0, #name, img.ModifiedTheme, false, #num)

		mock.addItem(0, #name, #num, true, Date())
		mock.Verify.AddItem(0, #name, img.FolderClosedModifiedTheme, true, #num)

		mock.addItem(0, #name, #num, true)
		mock.Verify.Times(2).AddItem(0, #name, img.FolderClosedModifiedTheme, true, #num)
		}
	}
