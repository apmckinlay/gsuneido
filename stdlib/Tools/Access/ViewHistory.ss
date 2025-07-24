// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (data, hwnd, historyFields)
	{
	if not Object?(historyFields)
		{
		Alert('No History Available', 'History', hwnd, MB.ICONINFORMATION)
		return
		}
	ModalWindow(Object('ViewHistory', fields: historyFields, :data),
		keep_size: 'ViewHistory' $ ' - ' $ data.transaction_type)
	}