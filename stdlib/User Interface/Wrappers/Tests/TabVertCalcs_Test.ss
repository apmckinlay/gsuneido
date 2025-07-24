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

	Test_left()
		{
		.spy()
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'left')

		calcCl.Resize(20, h = 200)
		Assert(calcCl.TabBarSize is: h)
		Assert(calcCl.FontOrientation is: 900)
		Assert(calcCl.TabVertCalcs_extraControlX is: -1)
		Assert(calcCl.TabVertCalcs_buttonOffset is: 0)

		// Tab render values
		.setTextMetrics(tab = Object(tabName: 'Tab1', hide?: false, data: false, image:))
		calcCl.CalcRenderRect(0, tab, 0)
		Assert(tab.renderRect.start is: 0)
		Assert(tab.renderRect.top is: 0)
		Assert(tab.renderRect.tipY is: 0)
		Assert(tab.renderRect.end is: 30)
		Assert(tab.renderRect.bottom is: 30)
		Assert(tab.renderRect.left is: 1)
		Assert(tab.renderRect.right is: 10)
		Assert(tab.renderRect.tipX is: 10)
		Assert(tab.renderWidth is: tab.width)

		// Tab draw values
		drawSpecs = calcCl.CalcDrawSpecs(300, 80)
		Assert(drawSpecs.ellipseSize is: 60)
		// Draw base round rectangle values
		Assert(drawSpecs.baseRound.left is: 0)
		Assert(drawSpecs.baseRound.top is: 0)
		Assert(drawSpecs.baseRound.right is: 300)
		Assert(drawSpecs.baseRound.bottom is: 80)
		// Draw base round rectangle fill values
		Assert(drawSpecs.baseFill.right is: 300)
		Assert(drawSpecs.baseFill.bottom is: 80)
		// Draw override bottom rounded corner values
		Assert(drawSpecs.overrideRect.left is: 30)
		Assert(drawSpecs.overrideRect.top is: 0)
		Assert(drawSpecs.overrideRect.right is: 300)
		Assert(drawSpecs.overrideRect.bottom is: 80)
		// Draw override bottom rounded corner fill values
		Assert(drawSpecs.overrideFill.left is: 31)
		Assert(drawSpecs.overrideFill.top is: 1)
		Assert(drawSpecs.overrideFill.right is: 30)
		Assert(drawSpecs.overrideFill.bottom is: 79)

		// Test text draw positions
		// Text draw position for selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, true)
		Assert(textSpecs.x is: 4)
		Assert(textSpecs.y is: 19)
		// Text draw position for not selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 4)
		Assert(textSpecs.y is: 17)
		// Text draw position for not selected tab with renderWidth differing from width
		tab.width += 2
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 4)
		Assert(textSpecs.y is: 19)
		tab.width -= 2

		// Test image rectangle values
		imageRect = calcCl.ImageRect(tab)
		Assert(imageRect.left is: 4)
		Assert(imageRect.top is: 13)
		Assert(imageRect.right is: 14)
		Assert(imageRect.bottom is: 23)

		// Test tab underline positions
		tabPoints = calcCl.CalcLinePoints()
		Assert(tabPoints.x1 is: 9)
		Assert(tabPoints.y1 is: 7)
		Assert(tabPoints.x2 is: 9)
		Assert(tabPoints.y2 is: h)

		// Test invalidate rectangle values
		invalidateOb = calcCl.InvalidateRect(0, tab)
		Assert(invalidateOb.left is: 1)
		Assert(invalidateOb.top is: -1)
		Assert(invalidateOb.right is: 10)
		Assert(invalidateOb.bottom is: 31)
		}

	Test_right()
		{
		.spy()
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'right')

		calcCl.Resize(20, h = 200)
		Assert(calcCl.TabBarSize is: h)
		Assert(calcCl.FontOrientation is: 2700)
		Assert(calcCl.TabVertCalcs_extraControlX is: 1)
		Assert(calcCl.TabVertCalcs_buttonOffset is: -2)

		// Tab render values
		.setTextMetrics(tab = Object(tabName: 'Tab1', hide?: false, data: false, image:))
		calcCl.CalcRenderRect(0, tab, 0)
		Assert(tab.renderRect.start is: 0)
		Assert(tab.renderRect.top is: 0)
		Assert(tab.renderRect.tipY is: 0)
		Assert(tab.renderRect.end is: 31)
		Assert(tab.renderRect.bottom is: 31)
		Assert(tab.renderRect.left is: 0)
		Assert(tab.renderRect.right is: 9)
		Assert(tab.renderRect.tipX is: 9)
		Assert(tab.renderWidth is: tab.width)

		// Tab draw values
		drawSpecs = calcCl.CalcDrawSpecs(300, 80)
		Assert(drawSpecs.ellipseSize is: 60)
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
		Assert(drawSpecs.overrideRect.right is: 275)
		Assert(drawSpecs.overrideRect.bottom is: 80)
		// Draw override bottom rounded corner fill values
		Assert(drawSpecs.overrideFill.left is: -1)
		Assert(drawSpecs.overrideFill.top is: 1)
		Assert(drawSpecs.overrideFill.right is: 275)
		Assert(drawSpecs.overrideFill.bottom is: 79)

		// Test text draw positions
		// Text draw position for selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, true)
		Assert(textSpecs.x is: 6)
		Assert(textSpecs.y is: 11)
		// Text draw position for not selected tab
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 6)
		Assert(textSpecs.y is: 13)
		// Text draw position for not selected tab with renderWidth differing from width
		tab.width += 2
		textSpecs = calcCl.CalcTextSpecs(tab, false)
		Assert(textSpecs.x is: 6)
		Assert(textSpecs.y is: 11)
		tab.width -= 2

		// Test image rectangle values
		imageRect = calcCl.ImageRect(tab)
		Assert(imageRect.left is: -4)
		Assert(imageRect.top is: 7)
		Assert(imageRect.right is: 6)
		Assert(imageRect.bottom is: 17)

		// Test tab underline positions
		tabPoints = calcCl.CalcLinePoints()
		Assert(tabPoints.x1 is: 0)
		Assert(tabPoints.y1 is: 0)
		Assert(tabPoints.x2 is: 0)
		Assert(tabPoints.y2 is: h)

		// Test invalidate rectangle values
		invalidateOb = calcCl.InvalidateRect(0, tab)
		Assert(invalidateOb.left is: 0)
		Assert(invalidateOb.top is: -1)
		Assert(invalidateOb.right is: 10)
		}

	Test_Resize?()
		{
		.spy()
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'left')

		// .H is initialized to false
		Assert(calcCl.Resize?('unused - w', false) is: false)
		Assert(calcCl.Resize?('unused - w', 0))

		calcCl.Resize(20, 200)
		Assert(calcCl.Resize?('unused - w', 199))
		Assert(calcCl.Resize?('unused - w', 200) is: false)
		Assert(calcCl.Resize?('unused - w', 201))
		}

	Test_ResizeExtraControl()
		{
		.spy()
		extraCtrl = Mock()
		extraCtrl.When.Resize([anyArgs:]).Do({ })
		programmerErrors = .SpyOn(ProgrammerError).Return('').CallLogs()

		// Vertical tabs left
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'left')
		calcCl.Resize(20, 200)
		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50)
		extraCtrl.Verify.Resize(-1, 5, 13, 50)
		Assert(programmerErrors isSize: 0)

		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50, xstretch:)
		Assert(programmerErrors isSize: 1)

		// Vertical tabs right
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'right')
		calcCl.Resize(20, 200)
		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50)
		extraCtrl.Verify.Resize(1, 5, 13, 50)
		Assert(programmerErrors isSize: 1)

		calcCl.ResizeExtraControl(extraCtrl, ctrlPos: 5, ctrlSize: 50, xstretch:)
		Assert(programmerErrors isSize: 2)
		}

	Test_ResizeButton()
		{
		.spy()
		button = Mock()
		button.When.Resize([anyArgs:]).Do({ })

		// Vertical tabs left
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'left')
		calcCl.Resize(20, 200)
		calcCl.ResizeButton(button, 10)
		button.Verify.Resize(5, 10, 4, 4)

		// Vertical tabs right
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'right')
		calcCl.Resize(20, 200)
		calcCl.ResizeButton(button, 10)
		button.Verify.Resize(3, 10, 4, 4)
		}

	Test_TabDragSpecs()
		{
		.spy()
		calcCl = TabVertCalcs(new .fakeTabControl, orientation: 'left')
		calcCl.Resize(20, h = 200)
		Assert(calcCl.TabBarSize is: h)

		// .W is checked against x, availableSpace is checked against y
		// x < 0 (outside of bar range), y is between 0 and 200 (inside of bar range)
		result = calcCl.TabDragSpecs(-1, y = 25, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.VSPLITBAR)
		Assert(result.check is: y)

		// x is 0 (just inside of bar range), y > 200 (just oustide of bar range)
		result = calcCl.TabDragSpecs(0, y = 201, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.VSPLITBAR)
		Assert(result.check is: y)

		// x is 0 (just inside of bar range), y is 200 (just inside of bar range)
		result = calcCl.TabDragSpecs(0, y = 200, 200)
		Assert(result.drag?)
		Assert(result.cursor is: IDC.VSPLITBAR)
		Assert(result.check is: y)

		// x is 20 (just inside of bar range), y is 100 (middle of bar range)
		result = calcCl.TabDragSpecs(20, y = 100, 200)
		Assert(result.drag?)
		Assert(result.cursor is: IDC.VSPLITBAR)
		Assert(result.check is: y)

		// x is 21 (just outside of bar range), y is 100 (middle of bar range)
		result = calcCl.TabDragSpecs(21, y = 100, 200)
		Assert(result.drag? is: false)
		Assert(result.cursor is: IDC.VSPLITBAR)
		Assert(result.check is: y)
		}
	}
