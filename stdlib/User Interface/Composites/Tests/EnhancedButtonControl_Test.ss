// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.mock = Mock(EnhancedButtonControl)
		.mock.When.EnhancedButtonControl_updateSize([anyArgs:]).CallThrough()
		.mock.When.CalcWidth([anyArgs:]).Return(100)
		.mock.Top = 5
		.mock.When.EnhancedButtonControl_updateSize([anyArgs:]).CallThrough()
		.mock.When.CalcWidth([anyArgs:]).Return(100)
		}
	Test_calcSizeNoImage()
		{
		if not Sys.Win32?()
			return
		// no text, no image, no width, no buttonWidth & buttonHeight, no buttonStyle
		.testCalcSize(false, xMin: 10, yMin: 10, text: false, width: false,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: false,
			ih: 10, iw: 0, iy: 0, ix: 0, tx: 0, tw: 0, xMin_result: 0, yMin_result: 10)
		Assert(.mock.Top is: 0)

		// text, no image, no width, no buttonWidth & buttonHeight, buttonStyle
		// should be equivalent to normal ButtonControl
		.testCalcSize(false, xMin: 20, yMin: 10, text: "test", width: false,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: true,
			ih: 10, iw: 0, iy: 0, ix: 0, tx: 0, tw: 20, xMin_result: 20, yMin_result: 10)

		// text, no image, width: 20 chars, no buttonWidth & buttonHeight, buttonStyle
		// The Xmin should be the length of width number of characters
		// Be consistent with normal ButtonControl
		.testCalcSize(false, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: true,
			ih: 10, iw: 0, iy: 0, ix: 40, tx: 40, tw: 20,
			xMin_result: 100, yMin_result: 10)

		// text, no image, width: 20 chars, buttonWidth: 50, buttonHeight: 20, buttonStyle
		// buttonWidth & buttonHeight should overwrite width
		.testCalcSize(false, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: 50, buttonHeight: 20, skip: 10, buttonStyle: true,
			ih: 20, iw: 0, iy: 0, ix: 15, tx: 15, tw: 20,
			xMin_result: 50, yMin_result: 20)
		}
	Test_calcSizeWithImage()
		{
		if not Sys.Win32?()
			return
		fakeImage = FakeObject(Height: 20, Width: 20)
		.mock.EnhancedButtonControl_imageW = 20
		.mock.EnhancedButtonControl_imageH = 20
		.mock.When.WithDC([anyArgs:]).CallThrough()

		// no text, image, no width, no buttonWidth & buttonHeight, no buttonStyle
		// image should stretch to Ymin
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: false, width: false,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: false,
			ih: 10, iw: 10, iy: 0, ix: 0, tx: 10, tw: 0, xMin_result: 10, yMin_result: 10)

		// no text, image, width: 20 chars, no buttonWidth & buttonHeight, no buttonStyle
		// image should stretch to Ymin and be centered
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: false, width: 20,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: false,
			ih: 10, iw: 10, iy: 0, ix: 45, tx: 55, tw: 0,
			xMin_result: 100, yMin_result: 10)

		// no text, image, width: 20 chars, buttonWidth: 50, buttonHeight: 30, buttonStyle
		// buttonWidth & buttonHeight should overwrite width
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: false, width: 20,
			buttonWidth: 50, buttonHeight: 30, skip: 10, buttonStyle: false,
			ih: 30, iw: 30, iy: 0, ix: 10, tx: 40, tw: 0,
			xMin_result: 50, yMin_result: 30)

		// text, image, no width, no buttonWidth & buttonHeight, buttonStyle
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: "test", width: false,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: true,
			ih: 10, iw: 10, iy: 0, ix: 7, tx: 13, tw: 20,
			xMin_result: 37, yMin_result: 10)

		// text, image, width: 20 chars, no buttonWidth & buttonHeight, buttonStyle
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: true,
			ih: 10, iw: 10, iy: 0, ix: 35, tx: 41, tw: 20,
			xMin_result: 100, yMin_result: 10)

		// text, image, width: 20 chars, no buttonWidth & buttonHeight, no buttonStyle
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: false, buttonHeight: false, skip: 10, buttonStyle: false,
			ih: 10, iw: 10, iy: 0, ix: 35, tx: 45, tw: 20,
			xMin_result: 100, yMin_result: 10)

		// text, image, width: 20 chars, buttonWidth: 50, buttonHeight: 30, buttonStyle
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: 50, buttonHeight: 30, skip: 10, buttonStyle: true,
			ih: 30, iw: 30, iy: 0, ix: 7, tx: 33, tw: 20,
			xMin_result: 57, yMin_result: 30)

		// text, image, width: 20 chars, buttonWidth: 50, buttonHeight: 30, no buttonStyle
		.testCalcSize(fakeImage, xMin: 20, yMin: 10, text: "test", width: 20,
			buttonWidth: 50, buttonHeight: 30, skip: 10, buttonStyle: false,
			ih: 30, iw: 30, iy: 0, ix: 0, tx: 30, tw: 20,
			xMin_result: 50, yMin_result: 30)
		}

	testCalcSize(fakeImage, xMin, yMin, text, width, buttonWidth, buttonHeight,
		skip, buttonStyle, ih, iw, iy, ix, tx, tw, xMin_result, yMin_result)
		{
		.mock.EnhancedButtonControl_imageObj = fakeImage
		.mock.Hwnd = 0

		.mock.Xmin = xMin
		.mock.Ymin = yMin
		.mock.EnhancedButtonControl_text = text
		.mock.EnhancedButtonControl_width = width
		.mock.EnhancedButtonControl_buttonWidth = buttonWidth
		.mock.EnhancedButtonControl_buttonHeight = buttonHeight
		.mock.EnhancedButtonControl_skip = skip
		.mock.EnhancedButtonControl_buttonStyle = buttonStyle
		.mock.Eval(EnhancedButtonControl.EnhancedButtonControl_calcSize)
		Assert(.mock.EnhancedButtonControl_ih is: ih)
		Assert(.mock.EnhancedButtonControl_iw is: iw)
		Assert(.mock.EnhancedButtonControl_iy is: iy)
		Assert(.mock.EnhancedButtonControl_ix is: ix)
		Assert(.mock.EnhancedButtonControl_tx is: tx)
		Assert(.mock.EnhancedButtonControl_tw is: tw)
		Assert(.mock.Xmin is: xMin_result)
		Assert(.mock.Ymin is: yMin_result)
		}
	}