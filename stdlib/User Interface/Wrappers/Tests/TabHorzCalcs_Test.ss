// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	fakeTabControl: class
		{
		Ymin: 			0
		Selected: 		0
		TrimChar:		'~'
		TrimChars:		1
		GetSelected()
			{
			return .Selected
			}
		}

	spy()
		{
		.SpyOn(TabCalcs.TabCalcs_initFont).Return('')
		.SpyOn(TabCalcs.ImageDimensions).Return([width: 10, height: 10])
		.SpyOn(ScaleWithDpiFactor).Return(6)
		}

	setTextMetrics(tab)
		{
		tab.textWidth = tab.textHeight = 12
		tab.textBoldOffset = 2
		tab.boldCharSizes = .charSizes(tab.tabName, 4)
		tab.baseCharSizes = .charSizes(tab.tabName, 4)
		tab.trimCharSize = 3
		}

	charSizes(tabName, charSize)
		{
		return tabName.Size().Of({ charSize })
		}

	Test_top()
		{
		.spy()
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'top')

		calcCl.Resize(w = 200, 20)
		Assert(calcCl.TabBarSize is: w)
		Assert(calcCl.FontOrientation is: 0)
		Assert(calcCl.TabHorzCalcs_extraControlY is: 0)

		// Tab render values
		.setTextMetrics(tab = Object(tabName: 'Tab1', hide?: false, data: false, image:))
		calcCl.CalcRenderRect(0, tab, 0)
		Assert(tab.renderRect.start is: 0)
		Assert(tab.renderRect.top is: 1)
		Assert(tab.renderRect.tipY is: 10)
		Assert(tab.renderRect.end is: 30)
		Assert(tab.renderRect.bottom is: 10)
		Assert(tab.renderRect.left is: 0)
		Assert(tab.renderRect.right is: 30)
		Assert(tab.renderRect.tipX is: 0)
		Assert(tab.renderWidth is: tab.width)

		// Tab draw values
		drawSpecs = calcCl.CalcDrawSpecs(300, 80)
		Assert(drawSpecs.ellipseSize is: 40)
		// Draw base round rectangle values
		Assert(drawSpecs.baseRound.left is: 0)
		Assert(drawSpecs.baseRound.top is: 0)
		Assert(drawSpecs.baseRound.right is: 300)
		Assert(drawSpecs.baseRound.bottom is: 80)
		// Draw base round rectangle fill values
		Assert(drawSpecs.baseFill.right is: 300)
		Assert(drawSpecs.baseFill.bottom is: 80)
		// Draw override bottom rounded corner values
		Assert(drawSpecs.overrideRect.left is: 0)
		Assert(drawSpecs.overrideRect.top is: 20)
		Assert(drawSpecs.overrideRect.right is: 300)
		Assert(drawSpecs.overrideRect.bottom is: 80)
		// Draw override bottom rounded corner fill values
		Assert(drawSpecs.overrideFill.left is: 1)
		Assert(drawSpecs.overrideFill.top is: 20)
		Assert(drawSpecs.overrideFill.right is: 299)
		Assert(drawSpecs.overrideFill.bottom is: 21)

		// Test text draw positions
		// Text draw position for selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, true)
		Assert(textSpecs.x is: 11)
		Assert(textSpecs.y is: 5)
		// Text draw position for not selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 13)
		Assert(textSpecs.y is: 5)
		// Text draw position for not selected tab with renderWidth differing from width
		tab.width += 2
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 11)
		Assert(textSpecs.y is: 5)
		tab.width -= 2

		// Test image rectangle values
		imageRect = calcCl.ImageRect(tab)
		Assert(imageRect.left is: 7)
		Assert(imageRect.top is: 4)
		Assert(imageRect.right is: 17)
		Assert(imageRect.bottom is: 14)

		// Test tab underline positions
		tabPoints = calcCl.CalcLinePoints()
		Assert(tabPoints.x1 is: 0)
		Assert(tabPoints.y1 is: 9)
		Assert(tabPoints.x2 is: 199)
		Assert(tabPoints.y2 is: 9)

		// Test invalidate rectangle values
		invalidateOb = calcCl.InvalidateRect(0, tab)
		Assert(invalidateOb.left is: 0)
		Assert(invalidateOb.top is: 0)
		Assert(invalidateOb.right is: 30)
		Assert(invalidateOb.bottom is: 10)
		}

	Test_bottom()
		{
		.spy()
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'bottom')

		calcCl.Resize(w = 200, 20)
		Assert(calcCl.TabBarSize is: w)
		Assert(calcCl.FontOrientation is: 0)
		Assert(calcCl.TabHorzCalcs_extraControlY is: 1)

		// Tab render values
		.setTextMetrics(tab = Object(tabName: 'Tab1', hide?: false, data: false, image:))
		calcCl.CalcRenderRect(0, tab, 0)
		Assert(tab.renderRect.start is: 0)
		Assert(tab.renderRect.top is: 0)
		Assert(tab.renderRect.tipY is: 9)
		Assert(tab.renderRect.end is: 30)
		Assert(tab.renderRect.bottom is: 9)
		Assert(tab.renderRect.left is: 0)
		Assert(tab.renderRect.right is: 30)
		Assert(tab.renderRect.tipX is: 0)
		Assert(tab.renderWidth is: tab.width)

		// Tab draw values
		drawSpecs = calcCl.CalcDrawSpecs(300, 80)
		Assert(drawSpecs.ellipseSize is: 40)
		// Draw base round rectangle values
		Assert(drawSpecs.baseRound.left is: 0)
		Assert(drawSpecs.baseRound.top is: 0)
		Assert(drawSpecs.baseRound.right is: 300)
		Assert(drawSpecs.baseRound.bottom is: 80)
		// Draw base round rectangle fill values
		Assert(drawSpecs.baseFill.right is: 300)
		Assert(drawSpecs.baseFill.bottom is: 80)
		// Draw override bottom rounded corner values
		Assert(drawSpecs.overrideRect.left is: 0)
		Assert(drawSpecs.overrideRect.top is: 0)
		Assert(drawSpecs.overrideRect.right is: 300)
		Assert(drawSpecs.overrideRect.bottom is: 20)
		// Draw override bottom rounded corner fill values
		Assert(drawSpecs.overrideFill.left is: 1)
		Assert(drawSpecs.overrideFill.top is: 19)
		Assert(drawSpecs.overrideFill.right is: 299)
		Assert(drawSpecs.overrideFill.bottom is: 20)

		// Test text draw positions
		// Text draw position for selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, true)
		Assert(textSpecs.x is: 11)
		Assert(textSpecs.y is: -6)
		// Text draw position for not selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 13)
		Assert(textSpecs.y is: -6)
		// Text draw position for not selected tab with renderWidth differing from width
		tab.width += 2
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 11)
		Assert(textSpecs.y is: -6)
		tab.width -= 2

		// Test image rectangle values
		imageRect = calcCl.ImageRect(tab)
		Assert(imageRect.left is: 7)
		Assert(imageRect.top is: -6)
		Assert(imageRect.right is: 17)
		Assert(imageRect.bottom is: 4)

		// Test tab underline positions
		tabPoints = calcCl.CalcLinePoints()
		Assert(tabPoints.x1 is: 0)
		Assert(tabPoints.y1 is: 0)
		Assert(tabPoints.x2 is: 199)
		Assert(tabPoints.y2 is: 0)

		// Test invalidate rectangle values
		invalidateOb = calcCl.InvalidateRect(0, tab)
		Assert(invalidateOb.left is: 0)
		Assert(invalidateOb.top is: 0)
		Assert(invalidateOb.right is: 30)
		Assert(invalidateOb.bottom is: 10)
		}

	Test_Resize?()
		{
		.spy()
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'bottom')

		// .W is initialized to false
		Assert(calcCl.Resize?(false, 'unused - h') is: false)
		Assert(calcCl.Resize?(0, 'unused - h'))

		calcCl.Resize(200, 20)
		Assert(calcCl.Resize?(0, 'unused - h'))
		Assert(calcCl.Resize?(199, 'unused - h'))
		Assert(calcCl.Resize?(200, 'unused - h') is: false)
		Assert(calcCl.Resize?(201, 'unused - h'))
		}

	Test_ResizeExtraControl()
		{
		.spy()
		extraCtrl = Mock()
		extraCtrl.When.Resize([anyArgs:]).Do({ })

		// Horizontal tabs top
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'top')
		calcCl.Resize(200, 20)
		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50)
		extraCtrl.Verify.Resize(5, 0, 50, 19)

		// Horizontal tabs bottom
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'bottom')
		calcCl.Resize(200, 20)
		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50)
		extraCtrl.Verify.Resize(5, 1, 50, 19)
		}

	Test_ResizeButton()
		{
		.spy()
		button = Mock()
		button.When.Resize([anyArgs:]).Do({ })

		// Horizontal tabs top
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'top')
		calcCl.Resize(200, 20)
		calcCl.ResizeButton(button, 10)
		button.Verify.Resize(10, 5, 4, 4)

		// Horizontal tabs bottom
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'bottom')
		calcCl.Resize(200, 20)
		calcCl.ResizeButton(button, 10)
		button.Verify.Resize(10, 1, 4, 4)
		}

	Test_TabDragSpecs()
		{
		.spy()
		calcCl = TabHorzCalcs(new .fakeTabControl, orientation: 'bottom')
		calcCl.Resize(w = 200, 20)
		Assert(calcCl.TabBarSize is: w)

		// availableSpace is checked against x, .H is checked against y
		// x < 0 (outside of bar range), y is between 0 and 200 (inside of bar range)
		result = calcCl.TabDragSpecs(x = -1, 25, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.HSPLITBAR)
		Assert(result.check is: x)

		// x is 0 (just inside of bar range), y > 200 (just oustide of bar range)
		result = calcCl.TabDragSpecs(x = 0, 21, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.HSPLITBAR)
		Assert(result.check is: x)

		// x is 0 (just inside of bar range), y is 200 (just inside of bar range)
		result = calcCl.TabDragSpecs(x = 0, 20, 200)
		Assert(result.drag?)
		Assert(result.cursor is: IDC.HSPLITBAR)
		Assert(result.check is: x)

		// x is 20 (middle of bar range), y is 10 (middle of bar range)
		result = calcCl.TabDragSpecs(x = 100, 10, 200)
		Assert(result.drag?)
		Assert(result.cursor is: IDC.HSPLITBAR)
		Assert(result.check is: x)

		// x is 201 (just outside of bar range), y is 10 (middle of bar range)
		result = calcCl.TabDragSpecs(x = 201, 10, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.HSPLITBAR)
		Assert(result.check is: x)
		}
	}
