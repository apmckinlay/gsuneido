// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// make ParamsSelectControl act as a data field
ParamsSelectControl
	{
	New()
		{
		super('date_no_prompt', emptyValue: 'none')
		}

	/// force date codes to be translated to dates
	DateControl_ConvertDateCodes()
		{
		return 0
		}
	}
