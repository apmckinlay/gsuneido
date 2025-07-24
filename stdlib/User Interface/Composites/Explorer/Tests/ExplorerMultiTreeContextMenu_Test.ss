// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	explorer: class
		{
		New()
			{
			.RootSelected? = true
			.Controller = new class
				{
				New() { .CanPaste? = .AllowRootDelete? = false }
				AllowRootDelete?() { .AllowRootDelete? }
				CanPaste?() { .CanPaste? }
				}
			}
		RootSelected?() { .RootSelected? }
		}
	Test_buildOpenCloseMenu()
		{
		instance = new ExplorerMultiTreeContextMenu(new .explorer, false)
		method = instance.ExplorerMultiTreeContextMenu_buildOpenCloseMenu

		selectionSize = 1
		container? = expanded? = children? = false
		method(menu = Object(), :selectionSize, :container?, :expanded?, :children?)
		Assert(menu is: #())

		container? = true
		method(menu = Object(), :selectionSize, :container?, :expanded?, :children?)
		Assert(menu is: #())

		expanded? = true
		method(menu = Object(), :selectionSize, :container?, :expanded?, :children?)
		Assert(menu isSize: 2)
		Assert(menu[0].name is: 'C&lose')
		Assert(menu[1].name is: '')
		Assert(menu[0].order is: 10)
		Assert(menu[1].order is: 11)

		instance.ExplorerMultiTreeContextMenu_segments = 0
		expanded? = false
		children? = true
		method(menu = Object(), :selectionSize, :container?, :expanded?, :children?)
		Assert(menu isSize: 2)
		Assert(menu[0].name is: '&Open')
		Assert(menu[1].name is: '')
		Assert(menu[0].order is: 10)
		Assert(menu[1].order is: 11)
		}

	Test_buildItemMenu_libView()
		{
		cl = ExplorerMultiTreeContextMenu
			{
			ExplorerMultiTreeContextMenu_isLibViewControl?()
				{ return true }
			}
		instance = new cl(explorer = new .explorer, false)
		method = instance.ExplorerMultiTreeContextMenu_buildItemMenu

		selectionSize = 0
		static? = true
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: '&New'],
				[order: 20, name: ''],
				[order: 21, name: 'Dump']
				])

		instance.ExplorerMultiTreeContextMenu_segments = 0
		selectionSize = 1
		static? = false
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: 'Copy To...'],
				[order: 21, name: '&Move To...'],
				[order: 30, name: ''],
				[order: 40, name: '&New'],
				[order: 50, name: ''],
				[order: 51, name: 'Dump']
			])

		instance.ExplorerMultiTreeContextMenu_segments = 0
		explorer.Controller.AllowRootDelete? = true
		explorer.Controller.CanPaste? = true
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: '&Paste'],
				[order: 30, name: 'Copy To...'],
				[order: 31, name: '&Move To...'],
				[order: 40, name: ''],
				[order: 50, name: '&New'],
				[order: 60, name: ''],
				[order: 61, name: 'Dump'],
				[order: 62, name: '&Delete'],
			])

		// Root isn't selected, so AllowDelete? shouldn't affect the outcome
		instance.ExplorerMultiTreeContextMenu_segments = 0
		explorer.Controller.AllowRootDelete? = false
		explorer.RootSelected? = false
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: '&Paste'],
				[order: 30, name: 'Copy To...'],
				[order: 31, name: '&Move To...'],
				[order: 40, name: ''],
				[order: 50, name: '&Delete'],
				[order: 51, name: 'Rena&me'],
				[order: 52, name: ''],
				[order: 60, name: '&New']
			])
		}

	Test_menu_build_order_clean_readonly()
		{
		cl = ExplorerMultiTreeContextMenu
			{
			ExplorerMultiTreeContextMenu_isLibViewControl?()
				{ return true }
			}
		instance = new cl(explorer = new .explorer, true)
		method = instance.ExplorerMultiTreeContextMenu_buildItemMenu

		selectionSize = 0
		static? = true
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: '&New'],
				[order: 20, name: ''],
				[order: 21, name: 'Dump']
				])
		instance.ExplorerMultiTreeContextMenu_orderMenu(menu, [])
		Assert(menu is: [[name: ''], [name: 'Dump']])
		instance.ExplorerMultiTreeContextMenu_cleanUpMenu(menu)
		Assert(menu	is: [[name: 'Dump']])

		instance.ExplorerMultiTreeContextMenu_segments = 0
		selectionSize = 1
		static? = false
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: 'Copy To...'],
				[order: 21, name: '&Move To...'],
				[order: 30, name: ''],
				[order: 40, name: '&New'],
				[order: 50, name: ''],
				[order: 51, name: 'Dump']
			])
		instance.ExplorerMultiTreeContextMenu_orderMenu(menu, [])
		Assert(menu is: [[name: ''], [name: ''], [name: 'Dump']])
		instance.ExplorerMultiTreeContextMenu_cleanUpMenu(menu)
		Assert(menu	is: [[name: 'Dump']])


		instance.ExplorerMultiTreeContextMenu_segments = 0
		explorer.Controller.AllowRootDelete? = true
		explorer.Controller.CanPaste? = true
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: '&Paste'],
				[order: 30, name: 'Copy To...'],
				[order: 31, name: '&Move To...'],
				[order: 40, name: ''],
				[order: 50, name: '&New'],
				[order: 60, name: ''],
				[order: 61, name: 'Dump'],
				[order: 62, name: '&Delete'],
			])
		instance.ExplorerMultiTreeContextMenu_orderMenu(menu, [])
		Assert(menu is: [[name: ''], [name: ''], [name: 'Dump']])
		instance.ExplorerMultiTreeContextMenu_cleanUpMenu(menu)
		Assert(menu	is: [[name: 'Dump']])
		}

	Test_buildItemMenu_notLibView()
		{
		cl = ExplorerMultiTreeContextMenu
			{
			ExplorerMultiTreeContextMenu_isLibViewControl?()
				{ return false }
			}
		instance = new cl(new .explorer, false)
		method = instance.ExplorerMultiTreeContextMenu_buildItemMenu

		instance.ExplorerMultiTreeContextMenu_segments = 0
		selectionSize = 1
		static? = false
		method(menu = [], :selectionSize, :static?)
		Assert(menu
			is: [
				[order: 10, name: 'C&ut'],
				[order: 11, name: '&Copy'],
				[order: 20, name: ''],
				[order: 30, name: '&New'],
				[order: 40, name: ''],
				[order: 41, name: 'Dump']
			])
		}
	}