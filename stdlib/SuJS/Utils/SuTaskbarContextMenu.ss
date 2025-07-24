// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
function (id, x, y rect)
	{
	i = ContextMenu(#('Close window')).Show(false, x, y,
		rcExclude: Object(left: -9999, right: 9999, top: rect.top, bottom: rect.bottom)
		buttonRect: rect)
	if i is 0
		return
	if false is window = SuRenderBackend().WindowManager.GetWindow(id)
		return
	window.CLOSE()
	}