// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	images: #()
	New(image1, image2, .highlighted = 0, gap = 1)
		{
		.images = Object(ImageResource(image1), ImageResource(image2))
		.highlightBrush = CreateSolidBrush(CLR.Highlight)
		.gap = ScaleWithDpiFactor(gap)
		}

	IsRasterImage()
		{
		return false
		}

	Width(hdc)
		{
		return .images.Map({ it.Width(hdc) }).Sum() + .gap
		}

	Height(hdc)
		{
		return Max(@.images.Map({ it.Height(hdc) }))
		}

	Draw(hdc, x, y, w = 0, h = 0, brushImage = false,
		brushBackground = false)
		{
		if .images.Empty?()
			return
		.images[0].Draw(hdc, x, y, w, h,
			.highlighted is 0 ? .highlightBrush : brushImage,
			brushBackground)
		.images[1].Draw(hdc, x + (w - .gap) / 2 + .gap, y, w, h,
			.highlighted is 1 ? .highlightBrush : brushImage,
			brushBackground)
		}

	Close()
		{
		.images.Each(#Close)
		.images = #()
		DeleteObject(.highlightBrush)
		}
	}
