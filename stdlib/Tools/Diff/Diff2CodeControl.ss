// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(title, listNew, listOld, titleNew, titleOld, lib = '', recName = '',
		gotoButton = false, comment = '', commentBgColor = false, newOnRight? = false,
		extraControls = #(Skip))
		{
		Window(
			Object(this, listNew, listOld, titleNew, titleOld, lib, recName,
				:gotoButton, :comment,
				:commentBgColor, :newOnRight?, :extraControls),
			:title, keep_placement: 'Diff2Control')
		}

	New(listNew, listOld, titleNew, titleOld, lib = '', recName = '',
		gotoButton = false, comment = '', commentBgColor = false, newOnRight? = false,
		extraControls = #(Skip))
		{
		super(.layout(listNew, listOld, titleNew, titleOld, lib, recName,
			:gotoButton, :comment, :commentBgColor, :newOnRight?, :extraControls))
		}

	layout(listNew, listOld, titleNew, titleOld, lib = '', recName = '',
		gotoButton = false, comment = '', commentBgColor = false, newOnRight? = false,
		extraControls = #(Skip))
		{
		return Object('Vert', Object('Diff2',
			listNew, listOld, lib, recName, titleNew, titleOld,
			:gotoButton, :comment, :commentBgColor, :newOnRight?, :extraControls))
		}
	}
