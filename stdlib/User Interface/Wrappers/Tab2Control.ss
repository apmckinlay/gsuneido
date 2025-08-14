// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 'Tab'
	Xstretch: 1
	New(@tabs)
		{
		.CreateWindow('SuBtnfaceArrow', '', WS.VISIBLE, exStyle: WS_EX.CONTROLPARENT)
		.SubClass()
		.initDefaults(tabs)
		.initExtraControl(tabs)
		.initBrushes()
		.initTabItems(tabs)
		.addTabButton(tabs)
		.initToolTips()
		}

	calcClass: false
	initDefaults(tabs)
		{
		// Initialize generic instance members
		.staticTabs = tabs.GetDefault(#staticTabs, #())
		.closeImage = tabs.GetDefault(#close_button, false) isnt false
			? Object(ImageResource(#close), CLR.BLACK)
			: false
		.fitTabs = tabs.GetDefault(#scrollTabs, false)
			? .scrollTabs
			: .shrinkTabs
		if false is .selectedTabColor = tabs.GetDefault(#selectedTabColor, false)
			.selectedTabColor = CLR.Highlight

		// Initialize the TabCalcs class
		orientation = tabs.orientation
		selectedWeight = tabs.GetDefault(#selectedTabBold, true) ? FW.BOLD : ''
		.calcClass = orientation in (#left, #right)
			? new TabVertCalcs(controller: this, :orientation, :selectedWeight)
			: new TabHorzCalcs(controller: this, :orientation, :selectedWeight)
		}

	initExtraControl(tabs)
		{
		ctrl = tabs.GetDefault(#extraControl, false)
		.extraControl = ctrl isnt false
			? .Construct(ctrl)
			: class // Fake control to simplify references / use
				{
				Ymin: 0
				Xmin: 0
				Default(@unused) { }
				GetChildren() { return Object() }
				}
		.Ymin = Max(.Ymin, .extraControl.Ymin + 1 /*= bottom line*/)
		}

	initTabItems(tabs)
		{
		i = 0
		tab = false
		for tabName in tabs.Values(list:)
			.tabItems.Add(tab = .initTab(i++, tabName, prevTab: tab))
		}

	getter_tabItems()
		{
		return .tabItems = Object()
		}

	initTab(i, tabName, data = false, image = -1, prevTab = false)
		{
		if data is false
			data = Object(tooltip: '')
		tab = Object(hide?: false, :data, image: .baseImage(tabName, image))
		if prevTab is false
			prevTab = .prevTab(i)
		.calcTabItem(i, tab, prevTab.renderRect.end, tabName)
		return tab
		}

	baseImage(tabName, image)
		{
		if image is -1
			image = false
		baseImage = .staticTabs.Has?(tabName) ? image : .closeImage
		return image isnt false and .imageList isnt false
			? .imageList[image]
			: baseImage
		}

	prevTab(i)
		{
		return .tabItems.GetDefault(i - 1, #(renderRect: (end: 0)))
		}

	calcTabItem(i, tab, end, tabName = false)
		{
		return .calcClass.CalcRenderRect(i, tab, end, tabName)
		}

	initBrushes()
		{
		.brushes[CLR.white] = .selectedBgBrush = CreateSolidBrush(CLR.white)
		.brushes[CLR.lightblue] = .selectedTopBrush = CreateSolidBrush(CLR.lightblue)
		.brushes[CLR.highlightblue] = .hoverBgBrush = CreateSolidBrush(CLR.highlightblue)
		.brushes[CLR.ButtonFace] = .unselectedBgBrush = CreateSolidBrush(CLR.ButtonFace)
		.brushes[COLOR.TRIDSHADOW] = .borderPen =
			CreatePen(PS.SOLID, 0, GetSysColor(COLOR.TRIDSHADOW))
		}

	getter_brushes()
		{
		return .brushes = Object()
		}

	initToolTips()
		{
		.Map = Object()
		.Map[TTN.SHOW] = 'TTN_SHOW'
		.tip = .Construct(ToolTipControl)
		.tip.SendMessage(TTM.SETMAXTIPWIDTH, 0, 400) /*= tip max width */
		.tip.SendMessage(TTM.SETDELAYTIME, TTDT.AUTOPOP, 30000) /*= tip delay time*/
		.tip.Activate(false)
		.tip.AddTool(.Hwnd, LPSTR_TEXTCALLBACK)
		.tip.SetFont(StdFonts.Mono())
		.SetRelay(.tip.RelayEvent)
		}

	tabButton: false
	defaultTip: 'Add Tab'
	addTabButton(args)
		{
		buttonTip = args.GetDefault(#buttonTip, .defaultTip)
		if args.GetDefault(#addTabButton?, false) or buttonTip isnt .defaultTip
			.tabButton = .constructButton('expand.emf', buttonTip)
		}

	constructButton(image, tip, target = false)
		{
		button = .Construct(Object('EnhancedButton', :image, command: tip,
			imagePadding: 0.15, mouseEffect:, :tip))
		if target
			button.SetCommandTarget(this)
		return button
		}

	w: false
	h: false
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		if .selectedIdx isnt -1
			.ensureVisibleIdx = .selectedIdx
		if .calcClass.Resize?(w, h)
			.calcRender(w, h)
		}

	calcRender(.w, .h)
		{
		.calcClass.Resize(w, h)
		if .requiredSpace() >= tabBarSize = .calcClass.TabBarSize
			{
			(.fitTabs)(tabBarSize)
			prevEnd = 0
			.iterateTabItems()
				{|unused, tab|
				if tab.hide?
					continue
				tab.renderRect.start = prevEnd
				tab.renderRect.end = prevEnd + tab.renderWidth + 1
				prevEnd += tab.renderWidth
				}
			}
		else if .navigationButtons.NotEmpty?()
			.destroyNavigationButtons()
		.resizeControls(tabBarSize)
		}

	requiredSpace()
		{
		prevEnd = 0
		.iterateTabItems({|i, tab| prevEnd = .calcTabItem(i, tab, prevEnd) })
		return prevEnd + .reservedSpace()
		}

	iterateTabItems(block)
		{
		Seq(.Count()).Each()
			{|i|
			block(i, .tabItems[i])
			}
		}

	reservedSpace()
		{
		preAllocated = .extraControl.Xmin + .calcClass.PaddingSide
		if .tabButton isnt false
			preAllocated  += .calcClass.ButtonSize + .calcClass.PaddingSide
		return .navButtonSize() + preAllocated
		}

	navButtonSize()
		{
		return .navigationButtons.Size() * .calcClass.ButtonSize
		}

	getter_navigationButtons()
		{
		return .navigationButtons = Object()
		}

	shrinkTabs(w)
		{
		if .navigationButtons.Empty?()
			.navigationButtons.Add(.constructButton('arrow_down', 'Go to Tab', target:))
		.shrink(.availableSpace(w))
		}

	availableSpace(w)
		{
		return w - .reservedSpace()
		}

	shrink(availableSpace)
		{
		items = .tabItems.
			Map2({|m, v| [idx: m, width: v.width] }).
				Filter({ it.width > 0 }).
				Add([idx: false, width: 0])
		over = .tabItems.SumWith({ it.width }) - availableSpace
		rank = items.Members().Sort!({|x, y| items[x].width >= items[y].width })
		for (i = 1; i < rank.Size() and over > 0; i++)
			{
			diff = items[rank[i - 1]].width - items[rank[i]].width
			sub = (diff * i <= over) ? diff : (over / i).Ceiling()
			for (j = 0; j < i; j++)
				{
				.tabItems[items[rank[j]].idx].renderWidth -= sub
				over -= sub
				}
			}
		}

	ensureVisibleIdx: 0
	scrollTabs(w)
		{
		if .navigationButtons.Empty?()
			.navigationButtons.Add(
				.constructButton(.calcClass.ScrollNextImage, 'Next', target:, navButton:),
				.constructButton(.calcClass.ScrollPrevImage, 'Previous',
					target:, navButton:))
		.scroll(.availableSpace(w))
		}

	scroll(availableSpace)
		{
		prevEnd = 0
		firstTab = .firstTab(availableSpace)
		.iterateTabItems()
			{|i, tab|
			imageWidth = .calcClass.ImageWidth(tab.image)
			if tab.hide? = .hideTab?(i, firstTab, prevEnd, imageWidth, availableSpace)
				{
				// overwrite renderRect to ensure hidden tabs do NOT overlap visible tabs
				tab.renderRect = Object().Set_default(0)
				continue
				}
			.calcClass.CalcRenderRect(i, tab, prevEnd)
			if tab.renderRect.end > availableSpace
				tab.renderWidth -= (tab.renderRect.end - availableSpace)
			prevEnd += tab.renderWidth
			}
		}

	firstTab(availableSpace)
		{
		ensureVisibleIdx = Min(.ensureVisibleIdx, .lastTabIdx())
		ensureVisibleTab = .tabItems[ensureVisibleIdx]
		if availableSpace - ensureVisibleTab.renderWidth < ensureVisibleTab.width
			return ensureVisibleIdx
		range = ensureVisibleTab.renderRect.end - availableSpace
		firstTab = .tabItems.FindIf({ it.renderRect.start > range })
		return firstTab is false ? ensureVisibleIdx : firstTab
		}

	lastTabIdx()
		{
		return .Count() - 1
		}

	hideTab?(i, firstTab, start, imageWidth, availableSpace)
		{
		return i >= firstTab
			? 2 * imageWidth + start > availableSpace
			: true
		}

	destroyNavigationButtons()
		{
		while false isnt button = .navigationButtons.Extract(0, false)
			button.Destroy()
		.tabItems.Each({ it.hide? = false })
		}

	resizeControls(tabBarSize)
		{
		pos = .resizeTabButton()
		pos = .resizeExtraControl(tabBarSize, pos + .navButtonSize())
		.resizeTabButtons(pos)
		}

	resizeTabButton() // Placed to the immediate right of the last tab
		{
		pos = .lastVisibleTabEnd()
		if .tabButton isnt false
			pos += .calcClass.ResizeButton(.tabButton, pos + .calcClass.PaddingSide)
		return pos
		}

	lastVisibleTabEnd()
		{
		tab = false
		.tabItems.Each()
			{
			if not it.hide?
				tab = it
			else if tab isnt false
				return tab.renderRect.end
			}
		return tab isnt false ? tab.renderRect.end : 0
		}

	resizeExtraControl(tabBarSize, pos) // Placed at the very end of the tab bar
		{
		ctrlSize = .extraControl.Xmin
		ctrlPos = tabBarSize - ctrlSize
		xstretch = .extraControl.GetChildren().Any?({ it.GetDefault(#Xstretch, 0) >= 0 })
		if xstretch
			ctrlSize = tabBarSize - ctrlPos = pos
		.calcClass.ResizeExtraControl(.extraControl, ctrlPos, ctrlSize, :xstretch)
		return ctrlPos
		}

	resizeTabButtons(pos) // Placed between the last tab and the .extraControl
		{
		pos -= .calcClass.ButtonSize
		.navigationButtons.Each({ pos -= .calcClass.ResizeButton(it, pos) })
		}

	hoverIdx: 		false
	selectedIdx: 	-1
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		if .tabItems.NotEmpty?()
			.paint(hdc)
		EndPaint(.Hwnd, ps)
		return 0
		}

	paint(hdc)
		{
		if .selectedIdx is -1
			.selectTab(0)
		else
			WithHdcSettings(hdc, [.borderPen, SetBkMode: TRANSPARENT])
				{
				DoWithHdcObjects(hdc, [.calcClass.Font()])
					{
					.iterateTabItems()
						{|i, tab|
						if i isnt .selectedIdx
							.drawTab(hdc, i, tab)
						}
					}
				// Draw the selected last to avoid the antialias effect from drawing
				// other tabs the selected tab is higher than the unselected tabs.
				// Drawing the unselected tabs after the selected tab causes a small
				// overlap which looks like a gap
				DoWithHdcObjects(hdc, [.calcClass.Font(selectedTab?:)])
					{
					.drawTab(hdc, .selectedIdx, .tabItems[.selectedIdx], selectedTab:)
					}
				.drawBottomLine(hdc)
				}
		}

	drawTab(hdc, i, tab, selectedTab = false)
		{
		if tab.hide?
			return
		.drawRect(i, hdc, tab.renderRect, selectedTab)
		.drawText(hdc, tab, selectedTab)
		.drawImage(hdc, tab)
		}

	drawRect(i, hdc, rect, selectedTab)
		{
		brush = selectedTab
			? .selectedBgBrush
			: .hoverIdx is i
				? .hoverBgBrush
				: .unselectedBgBrush
		rc = Rect(rect.left, rect.top, rect.right - rect.left, rect.bottom - rect.top)
		Image.PaintWithAntialias(hdc, rc.GetWidth(), rc.GetHeight(), rc)
			{|hdcLarge, wLarge, hLarge|
			.drawRoundRect(hdcLarge, wLarge, hLarge, brush)
			}
		}

	drawRoundRect(hdc, wLarge, hLarge, brush)
		{
		DoWithHdcObjects(hdc, [brush])
			{
			// Calculate the required rectangles based on the calcClass
			specs = .calcClass.CalcDrawSpecs(wLarge, hLarge)

			// Draw the rounded rectangle
			FillRect(hdc, specs.baseFill, .unselectedBgBrush)
			.roundRect(hdc, specs.baseRound, specs.ellipseSize, specs.ellipseSize)

			// Override the bottom round corners with squared corners
			.roundRect(hdc, specs.overrideRect)

			// Erase the extra line from the previous step
			FillRect(hdc, specs.overrideFill, brush)
			}
		}

	roundRect(hdc, rect, width = 0, height = 0)
		{
		RoundRect(hdc, rect.left, rect.top, rect.right, rect.bottom, width, height)
		}

	drawBottomLine(hdc)
		{
		points = .calcClass.CalcLinePoints()
		MoveTo(hdc, points.x1, points.y1)
		LineTo(hdc, points.x2, points.y2)
		}

	drawText(hdc, tab, selectedTab)
		{
		if false is text = .calcClass.CalcTextSpecs(tab, selectedTab)
			return
		hdcSettings = Object()
		if selectedTab
			hdcSettings.SetTextColor = .selectedTabColor
		WithHdcSettings(hdc, hdcSettings)
			{
			TextOut(hdc, text.x, text.y, text.text, text.text.Size())
			}
		}

	drawImage(hdc, tab)
		{
		if tab.image is false
			return
		imageSpecs = .calcClass.CalcImageSpecs(tab)
		brushImage = .brush(tab.image.GetDefault(1, false))
		tab.image[0].Draw(hdc, imageSpecs.x, imageSpecs.y,
			imageSpecs.size, imageSpecs.size, :brushImage,
			orientation: .calcClass.FontOrientation)
		}

	brush(clr)
		{
		return clr isnt false ? .brushes.GetInit(clr, { CreateSolidBrush(clr) }) : false
		}

	MOUSEMOVE(lParam)
		{
		i = .getTabFromXY(x = LOSWORD(lParam), y = HISWORD(lParam))
		if .draggedIdx isnt false
			.draggingTab(i, x, y)
		else if i isnt false and .hoverIdx isnt i
			.hoverNewTab(.hoverIdx, i)
		return 0
		}

	getTabFromXY(x, y)
		{
		return .tabItems.FindIf({
			x >= it.renderRect.left and x <= it.renderRect.right and
			y >= it.renderRect.top and y <= it.renderRect.bottom })
		}

	draggingTab(i, x, y)
		{
		dragSpecs = .calcClass.TabDragSpecs(x, y, .availableSpace(.calcClass.TabBarSize))
		.attemptDrag(i, dragSpecs.check)
		cursor = .separateTabs?
			? dragSpecs.drag? ? dragSpecs.cursor : IDC.DRAG1COPY
			: dragSpecs.cursor
		SetCursor(LoadCursor(ResourceModule(), cursor))
		}

	attemptDrag(i, check)
		{
		if not .dragRequired?(i, check, draggedTab = .tabItems[.draggedIdx])
			return false
		previous? = check <= (draggedTab.renderRect.start + draggedTab.renderRect.end) / 2
		// i can be false when dragging outside of the tab bar or over "vacant" space
		// - IE: In between the last tab and the tab navigation buttons
		if i is false
			i = previous?
				? Max(.draggedIdx - 1, 0)
				: Min(.draggedIdx + 1, .lastTabIdx())
		if dragTab? = .dragTab?(check, previous?, .tabItems[i], draggedTab)
			.dragTab(i, previous?)
		return dragTab?
		}

	dragRequired?(i, check, draggedTab)
		{
		return i isnt .draggedIdx and
			check < draggedTab.renderRect.start or
			check > draggedTab.renderRect.end
		}

	dragTab?(check, previous?, nextTab, draggedTab)
		{
		return previous?
			? nextTab.renderRect.start > check - draggedTab.renderWidth
			: nextTab.renderRect.end < check + draggedTab.renderWidth
		}

	dragTab(i, previous?)
		{
		range = .visibleTabRange()
		if .draggedIdx is range.first and range.first isnt 0
			if not .tabFullyVisible?(.tabItems[.ensureVisibleIdx = range.last])
				--.ensureVisibleIdx
		.Send(#MoveTab, .draggedIdx, i)
		if not .tabFullyVisible?(.tabItems[.draggedIdx = i])
			this[#On_ $ (previous? ? #Previous : #Next)]()
		}

	tabFullyVisible?(tab)
		{
		return not tab.hide? and tab.renderWidth is tab.width
		}

	getter_separateTabs?()
		{
		result = .Send(#Collect, #Tab_SupportSeparate?)
		return .separateTabs? = Object?(result) ? result.Any?({ it is true }) : false
		}

	hoverNewTab(oldHover, newHover)
		{
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE, hwndTrack: .Hwnd))
		.restoreTabImage(oldHover)
		.setCloseImage(.hoverIdx = newHover)
		.repaintRect(newHover)
		.setTipText()
		}

	restoreTabImage(i)
		{
		if not .tabItems.Member?(i)
			return
		tab = .tabItems[i]
		tab.image = tab.Extract(#prevImage, tab.image)
		.repaintRect(i)
		}

	repaintRect(i)
		{
		if .tabItems.Member?(i) and .tabItems[i].Member?(#renderRect)
			InvalidateRect(.Hwnd, .calcClass.InvalidateRect(i, .tabItems[i]) true)
		}

	setCloseImage(i)
		{
		tab = .tabItems[i]
		if .closeImage is false or tab.image is .closeImage or
			.staticTabs.Has?(tab.tabName)
			return
		tab.prevImage = tab.image
		tab.image = .closeImage
		}

	setTipText()
		{
		tab = .tabItems[.hoverIdx]
		tip = tab.GetDefault(#data, []).GetDefault(#tooltip, '')
		if tip is '' and not .tabFullyVisible?(tab)
			tip = tab.tabName
		.tip.UpdateTipText(.Hwnd, tip)
		.tip.Activate(true)
		}

	MOUSELEAVE()
		{
		.restoreTabImage(.hoverIdx)
		.hoverIdx = false
		.tip.Activate(false)
		return 0
		}

	LBUTTONDOWN(lParam)
		{
		if false isnt i = .getTabFromXY(x = LOSWORD(lParam), y = HISWORD(lParam))
			.tabClicked(i, x, y)
		return 0
		}

	tabClicked(i, x, y)
		{
		if .closeTab?(i, x, y)
			{
			.Send(#Tab_Close, i)
			.hoverIdx = false
			return
			}
		.goto(i)
		.dragStart(i)
		}

	closeTab?(i, x, y)
		{
		tab = .tabItems[i]
		if tab.image is false or tab.image isnt .closeImage
			return false
		posOb = .calcClass.ImageRect(tab)
		return x >= posOb.left and x <= posOb.right and
			y >= posOb.top and y <= posOb.bottom
		}

	goto(clicked)
		{
		if clicked is false or true is .Send(#TabControl_SelChanging)
			return 0

		.focusTab(clicked)
		if clicked is .selectedIdx
			.Send(#TabClick, .selectedIdx)
		else
			{
			.hoverIdx = false
			.selectTab(clicked)
			.Send(#SelectTab, .selectedIdx)
			}
		return 0
		}

	focusTab(i)
		{
		tab = .tabItems[i]
		if tab.hide?
			{
			.ensureVisibleIdx = i
			.repaintTabs()
			}
		else if tab.renderWidth isnt tab.width
			.On_Next()
		}

	selectTab(newIdx)
		{
		if newIdx is .selectedIdx
			return
		oldIdx = .selectedIdx
		.recalcTabSelect(.selectedIdx = newIdx)
		.recalcTabSelect(oldIdx)
		.Repaint()
		}

	recalcTabSelect(i)
		{
		if false isnt tab = .tabItems.GetDefault(i, false)
			.calcClass.CalcSelectChange(i, tab)
		}

	dragStart(i)
		{
		if not .allowDraggingTabs?
			return
		.draggedIdx = i
		SetCapture(.Hwnd)
		}

	getter_allowDraggingTabs?()
		{
		return .allowDraggingTabs? = .Send(#Tab_AllowDrag) is true
		}

	LBUTTONUP(lParam)
		{
		if .draggedIdx isnt false
			.dragEnd(lParam, .draggedIdx)
		return 0
		}

	draggedIdx: false
	dragEnd(lParam, draggedIdx)
		{
		.draggedIdx = false // Prevent futher calls to: .draggingTab
		ReleaseCapture()
		SetCursor(LoadCursor(ResourceModule(), IDC.ARROW))
		dragSpecs = .calcClass.TabDragSpecs(LOSWORD(lParam), HISWORD(lParam),
			.availableSpace(.calcClass.TabBarSize))
		if not dragSpecs.drag?
			.Send(#Tab_Separate, draggedIdx)
		}

	TTN_SHOW(lParam)
		{
		if .hoverIdx is false
			return true
		tab = .tabItems[.hoverIdx]
		.tip.AdjustRect(false, tab.renderRect.Copy())
		ClientToScreen(.Hwnd, p = [x: tab.renderRect.tipX, y: tab.renderRect.tipY])
		SetWindowPos(NMHDR(lParam).hwndFrom, 0,
			p.x, p.y, 0, 0, // rect
			SWP.NOACTIVATE | SWP.NOSIZE | SWP.NOZORDER)
		return true
		}

	Insert(i, tabName, data = false, image = -1)
		{
		.tabItems.Add(.initTab(i, tabName, data, image), at: i)
		.repaintTabs()
		}

	repaintTabs()
		{
		.calcRender(.w, .h)
		.Repaint()
		}

	SetData(i, data)
		{
		.tabItems[i].data = data
		}

	SetText(i, text)
		{
		tab = .tabItems[i]
		if tab.tabName is text
			return
		.calcTabItem(i, tab, tab.renderRect.start, text)
		.repaintTabs()
		}

	GetData(i)
		{
		return .tabItems[i].GetDefault(#data, [])
		}

	Remove(i)
		{
		.remove(i)
		.repaintTabs()
		}

	remove(i)
		{
		if i < .selectedIdx or .selectedIdx is .lastTabIdx()
			.selectedIdx--
		if i is .hoverIdx
			.hoverIdx = false
		.tabItems.Delete(i)
		}

	Select(i)
		{
		.focusTab(i)
		.selectTab(i)
		}

	imageList: false
	SetImageList(.imageList)
		{ }

	SetImage(i, image)
		{
		tab = .tabItems[i]
		image = .baseImage(tab.tabName, image)
		if i is .hoverIdx
			{
			tab.prevImage = image
			return
			}
		tab.Delete(#prevImage)
		repaint? = .calcClass.ImageWidth(image) isnt .calcClass.ImageWidth(tab.image)
		tab.image = image
		if repaint?
			.repaintTabs()
		else
			.repaintRect(i)
		}

	GetText(i)
		{
		return .tabItems[i].tabName
		}

	Count()
		{
		return .tabItems.Size()
		}

	GetSelected()
		{
		return .selectedIdx
		}

	ContextMenu(x, y)
		{
		return .hoverIdx is false ? 0 : .Send(#TabContextMenu, x, y, hover: .hoverIdx)
		}

	// Prevents tabs from being disabled
	SetReadOnly(unused)
		{ }

	// Read-only is not applicable to tabs
	GetReadOnly()
		{
		return true
		}

	On_Go_to_Tab()
		{
		GetCursorPos(pt = Object())
		list = .tabItems.Map({ it.tabName })
		if 0 isnt i = ContextMenu(list).Show(.Window.Hwnd, pt.x, pt.y)
			.goto(i - 1)
		}

	On_Previous()
		{
		if .tabFullyVisible?(.tabItems[0])
			return
		prevVisible = .ensureVisibleIdx
		if prevVisible isnt .ensureVisibleIdx = Max(.visibleTabRange().first - 1, 0)
			.repaintTabs()
		}

	visibleTabRange()
		{
		first = last = false
		.iterateTabItems()
			{|i, tab|
			if first is false and not tab.hide?
				first = i
			else if first isnt false and tab.hide?
				{
				last = i - 1
				break
				}
			}
		return [
			first: first is false ? 0 : first,
			last:   last is false ? .lastTabIdx() : last
			]
		}

	On_Next()
		{
		if .tabFullyVisible?(.tabItems[last = .lastTabIdx()])
			return
		prevVisibleIdx = .ensureVisibleIdx
		if prevVisibleIdx isnt .ensureVisibleIdx = Min(.visibleTabRange().last + 1, last)
			.repaintTabs()
		}

	Move(i, newPos)
		{
		.hoverIdx = false
		tab = .tabItems.Extract(i)
		tab.image = tab.Extract(#prevImage, tab.image)
		.calcTabItem(newPos, tab, .prevTab(newPos).renderRect.end)
		.selectedIdx = newPos
		if .draggedIdx is false
			.ensureVisibleIdx = .selectedIdx
		.tabItems.Add(tab, at: newPos)
		.repaintTabs()
		}

	ForEachTab(block)
		{
		.iterateTabItems({|i, tab| block(tab.GetDefault(#data, []), idx: i) })
		}

	Destroy()
		{
		.brushes.Each(DeleteObject)
		.brushes.Delete(all:)
		if .calcClass isnt false
			.calcClass.Destroy()
		if .tabButton isnt false
			.tabButton.Destroy()
		.extraControl.Destroy()
		.destroyNavigationButtons()
		.tip.Destroy()
		super.Destroy()
		}
	}