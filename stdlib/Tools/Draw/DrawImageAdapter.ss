// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// TODO: gif is not working, or maybe only allows jpg
function (x1, y1)
	{
	maxFileSize = 100.Kb()
	image = Ask(prompt: 'Choose an image',
		ctrl: Object('Vert'
			Object('OpenFile' title: 'Image' filter: "JPEG (*.jpg)\x00*.jpg",
				name: 'file')))
	if image is false or image.file is ''
		return false
	if false is image.file.AfterLast('.').Lower() in ('jpg', 'jpe', 'jpeg')
		{
		Alert("File format must be JPEG (ie. .jpg, .jpe, .jpeg)")
		return false
		}
	else if FileSize(image.file) > maxFileSize
		{
		Alert("Selected file exceeds the maximum size limitation (" $
			ReadableSize(maxFileSize) $ ")")
		return false
		}
	else
		{
		Jpeg.RunWithCatch()
			{
			return CanvasImage(image.file, x1, y1, x1, y1, useDefaultSize:)
			}
		Alert("Unable to load the file, it may be corrupt")
		return false
		}
	}