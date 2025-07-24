// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_TableNames()
		{
		matches = .cl.TableNames("")
		Assert(matches.Size() is: 5, msg: "no prefix")

		matches = .cl.TableNames("e")
		Assert(matches.Size() is: 2, msg: "prefix e")

		matches = .cl.TableNames("etaorderc")
		Assert(matches.Size() is: 1, msg: "prefix etaorderc")
		Assert(matches[0] is: "eta_order_charges")

		matches = .cl.TableNames("eta_order_")
		Assert(matches.Size() is: 1, msg: "prefix eta_order_")
		Assert(matches[0] is: "eta_order_charges")

		matches = .cl.TableNames("titicket")
		Assert(matches.Size() is: 2, msg: "prefix titicket")
		Assert(matches[0] is: "titickets", msg: "exact should be first")
		Assert(matches[1] is: "ti_tickets")

		matches = .cl.TableNames("nomatches")
		Assert(matches.Size() is: 0, msg: "no matches")
		}
	cl: SchemaLocateControl
		{
		SchemaLocateControl_tables()
			{
			return [
				[name: "eta_orders", trimmed: "etaorders"],
				[name: "eta_order_charges", trimmed: "etaordercharges"],
				[name: "titickets", trimmed: "titickets"],
				[name: "ti_tickets", trimmed: "titickets"],
				[name: "other_table", trimmed: "othertable"]
				]
			}
		}
	}