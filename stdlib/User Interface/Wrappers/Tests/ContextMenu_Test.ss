// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_AddItems()
		{
		items = #(
			#(name: "Set Status"),
			#("Available", "Completed", "Dispatched", "Overdue", "Preplanned", "Quote"),
			"Equipment/Drivers",
			"Order Confirmation",
			"Bill Of Lading",
			"Email Order...",
			"Send CAA Action",
			"",
			"Sync to Link App",
			#("Default Order Message", "Order Confirmation", "", "Delete From Link App"))

		cl = ContextMenu
			{
			Getter_Menu()
				{
				return .Menu = Object()
				}
			InsertItem(item = false, pos = 0, newhandle = false, prefixes = #())
				{
				.Menu[newhandle][pos] = Object(:item, :prefixes)
				}
			ContextMenu_createPopupMenu()
				{
				newMenu = .Menu.Size()
				.Menu[newMenu] = Object()
				return newMenu
				}
			ContextMenu_setMenuItemInfo(@unused) { }
			ContextMenu_getMenuItemCount(menu)
				{
				return .Menu[menu].Size()
				}
			}
		myContextMenu = new cl(items)
		Assert(myContextMenu.Menu isSize: 3)
		Assert(myContextMenu.Menu[0] isSize: 8)
		Assert(myContextMenu.Menu[1] isSize: 6)
		Assert(myContextMenu.Menu[2] isSize: 4)

		Assert(myContextMenu.Menu[0][0] is: #(item: #(name: "Set Status"), prefixes: #()))
		Assert(myContextMenu.Menu[0][2]
			is: #(item: #(name: "Order Confirmation"), prefixes: #()))
		Assert(myContextMenu.Menu[1][4]
			is: #(item: #(name: "Preplanned"), prefixes: #("Set_Status")))
		Assert(myContextMenu.Menu[2][1]
			is: #(item: #(name: "Order Confirmation"), prefixes: #("Sync_to_Link_App")))
		}
	}
