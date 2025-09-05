// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Image throws errors with prefix "Image:" as known errors
// the caller could either catch them, like: try Image(file) catch (err, "Image:") { }
// or use Image.RunWithErrorLog() { } to ignore known errors
class
	{
	CallClass(dataOrFilename)
		{
		if Paths.IsValid?(dataOrFilename) or not .saveInMemory?(dataOrFilename)
			return new this(dataOrFilename)

		return (.inMemory)(dataOrFilename)
		}

	RunWithErrorLog(block)
		{
		try
			{
			_throwForErrorLog = true
			block()
			return true
			}
		catch (err)
			{
			if not .comObjectErr(err, false) and not err.Prefix?('Image:')
				ProgrammerError(err)
			return false
			}
		}

	comObjectErr(err, _throwForErrorLog = false)
		{
		if commErr? = err.Prefix?('COMobject') and err.Has?('800a017c')
			{
			if throwForErrorLog
				throw err
			SuneidoLog('ERRATIC: ' $ err)
			}
		return commErr?
		}

	sizeThreshold: 5000
	saveInMemory?(data)
		{
		return data.Size() < .sizeThreshold
		}

	picObject: false

	New(dataOrFilename)
		{
		Assert(dataOrFilename, isString:)
		tuple = false
		pstm = Object()
		if Paths.IsValid?(dataOrFilename)
			tuple = .imageFromFile(dataOrFilename, pstm)
		else
			{
			tuple = .copydata(dataOrFilename)
			hr = CreateStreamOnHGlobal(tuple.hGlobal, true, pstm)
			if not SUCCEEDED(hr) or NULL is pstm.x
				.cleanupAndThrow("Image: failed to create stream on hGlobal", tuple)
			// Remove hGlobal from our cleanup list because passing a true value to
			// the fDeleteOnRelease parameter of CreateStreamOnHGlobal takes care of this
			tuple.Delete('hGlobal')
			}

		hr = OleLoadPicture(pstm.x, tuple.size, false, IID_IUnknown, pppic = Object())
		COMobject(pstm.x).Release()
		if not SUCCEEDED(hr) or NULL is pppic.x
			.cleanupAndThrow("Image: couldn't load picture, hr: " $ Display(hr), tuple)

		.picObject = COMobject(pppic.x)
		Assert(.picObject.Dispatch?())
		}

	ValidImageExtension: #(bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif, tiff, png)
	maxImageSize: 5 // mb, matches email attachment limit
	imageFromFile(filename, pstm)
		{
		if not .ValidImageExtension.Has?(ext = filename.AfterLast('.'))
			throw 'Image: invalid file extension: ' $ Display(ext)
		tuple = Object(size: .FileSize(filename))
		if tuple.size > .maxImageSize.Mb()
			.cleanupAndThrow("Image: " $ filename $
				": exceeds max image size: " $ ReadableSize(tuple.size), tuple)
		hr = .createFileStream(filename, pstm)
		if not SUCCEEDED(hr) or NULL is pstm.x
			.cleanupAndThrow("Image: failed to create stream on " $ filename, tuple)
		return tuple
		}

	FileSize(filename)
		{
		try
			return FileSize(filename)
		catch (e, 'FileSize:')
			{
			SuneidoLog('ERROR: (CAUGHT) ' $ e,
				caughtMsg: 'FileSize Error; Returned false')
			return false
			}
		}

	createFileStream(filename, pstm) // Extracted for tests
		{
		return SHCreateStreamOnFile(filename, 0, pstm)
		}

	Width(hdc = false)
		{
		.ckopen()
		hmWidth = .width
		return false is hdc
			? hmWidth
			: MulDiv(hmWidth, GetDeviceCaps(hdc, GDC.LOGPIXELSX), HIMETRIC_INCH)
		}
	getter_width()
		{
		return .width = .picObject.Width
		}
	Height(hdc = false)
		{
		.ckopen()
		hmHeight = .height
		return false is hdc
			? hmHeight
			: MulDiv(hmHeight, GetDeviceCaps(hdc, GDC.LOGPIXELSY), HIMETRIC_INCH)
		}
	getter_height()
		{
		return .height = .picObject.Height
		}
	Draw(hdc, x, y, w = 0, h = 0, brushImage = false, brushBackground = false)
		{
		hmWidth = .picObject.Width
		x = x.Int()         // Picture.Render() can't handle non-integer values
		y = y.Int()         // Picture.Render() can't handle non-integer values
		if (w is 0)
			w = MulDiv(hmWidth, GetDeviceCaps(hdc, GDC.LOGPIXELSX), HIMETRIC_INCH)
		else
			w = w.Int()     // Picture.Render() can't handle non-integer values
		hmHeight = .picObject.Height
		if (h is 0)
			h = MulDiv(hmHeight, GetDeviceCaps(hdc, GDC.LOGPIXELSY), HIMETRIC_INCH)
		else
			h = h.Int()     // Picture.Render() can't handle non-integer values
		rc = Rect(x, y, w, h)
		hmSize = Object(height: hmHeight, width: hmWidth) // hm = HiMetric
		if .IsRasterImage()
			.draw(hdc, rc, hmSize)
		else
			.drawWithAntialias(hdc, rc, hmSize, brushImage, brushBackground)
		}
	PICTYPE_BITMAP: 1
	PICTYPE_ICON: 3
	IsRasterImage()
		{
		return .picObject.Type is .PICTYPE_BITMAP or .picObject.Type is .PICTYPE_ICON
		}
	PaintWithAntialias(hdc, imageW, imageH, rc, block)
		{
		// pseudo antialias
		// http://www.codeproject.com/Articles/21520/Antialiasing-Using-Windows-GDI
		// http://www.gamedev.net/topic/617849-win32-draw-to-bitmap/
		scaleFactor = 4
		largeImageW = imageW * scaleFactor
		largeImageH = imageH * scaleFactor
		WithCompatibleDC(hdc, largeImageW, largeImageH)
			{|hdcBmp|
			block(hdcBmp, largeImageW, largeImageH)

			SetStretchBltMode(hdc, STRETCH.HALFTONE)
			StretchBlt(hdc,
				rc.GetX(), rc.GetY(), rc.GetWidth(), rc.GetHeight(),
				hdcBmp, 0, 0, largeImageW, largeImageH,
				ROP.SRCCOPY)
			}
		}

	copydata(data)
		{
		return Object(size: data.Size(), hGlobal: GlobalAllocData(data))
		}
	ckopen()
		{
		if false is .picObject
			throw "Image: already closed"
		}
	draw(hdc, rc, hmSize)
		{
		try
			.picObject.Render(
				hdc,
				rc.GetX(), rc.GetY(), rc.GetWidth(), rc.GetHeight(),
				0, hmSize.height - 1, hmSize.width, -hmSize.height, 0)
		catch (err)
			if not .comObjectErr(err)
				throw err
		}
	drawWithAntialias(hdc, rc, hmSize, brushImage, brushBackground)
		{
		.PaintWithAntialias(hdc, rc.GetWidth(), rc.GetHeight(), rc)
			{|hdcLarge, wLarge, hLarge|

			brush = CreateSolidBrush(CLR.WHITE)
			FillRect(hdcLarge, [right: wLarge, bottom: hLarge], brush)
			DeleteObject(brush)

			.draw(hdcLarge, Rect(0, 0, wLarge, hLarge), hmSize)

			.paintColors(hdcLarge, wLarge, hLarge, brushImage, brushBackground)
			}
		}
	paintColors(hdcDst, w, h, brushImage, brushBackground)
		{
		WithCompatibleDC(hdcDst, w, h)
			{ |hdcCopy|
			BitBlt(hdcCopy, 0, 0, w, h, hdcDst, 0, 0, ROP.SRCCOPY)

			// msdn.microsoft.com/en-us/library/windows/desktop/dd145130(v=vs.85).aspx
			// 0x0030032A is the raster operation code "PSna", which combines
			//	the colors of the brush currently selected in hdcDst, with the colors
			//	of the inverted source rectangle by using the Boolean AND operator.
			curBrush = brushImage is false ? CreateSolidBrush(CLR.BLACK) : brushImage
			DoWithHdcObjects(hdcDst, [curBrush])
				{
				BitBlt(hdcDst, 0, 0, w, h, hdcCopy, 0, 0, 0x0030032A)
				}
			if brushImage is false
				DeleteObject(curBrush)

			curbrush = brushBackground is false
				? CreateSolidBrush(CLR.WHITE) : brushBackground
			WithCompatibleDC(hdcCopy, w, h)
				{ |hdcBmpTempBk|
				DoWithHdcObjects(hdcBmpTempBk, [curbrush])
					{
					BitBlt(hdcBmpTempBk, 0, 0, w, h, hdcCopy, 0, 0, ROP.MERGECOPY)

					BitBlt(hdcDst, 0, 0, w, h, hdcBmpTempBk, 0, 0,
						ROP.SRCINVERT)
					}
				}
			if brushBackground is false
				DeleteObject(curbrush)
			}
		}
	cleanupAndThrow(error, tuple)
		{
		.cleanup(tuple)
		throw error
		}
	cleanup(tuple)
		{
		if tuple.Member?('hFile')
			{
			CloseHandle(tuple.hFile)
			tuple.Delete('hFile')
			}
		if tuple.Member?('hGlobal') and NULL isnt tuple.hGlobal
			{
			GlobalFree(tuple.hGlobal)
			tuple.Delete('hGlobal')
			}
		}

	inMemory: class
		{
		CallClass(data)
			{
			images = Suneido.GetInit("Images", Object())
			m = Adler32(data)
			if images.Member?(m)
				return images[m]
			imgOb = new this(new Image(data))
			return images[m] = imgOb
			}

		New(.image)
			{
			}

		Close()
			{
			}

		Destroy()
			{
			.image.Close()
			}

		Default(@args)
			{
			method = args[0]
			return .image[method](@+1 args)
			}
		}

	DestroyAllInMemory()
		{
		images = Suneido.GetDefault("Images", Object())
		for imgOb in images
			try imgOb.Destroy()
		Suneido.Images = Object()
		}

	Close()
		{
		.picObject.Release()
		.Delete('picObject')
		}
	}
