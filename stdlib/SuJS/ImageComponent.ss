// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
// TODO: handle drap and drop(?)
Component
	{
	Name: 'Image'
	Unsortable: true
	Xmin: false
	Ymin: false
	ContextMenu: true
	New(.color = false, .stretch = false)
		{
		.CreateElement('div')
		.SetStyles(#('align-self': 'flex-start', outline: 'none'))
		.El.draggable = true

		.origX = .Xmin
		.origY = .Ymin
		.sizeControlToImage()
		.SetMinSize()
		.setBkColor()

		.El.AddEventListener('mouseup', .mouseup)
		.El.AddEventListener('dblclick', .doubleClick)

		.El.AddEventListener(#dragstart, .dragstart)
		.El.AddEventListener(#dragend, .dragend)
		.El.AddEventListener('dragenter', .dragenter)
		.El.AddEventListener('dragover', .dragover)
		.El.AddEventListener('dragleave', .dragleave)
		.El.AddEventListener('drop', .drop)
		.El.tabIndex = "0"
		}

	mouseup(event)
		{
		if event.button isnt 0
			return
		.Event("ImageClick")
		}

	doubleClick(event/*unused*/)
		{
		.Event('LBUTTONDBLCLK')
		}

	dragstart(event)
		{
		if .GetReadOnly() or event.target isnt .El
			{
			event.StopPropagation();
			event.PreventDefault();
			return
			}

		event.dataTransfer.setData('application/suneido', 'dummy')
		.Event("ImageStartDrag")
		}

	dragend(event)
		{
		if .GetReadOnly() or event.target isnt .El
			return

		.Event("ImageEndDrag")
		}

	dragenter(event)
		{
		if .GetReadOnly() or event.target isnt .El
			return

		if not event.dataTransfer.types.Has?("Files")
			return

		event.StopPropagation();
		event.PreventDefault();
		.setBkColor('lightgreen')
		}

	dragover(event)
		{
		if .GetReadOnly() or event.target isnt .El
			return

		if not event.dataTransfer.types.HasIf?({ it in ("Files", "application/suneido") })
			return

		event.StopPropagation();
		event.PreventDefault();
		}

	dragleave(event)
		{
		if .GetReadOnly() or event.target isnt .El
			return

		event.StopPropagation();
		event.PreventDefault();
		.setBkColor()
		}

	drop(event)
		{
		if .GetReadOnly() or event.target isnt .El
			return

		dt = event.dataTransfer

		if dt.types.Has?("application/suneido")
			{
			.Event("ImageFinishDrag")
			return
			}

		files = dt.files
		if files.length is 0
			return

		hDrop = SuRender().AddDropFiles(files)
		.Event(#DROPFILES, hDrop)
		.setBkColor()

		event.StopPropagation();
		event.PreventDefault();
		}

	iw: 100
	ih: 100
	el: false
	type: false
	SetImage(img, .type, .highlight? = false, color = false)
		{
		.El.innerHTML = ""
		.el = false
		switch (.type)
			{
		case #imageFont:
			.el = CreateElement('span', .El)
			.el.SetAttribute('translate', 'no')
			.el.textContent = img.char
			.SetStyles(Object('font-family': img.font,
				'font-style': 'normal',
				'font-weight': 'normal'), .el)
			if color isnt false
				.color = color
			if .color isnt false
				.el.SetStyle('color', ToCssColor(.color))
			if false isnt fontSize = .Xmin is false ? .Ymin : .Xmin
				.el.SetStyle('font-size', fontSize $ 'px')
		case #bookImage:
			.el = CreateElement('img', .El)
			book = img.BeforeFirst('%')
			name = img.AfterFirst('%')
			.el.AddEventListener('load', .onload)
			.el.AddEventListener('error', .onerror)
			.el.SetAttribute('src', 'Res' $ Url.BuildQuery(Object(:book, :name)))
			.el.SetAttribute('alt', img)
			.el.SetStyle('overflow-wrap', 'break-word')
		case #message:
			.el = CreateElement('div', .El)
			.el.SetStyle('width', .Xmin $ 'px')
			.el.SetStyle('overflow-wrap', 'break-word')
			.el.SetStyle('pointer-events', 'none')
			.el.textContent = img
			.AddToolTip(img, .El)
			}
		.style()
		}

	sizeControlToImage()
		{
		if .Xstretch isnt false or .Ystretch isnt false
			return // size will be set by parent

		if .origX is false and .origY is false
			{
			.Xmin = .iw
			.Ymin = .ih
			}
		else if .origY is false
			.Ymin = (.ih / .iw) * .origX
		else if .origX is false
			.Xmin = (.iw / .ih) * .origY
		}

	onload(event/*unused*/)
		{
		if .el is false
			return

		.iw = .el.naturalWidth
		.ih = .el.naturalHeight
		if .origX is false or .origY is false
			{
			curX = .Xmin
			curY = .Ymin
			.sizeControlToImage()
			if curX isnt .Xmin or curY isnt .Ymin
				{
				if not .Destroyed?()
					.Window.Refresh()
				}
			}
		.sizeImageToControl()
		}

	sizeImageToControl()
		{
		if .stretch is true
			{
			.el.SetAttribute(#width, .Xmin)
			.el.SetAttribute(#height, .Ymin)
			}
		// not handle xStretch or yStretch isnt false
		else if .iw / .ih < .Xmin / .Ymin
			.el.SetAttribute(#height, .Ymin)
		else
			.el.SetAttribute(#width, .Xmin)
		}

	onerror(event/*unused*/)
		{
		if .el isnt false and .el.tagName.Lower() is 'img'
			{
			.iw = .ih = 100
			.el.SetAttribute(#width, .Xmin)
			.el.SetAttribute(#height, .Ymin)
			}
		}

	highlight?: false
	style()
		{
		if .type not in (#bookImage, #message)
			return

		if .highlight?
			.el.SetStyle('color', ToCssColor(CLR.Highlight))
		else if .readonly
			.el.SetStyle('color', ToCssColor(CLR.Inactive))
		else
			.el.SetStyle('color', '')
		}

	readonly: false
	SetReadOnly(.readonly)
		{
		.style()
		.setBkColor()
		}
	GetReadOnly()			// read-only not applicable to image
		{
		return .readonly
		}
	setBkColor(color = false)
		{
		color = color isnt false
			? color
			: .readonly
				? ''
				: 'white'
		.El.SetStyle('background-color', color)
		}
	}
