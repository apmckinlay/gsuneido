// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_NumberControl
	{
	ComponentName: "Number"
	New(.mask = '-###,###,###', readonly = false,
		rangefrom = false, rangeto = false, width = false,
		set = false, mandatory = false, status = "", justify = "RIGHT",
		font = "", size = "", weight = "", underline = false,
		hidden = false, tabover = false)
		{
		super(mask, readonly, rangefrom, rangeto, width,
		set, mandatory, status, justify, font, size, weight, underline, hidden, tabover)

		.ComponentArgs = Object(.mask, readonly, width, justify,
			font, size, weight, underline, tabover)
		}

	focused?: false
	EN_SETFOCUS()
		{
		skipSetText? = .focused?
		.focused? = true
		super.EN_SETFOCUS(:skipSetText?)
		}
	KillFocus()
		{
		.focused? = false
		super.KillFocus()
		}
	}