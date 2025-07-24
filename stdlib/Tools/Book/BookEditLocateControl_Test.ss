// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: BookEditLocateControl
		{
		BookEditLocateControl_records(table/*unused*/)
			{
			return [
				[path: "/Trucking", name: "Orders",
					trimmed: "Orders"],
				[path: "/Trucking" name: "Order Information",
					trimmed: "OrderInformation"],
				[path: "/Trucking", name: "Rate Tables",
					trimmed: "RateTables"],
				[path: "/Tickets" name: "Tickets",
					trimmed: "Tickets"],
				[path: "/Tickets/Link App" name: "Ticket Activity",
					trimmed: "TicketActivity"],
				[path: "/Tickets", name: "Rate Tables",
					trimmed: "RateTables"],
				[path: "/Other" name: "other_table",
					trimmed: "othertable"]
				]
			}
		}

	Test_RecordNames()
		{
		matches = .cl.RecordNames("")
		Assert(matches.Size() is: 7, msg: "no prefix")

		matches = .cl.RecordNames("Order")
		Assert(matches.Size() is: 2, msg: "prefix Order")
		Assert(matches[0] is: "Trucking > Order Information")
		Assert(matches[1] is: "Trucking > Orders")

		matches = .cl.RecordNames("OrderInformation")
		Assert(matches.Size() is: 1, msg: "prefix OrderInformation")
		Assert(matches[0] is: "Trucking > Order Information")

		matches = .cl.RecordNames("Rate")
		Assert(matches.Size() is: 2, msg: "prefix Rate")
		Assert(matches[0] is: "Tickets > Rate Tables")
		Assert(matches[1] is: "Trucking > Rate Tables")

		matches = .cl.RecordNames("RateTable")
		Assert(matches.Size() is: 2, msg: "prefix RateTable")
		Assert(matches[0] is: "Tickets > Rate Tables")
		Assert(matches[1] is: "Trucking > Rate Tables")

		matches = .cl.RecordNames("Ticket")
		Assert(matches.Size() is: 2, msg: "prefix Ticket")
		Assert(matches[0] is: "Tickets > Link App > Ticket Activity")
		Assert(matches[1] is: "Tickets > Tickets")

		matches = .cl.RecordNames("nomatches")
		Assert(matches.Size() is: 0, msg: "no matches")
		}
	}