// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
VirtualListThumbImageButtonControl
	{
	Name: 'VirtualListExpandButton'
	New(.width = false)
		{
		super(@.layout())
		.Pushed?(false)
		}

	expandImg: 'forward.emf'
	collapseImg: 'next.emf'
	layout()
		{
		.height = .width
		return Object(command: 'VirtualListThumb_Expand', image: .expandImg,
			buttonWidth: .width, buttonHeight: .width,
			tip: 'Expand (Ctrl + Plus) / Collapse (Ctrl + Minus)',
			mouseEffect:, imagePadding: 0.1, imageColor: CLR.Inactive)
		}

	MOUSEMOVE(lParam)
		{
		.Send('VirtualListExpandBarButton_MouseMove', :lParam)
		super.MOUSEMOVE()
		}

	curRowNum: false
	LBUTTONDOWN()
		{
		super.LBUTTONDOWN()
		.Send('VirtualListThumb_Expand', .curRowNum, expand: .GetImage() is .expandImg)
		.SetImage(.GetImage() is .expandImg ? .collapseImg : .expandImg)
		return 0
		}

	UpdateButton(row_num, minus = false)
		{
		image = minus ? .collapseImg : .expandImg
		if .curRowNum is row_num and image is .GetImage()
			return
		.SetImage(image)
		.curRowNum = row_num
		}

	Reset()
		{
		.curRowNum = false
		}
	}
