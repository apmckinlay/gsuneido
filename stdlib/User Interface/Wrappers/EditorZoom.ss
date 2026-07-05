// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
function(ctrl, zoom, zoomDialog = false, font = '', size = '', textLimit = false)
	{
	if zoom
		{
		ctrl.Send("CloseZoom")
		return
		}
	if zoomDialog is false
		zoomDialog = ZoomControl

	ctrl.Hasfocus? = true
	hwnd = ctrl.Window.MainHwnd()
	text = ctrl.Get()
	readonly = ctrl.GetReadOnly()
	zoom_text = zoomDialog(hwnd, text, :readonly, :font, :size, :textLimit)
	if not readonly and String?(zoom_text) and zoom_text isnt text
		{
		ctrl.Set(zoom_text)
		ctrl.Dirty?(true)
		}
	ctrl.Hasfocus? = false
	if ctrl.Method?('EnsureSelect')
		ctrl.EnsureSelect()
	}