// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
// TODO: handle mouse events
Control
	{
	Name: 'Image'
	ComponentName: 'Image'
	Unsortable: true
	Xmin: false
	Ymin: false

	New(image = '', message = 'no image', .stretch = false, .acceptDrop = false,
		.color = false, .bgndcolor = false)
		{
		.orig_message = message
		.Set(image)

		.ComponentArgs = Object(color, stretch)
		}

	image: false
	Set(image, highlight? = false, message = false, color = false)
		{
		if not String?(image)
			throw "Cannot handle image of type: " $ Type(image)

		if image is .image
			return

		if image is ''
			.Act('SetImage', .orig_message, :highlight?, type: 'message')
		else if false isnt charCode = IconFont().MapToCharCode(image)
			.Act('SetImage', charCode, :color, type: 'imageFont')
		else if .isBookImage?(image) and not  .isInvalidType?(image)
			.Act('SetImage', image, :highlight?, type: 'bookImage')
		else
			// TODO: handle image in file
			.Act('SetImage', message is false ? image : message, :highlight?,
				type: 'message')
		.image = image
		}

	isInvalidType?(image)
		{
		return Paths.IsValid?(image) and
			image !~ "(?i)[.](bmp|gif|jpg|jpe|jpeg|ico|emf|wmf|png)$"
		}

	isBookImage?(image)
		{
		if Paths.IsValid?(image) and image.Has?('%')
			{
			book = image.BeforeFirst('%')
			return TableExists?(book)
			}
		return false
		}

	readonly: false
	SetReadOnly(ro)
		{
		ro = ro is true
		.readonly = ro
		.Act(#SetReadOnly, ro)
		}
	GetReadOnly()			// read-only not applicable to image
		{
		return .readonly
		}

	ImageClick()
		{
		.Send("ImageClick")
		}

	LBUTTONDBLCLK()
		{
		.Send("ImageDoubleClick")
		return 0
		}

	DROPFILES(wParam)
		{
		.Send("ImageDropFiles", wParam)
		return 0
		}

	ImageStartDrag()
		{
		.Send("ImageStartDrag")
		}

	ImageEndDrag()
		{
		.Send("ImageEndDrag")
		}

	ImageFinishDrag()
		{
		.Send("ImageFinishDrag")
		}

	ContextMenu(x, y)
		{
		.Send("ImageContextMenu", x, y)
		return 0
		}

	// TODO: Implement me
	SetTip(tip /*unused*/) {}
	}