// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
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

		myContextMenu = new ContextMenu(items)
		Assert(myContextMenu.ContextMenu_menu
			is: #(
				#(name: "Set Status", id: 1, submenu: #(
					#(name: "Available", id: 2),
					#(name: "Completed", id: 3),
					#(name: "Dispatched", id: 4),
					#(name: "Overdue", id: 5),
					#(name: "Preplanned", id: 6),
					#(name: "Quote", id: 7))),
				#(name: "Equipment/Drivers", id: 8),
				#(name: "Order Confirmation", id: 9),
				#(name: "Bill Of Lading", id: 10),
				#(name: "Email Order...", id: 11),
				#(name: "Send CAA Action", id: 12),
				"",
				#(name: "Sync to Link App", id: 14, submenu: #(
					#(name: "Default Order Message", id: 15),
					#(name: "Order Confirmation", id: 16), "",
					#(name: "Delete From Link App", id: 18)))))
		}
	}