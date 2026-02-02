// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
// REFERENCE:
//	https://docs.microsoft.com/en-us/windows/win32/gdi/capturing-an-image
class
	{
	CallClass(hwnd) // save and convert to png
		{
		SetFocus(NULL)
		if false is data = OkCancel(.ask, .ask.Title, hwnd)
			return
		bmp = GetAppTempFullFileName('bookscreenshot')
		if false is .SaveBmp(hwnd, bmp)
			{
			DeleteFile(bmp)
			Alert('Save Screenshot failed.', 'Save Screenshot', hwnd, MB.ICONSTOP)
			return
			}
		.save(data, bmp, hwnd)
		}

	save(data, bmp, hwnd)
		{
		path = GetAppTempPath()
		screenshotName = data.screenshot_name isnt ''
			? path $ data.screenshot_name
			: data.screenshot_file
		jpg = screenshotName.Replace('png$', 'jpg')

		magick = OptContribution('ImageMagickApp', { 'magick' })()
		Spawn(P.WAIT, magick, bmp, '-strip',
			'-dither', 'FloydSteinberg', // Dithering method
			'-colors', '256', // Reduce to 8-bit palette (similar to pngquant)
			'-depth', '8', '-define', 'png:compression-level=9', screenshotName)

		table = Suneido.CurrentBook $ 'Help'
		if data.screenshot_name isnt ''
			ImportSvcTableText(screenshotName, table, '/res', :hwnd, quiet:)

		DeleteFile(bmp)
		DeleteFile(jpg)
		if data.screenshot_name isnt ''
			DeleteFile(screenshotName)
		displayMsg = data.screenshot_name isnt ''
			? data.screenshot_name $ ' is saved into ' $ table
			: 'Screenshot is saved to ' $ data.screenshot_file
		Alert(displayMsg, 'Save Screenshot', hwnd, MB.ICONINFORMATION)
		}

	ask: Controller
		{
		Title: 'Save Screenshot (png)'
		New()
			{
			.FindControl('screenshot_name').SetFocus()
			}
		Controls: (Record (Vert
			(RadioGroups
				(Pair
					(Static 'Screenshot Name')
					(Field, width: 20, name: 'screenshot_name')
					label: 'Save screenshot to Help')
				(Pair
					(Static 'File Name')
					(SaveFile width: 20, name: 'screenshot_file')
					label: 'Save screenshot to File')
				name: 'screenshot')
			'Skip'
			'OkCancel'
			'Skip'
			(Status name: 'status')))
		On_OK()
			{
			.Send('On_OK')
			}
		OK()
			{
			data = .Data.Get()
			status = .FindControl('status')
			status.SetValid(true)
			status.Set('')
			if data.screenshot_name is '' and data.screenshot_file is ''
				{
				status.SetValid(false)
				status.Set('Please enter either a Screenshot Name or File Name')
				return false
				}
			if data.screenshot_name isnt '' and not data.screenshot_name.Suffix?('.png')
				{
				status.SetValid(false)
				status.Set('Screenshot Name must be a png file.')
				return false
				}
			if data.screenshot_file isnt '' and not data.screenshot_file.Suffix?('.png')
				{
				status.SetValid(false)
				status.Set('File Name must be a png file.')
				return false
				}
			return [screenshot_name: data.screenshot_name,
				screenshot_file: data.screenshot_file]
			}
		}

	SaveBmp(hwnd, bmpFile)
		{
		save = new this(hwnd, bmpFile)
		return save.Status
		}

	New(.hwnd, bmpFile)
		{
		try
			.saveToBmp(hwnd, bmpFile)
		catch(err)
			{
			.Status = false
			.cleanUp()
			throw err
			}
		.cleanUp()
		}

	saveToBmp(hwnd, bmpFile)
		{
		.initHdcValues(hwnd, rcClient = Object())
		DoWithHdcObjects(.hdcMemDC, [.hbmScreen])
			{
			Assert(BitBlt(.hdcMemDC,
				0, 0, rcClient.right-rcClient.left, rcClient.bottom-rcClient.top,
				.hdcWindow,
				rcClient.left, rcClient.top, ROP.SRCCOPY))
			bmpScreen = GetObjectBitmap(.hbmScreen)
			bmi = .buildBitMapInfo(bmpScreen)
			lpbitmap = .getDIBits(bmpScreen, bmi)
			.saveBmpFile(bmi, lpbitmap, bmpFile)
			}
		.Status = true
		}

	initHdcValues(hwnd, rcClient)
		{
		.hdcWindow = GetWindowDC(NULL)
		.hdcMemDC = CreateCompatibleDC(.hdcWindow)
		Assert(.hdcMemDC isnt: 0)

		dwmwa_EXTENDED_FRAME_BOUNDS = 9
		DwmGetWindowAttributeRect(hwnd, dwmwa_EXTENDED_FRAME_BOUNDS, rcClient,
			RECT.Size())

		.hbmScreen = CreateCompatibleBitmap(.hdcWindow,
			rcClient.right-rcClient.left, rcClient.bottom-rcClient.top) // HBITMAP
		Assert(.hbmScreen isnt: 0)
		}

	buildBitMapInfo(bmpScreen)
		{
		bmi = Object() // BITMAPINFO
		bmi.bmiHeader = Object() // BITMAPINFOHEADER
		bmi.bmiHeader.biSize = BITMAPINFOHEADER.Size()
		bmi.bmiHeader.biWidth = bmpScreen.bmWidth
		bmi.bmiHeader.biHeight = bmpScreen.bmHeight
		bmi.bmiHeader.biPlanes = 1
		bmi.bmiHeader.biBitCount = 32
		return bmi
		}

	getDIBits(bmpScreen, bmi)
		{
		dwBmpSize = ((bmpScreen.bmWidth * bmi.bmiHeader.biBitCount + 31) / 32).Int() /*=
			Compute the number of bytes in the array of color
			indices and store the result in biSizeImage.
			The width must be DWORD aligned unless the bitmap is RLE
			compressed. */
		dwBmpSize *= 4 * bmpScreen.bmHeight /*= color */

		.hDIB = GlobalAlloc(GMEM.MOVEABLE | GMEM.ZEROINIT, dwBmpSize)
		lpbitmap = GlobalLock(.hDIB)

		lines = GetDIBits(.hdcMemDC, .hbmScreen, 0, bmi.bmiHeader.biHeight,
			lpbitmap, bmi, 0 /*DIB_RGB_COLORS*/)
		Assert(lines greaterThan: 0)

		return lpbitmap
		}

	saveBmpFile(bmi, lpbitmap, bmpFile)
		{
		bmfHeader = Object() // BITMAPFILEHEADER
		bmfHeader.bfOffBits = BITMAPFILEHEADER.Size() + BITMAPINFOHEADER.Size()
		bmfHeader.bfSize = bmfHeader.bfOffBits + bmi.bmiHeader.biSizeImage
		bmfHeader.bfType = 0x4d42 //BM

		// using Win32 file so we can write lpbitmap directly
		hFile = CreateFile(bmpFile, GENERIC_READ | GENERIC_WRITE, 0,
			NULL, CREATE_ALWAYS, FILE_ATTRIBUTE.NORMAL, NULL)
		dwBytesWritten = Object()
		WriteFile(hFile, BITMAPFILEHEADER(bmfHeader), BITMAPFILEHEADER.Size(),
			dwBytesWritten, NULL)
		WriteFile(hFile, BITMAPINFOHEADER(bmi.bmiHeader), BITMAPINFOHEADER.Size(),
			dwBytesWritten, NULL)
		WriteFilePtr(hFile, lpbitmap, bmi.bmiHeader.biSizeImage, dwBytesWritten, NULL)
		CloseHandle(hFile)
		}

	cleanUp()
		{
		GlobalUnlock(.hDIB)
		GlobalFree(.hDIB)
		DeleteObject(.hbmScreen)
		DeleteObject(.hdcMemDC)
		ReleaseDC(.hwnd, .hdcWindow)
		.hDIB = .hbmScreen = .hdcMemDC = .hdcWindow = NULL
		}
	}
