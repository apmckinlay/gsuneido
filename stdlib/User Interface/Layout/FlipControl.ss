// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	// data:
	rect:		false
	controls:	false
	current:	false
	visible:	true
	Xmin:		0
	Ymin:		0
	Xstretch:	0
	Ystretch:	0
	Name:		"Flip"
	// interface:
	New(@ctrlList)
		// ctrlList is a list of .Construct -able control specifications (named
		// members are ignored)
		{
		.controls = Object()
		size = ctrlList.Size(list:)
		for (i = size - 1; 0 <= i; --i)		// add backward to minimize flicker
			{
			.controls.Add(ctrl = .Construct(ctrlList[i]), at: 0)
			if 0 < i
				ctrl.SetVisible(false)
			.Xmin = Max(ctrl.Xmin, .Xmin)
			.Ymin = Max(ctrl.Ymin, .Ymin)
			.Xstretch = Max(ctrl.Xstretch, .Xstretch)
			.Ystretch = Max(ctrl.Ystretch, .Ystretch)
			}
		.current = .controls.Size() is 0 ? false : 0
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.rect = Rect(x, y, w, h)
		// need Resize on control even if rect did not change. This handles
		// proper drawing of controls inserted/removed on the fly (dispatch filter)
		if false isnt .current
			.controls[.current].Resize(x, y, w, h)
		}
	GetControl(idx)
		{ return .controls.Member?(idx) ? .controls[idx] : false }
	GetCurrent()
		// Returns index of current control, or false if there is none
		{ return .current is false ? false : .current }
	GetCurrentControl()
		// Returns reference to current control, or false if there is none
		{ return .current is false ? false : .controls[.current] }
	SetControl(ctrl)
		{
		idx = .controls.Find(ctrl)
		.SetCurrent(idx)
		}
	SetCurrent(idx)
		// Sets the current control index to 'idx', an integer in the range
		// [0, N) , where N is the number of children of this flip control.
		{
		Assert(idx, isIntInRange: Object(0, .controls.Size()))
		prev = .controls[.current]
		prev.SetVisible(false)
		.current = idx
		.updateCurrentCtrl()
		return prev
		}
	Flip()
		// If posssible, next control sequentially in list becomes visible and a
		// reference to the previously visible control is returned. If no
		// flipping is possible, false is returned.
		{
		return false is .current
			? false
			: .SetCurrent((.current + 1) % .controls.Size()) // NOTE: false + 1 => 1
		}
	GetChildren()
		{ return .controls.Copy() }
	GetClientRect()
		{
		return .current isnt false
			? .controls[.current].GetClientRect()
			: super.GetClientRect()
		}
	GetRect()
		{ return .rect isnt false ? .rect.Copy() : super.GetRect() }
	SetVisible(visible)
		{
		Assert(visible, isBoolean:)
		if .visible isnt visible
			{
			.visible = visible
			.updateCurrentCtrl()
			}
		}
	updateCurrentCtrl()
		{
		ctrl = .GetCurrentControl()
		if false isnt ctrl
			{
			if .visible and false isnt .rect
				ctrl.Resize(.rect.GetX(), .rect.GetY(), .rect.GetWidth(),
					.rect.GetHeight())
			ctrl.SetVisible(.visible)
			}
		}
	}
