// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function (images, image, color, w, h, dark? = false, padding = 0)
	{
	w = ScaleWithDpiFactor(w)
	h = ScaleWithDpiFactor(h)
	hdc = GetWindowDC(NULL)
	memdc = CreateCompatibleDC(hdc)
	hBmp = CreateCompatibleBitmap(hdc, w, h)
	ReleaseDC(NULL, hdc)
	if color is CLR.BLACK and dark?
		color = CLR.WHITE
	brushImage = CreateSolidBrush(color)
	back = dark? ? CLR.BLACK : CLR.WHITE
	brushBackground = CreateSolidBrush(back)
	image = ImageResource(image)
	DoWithHdcObjects(memdc, [hBmp])
		{
		paddingW = (w * padding).Round(0)
		paddingH = (h * padding).Round(0)
		FillRect(memdc, [right: w, bottom: h], brushBackground)
		image.Draw(memdc, paddingW, paddingH, w - paddingW * 2, h - paddingH * 2,
			:brushImage, :brushBackground)
		}
	DeleteDC(memdc)
	DeleteObject(brushImage)
	DeleteObject(brushBackground)
	idx = ImageList_AddMasked(images, hBmp, back)
	DeleteObject(hBmp)
	image.Close()
	return idx
	}
