// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	Valid: true
	aspectRatio: 1
	New(image, x1, y1, x2, y2, useDefaultSize = false, .name = '')
		{
		.text = .isfilename(image) ? GetFile(image) : Base64.Decode(image)
		.resourceID = 'V' $ UuidString().Replace('-', '_')
		.image = new Jpeg(.text)
		.posLocked? = false
		.aspectRatio = .image.GetWidth() / .image.GetHeight()
		if useDefaultSize is true
			{
			defaultSize = GetDpiFactor() * WinDefaultDpi
			if .aspectRatio > 1
				{
				x2 = x1 + defaultSize
				y2 = y1 + defaultSize / .aspectRatio
				}
			else
				{
				x2 = x1 + defaultSize * .aspectRatio
				y2 = y2 + defaultSize
				}
			}
		.sortPoints(x1, y1, x2, y2)
		}
	sortPoints(x1, y1, x2, y2)
		{
		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)
		}
	Paint()
		{
		width = (.x1 - .x2).Abs()
		height = (.y1 - .y2).Abs()


		if not Image.RunWithErrorLog({ _report.AddImage(.x1, .y1, width, height, .text) })
			{
			Alert('There seems to be a problem with the jpeg being used\r\n' $
				'Please click the triangle and then attempt to re-import the image',
				'Bad Jpeg')
			.text = Query1('imagebook', name: 'triangle-warning.emf').text
			_report.AddImage(.x1, .y1, width, height, .text)
			.Valid = false
			}
		}
	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}
	SetSize(x1, y1, x2, y2)
		{
		.x1 = x1
		.y1 = y1
		.x2 = x2
		.y2 = y2
		}
	ResetSize()
		{
		result = ResetSizeControl(0, Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2))
		if (result is false)
			return
		x1 = Number(result.x1)
		y1 = Number(result.y1)
		x2 = Number(result.x2)
		y2 = Number(result.y2)
		.sortPoints(x1, y1, x2, y2)

		width = (.x2 - .x1).Abs()
		height = (.y2 - .y1).Abs()
		if width / height < .aspectRatio
			.x2 = .x1 + height * .aspectRatio
		else
			.y2 = .y1 + width / .aspectRatio
		}
	StringToSave()
		{
		return 'CanvasImage(image: .' $ .resourceID $ ', x1: ' $ Display(.x1) $
			', y1: ' $ Display(.y1) $ ', x2: ' $ Display(.x2) $
			', y2: ' $ Display(.y2) $ ')'
		}
	ObToSave()
		{
		return Object('CanvasImage', Base64.Encode(.text), .x1, .y1, .x2, .y2)
		}
	Resize(origx, origy, x, y)
		{
		varyx = varyy = 'none'
		if .Resizing?(.x1, origx)
			{
			varyx = 'left'
			.x1 = x
			}
		if .Resizing?(.y1, origy)
			{
			varyy = 'top'
			.y1 = y
			}
		if .Resizing?(.x2, origx)
			{
			varyx = 'right'
			.x2 = x
			}
		if .Resizing?(.y2, origy)
			{
			varyy = 'bottom'
			.y2 = y
			}

		rect = Object(left: .x1, right: .x2, top: .y1, bottom: .y2)
		CanvasImage_UpdateWithAspectRatio(varyx, varyy, rect, .aspectRatio)
		.x1 = rect.left
		.x2 = rect.right
		.y1 = rect.top
		.y2 = rect.bottom

		.sortPoints(.x1, .y1, .x2, .y2)
		}
	Move(dx, dy)
		{
		if .posLocked?
			return
		.x1 += dx
		.x2 += dx
		.y1 += dy
		.y2 += dy
		}
	Scale(by)
		{
		.x1 *= by
		.x2 *= by
		.y1 *= by
		.y2 *= by
		.sortPoints(.x1, .y1, .x2, .y2)
		}

	DoResizeMove(x, y, varyx, varyy, rect)
		{
		if varyx isnt 'none'
			rect[varyx] = x
		if varyy isnt 'none'
			rect[varyy] = y
		CanvasImage_UpdateWithAspectRatio(varyx, varyy, rect, .aspectRatio)
		}

	GetName()
		{ return .name }
	GetResource()
		{ return Object(Object(text: Base64.Encode(.text), id: .resourceID)) }
	isfilename(str)
		{
		fileNameLimit = 256
		if fileNameLimit < str.Size()
			return false
		for (c in str)
			if c < ' ' or '~' < c
				return false
		return true
		}

	GetSuJSObject()
		{
		return Object('SuCanvasImage', .x1, .y1, .x2, .y2, Base64.Encode(.text),
			.aspectRatio, id: .Id)
		}
	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}
