// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
MultiAutoChooseControl
	{
	Name: "AttachmentLabels"
	New(width = 50, height = 3, style = 16 /*lowercase style*/)
		{
		super(.list, :width, :height, allowOther:, :style)
		AttachmentLabels.Ensure()
		.FieldMenu = .FieldMenu.Copy().Append(#('', 'Access'))
		}
	list(prefix)
		{
		return AttachmentLabels.GetLabels(prefix, limit: 20)
		}
	On_Access()
		{
		AccessGoTo('Biz_AttachmentLabelsControl', 'attachlbl_label', '', .Hwnd,
			onDestroy: { SetFocus(.Hwnd) })
		}
	}