// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* USAGE:
Important Arguments / Class members:
- .Tab:		 	This will point to either Tab or Tab2, and is the literal "tabs"
- .ctrls: 		A object of constructed controls. This is what is under the tabs
- .controls: 	A object of "what" each tab is to consist of (layout and related data)
- .ctrl:		The current visible control (equivalent to: .ctrls[.Tab.GetSelected()])

WARNING:
- These objects MUST remain synced.
-- If not, the Destroy / Select will not behave as expected.

Example: TabsControl(#(Static 'this is one'  Tab: 'One') #(Editor, Tab: 'Editor'))
*/
PassthruController
	{
	Name: "Tabs"
	Xstretch: 1
	Ystretch: 1

	New(@args)
		{
		super(false)
		.controls = .controlArgs(args)
		.destroyOnSwitch = args.GetDefault(#destroyOnSwitch, false)
		.border = args.GetDefault(#border, 5 /*= default border size*/)
		.addCustomTabs(.controls, .customizable? = args.GetDefault(#customizable?, false))

		.Construct(.tabControl(.controls, args))

		.ctrls = Object().Set_default(false)
		.initialConstruct(args, startTab = args.GetDefault(#startTab, 0))
		.Recalc()
		.Send('AccessObserver', .AccessChanged)

		if not .NoCtrls?() and startTab isnt .GetSelected()
			.Tab.Select(startTab)
		}

	controlArgs(args)
		{
		controls = Object()
		for x in args.Values(list:)
			if not x.Member?(#Hide?) or x.Hide? isnt true
				controls.Add(x)
		return controls
		}

	addCustomTabs(controls, customizable?)
		{
		if not customizable?
			return

		if false is cl = OptContribution('CustomTabPermissions', false)
			return

		if not String?(tableName = .Send('GetCustomizableName'))
			tableName = .Send('GetTableName')

		cl.WithPermissableTabs(tableName)
			{|tab|
			controls.Add(Object('Customizable', Tab: tab, tabName: tab))
			}
		}

	customizableTabs: false
	tabControl(controls, args)
		{
		themed = args.GetDefault(#themed, true)
		close_button = args.GetDefault(#close_button, false)
		addTabButton? = args.GetDefault(#addTabButton?, false) is true
		buttonTip = args.GetDefault(#buttonTip, 'Add Tab')
		extraControl = args.GetDefault(#extraControl, false)
		staticTabs = args.GetDefault(#staticTabs, #())
		scrollTabs = args.GetDefault(#scrollTabs, false)
		selectedTabColor = args.GetDefault(#selectedTabColor, false)
		selectedTabBold = args.GetDefault(#selectedTabBold, true)
		orientation = .validateOrientation(args.GetDefault(#orientation, #top))
		.vertical = orientation in (#left, #right)
		.alternativePos = orientation in (#bottom, #right)
		.customizableTabs = Object()

		tabs = Object('Tab2', :addTabButton?, :buttonTip, :extraControl, :staticTabs,
			:scrollTabs, :selectedTabColor, :selectedTabBold, :orientation,
			:themed, :close_button)
		for x in controls
			{
			.collectCustomizableTabs(x, x.Tab)
			tabs.Add(x.Tab)
			}
		return tabs
		}

	validateOrientation(orientation)
		{
		// will throw if an invalid value is supplied for orientation
		switch (orientation)
			{
			case #top:
			case #left:
			case #right:
			case #bottom:
			}
		return orientation
		}

	collectCustomizableTabs(control, parentTab)
		{
		if not Object?(control)
			return
		if control.Has?('Customizable')
			.customizableTabs.Add(control.GetDefault('tabName', parentTab))
		else
			control.Filter(Object?).Each({ .collectCustomizableTabs(it, parentTab) })
		}

	initialConstruct(args, startTab)
		{
		for (i = args.Size(list:) - 1; i >= 0; --i)
			{
			if i is startTab
				{
				.construct(i)
				.ctrls[i].SetVisible(i is startTab)
				}
			else
				.ctrls[i] = false
			}
		.ctrl = .ctrls[startTab]
		}

	Recalc()
		{
		.Xmin = Max(.ctrlXmin + .tabOffset, .Tab.Xmin)
		.Ymin = .ctrlYmin + .Tab.Ymin + .tabOffset
		}

	getter_ctrlXmin()
		{ return .ctrl is false ? 0 : .ctrl.Xmin }

	getter_ctrlYmin()
		{ return .ctrl is false ? 0 : .ctrl.Ymin }

	NoCtrls?()
		{ return .ctrls.Empty?() or .ctrl is false }

	Remove(i)
		{
		.controls.Delete(i)
		.destroyTab(i)
		.Tab.Remove(i)
		}

	Append(label, layout, image = -1, noSelect = false)
		{
		.Insert(label, layout, :image, :noSelect)
		}

	Insert(label, layout, at = false, image = -1, data = false, noSelect = false)
		{
		i = at is false ? .Tab.Count() : at
		if data is false
			data = Object(tooltip: '')
		data.image = image
		.Tab.Insert(i, label, :data, :image)
		.controls.Add(layout, at: i)
		.ctrls.Add(false, at: i)
		if noSelect is false
			.Select(i)
		}

	On_NextTab()
		{
		if .Destroyed?()
			return
		if not .focusInTab?()
			{
			.FocusFirst(.ctrl.Hwnd)
			if .focusInTab?()
				return
			}
		.Select((.GetSelected() + 1) % .Tab.Count())
		}

	On_PrevTab()
		{
		if .Destroyed?()
			return
		n = .Tab.Count()
		.Select((.GetSelected() + n - 1) % n)
		}

	HasFocus?()
		{
		return .focusInTab?()
		}

	focusInTab?()
		{
		hwnd = GetFocus()
		while hwnd isnt NULL
			{
			if hwnd is .ctrl.Hwnd
				return true
			hwnd = GetParent(hwnd)
			}
		return false
		}

	GetSelected()
		{
		return .Tab.GetSelected()
		}

	GetTabData(idx = false, def = false)
		{
		idx = idx is false ? .Tab.GetSelected() : idx
		return idx < 0
			? Object().Set_default(def)
			: .Tab.GetData(idx)
		}

	SetTabData(idx, newData, name = false)
		{
		oldData = .GetTabData(idx)
		if oldData.Readonly?()
			oldData = oldData.Copy()
		if name isnt false
			{
			.controls[idx].Tab = name
			.Tab.SetText(idx, name)
			}
		.Tab.SetData(idx, oldData.Merge(newData))
		}

	Select(i, keepFocus = false)
		{
		if .GetSelected() is i and .ctrl is .ctrls[i] and .ctrl isnt false
			return
		.Tab.Select(i)
		.SelectTab(i, :keepFocus)
		}

	tabs_changing: false
	SelectTab(i, source = false, keepFocus = false)
		{
		if not .allowSelectTab?(source, i)
			return
		.tabs_changing = true
		.Send('SelectTab', i)
		if not keepFocus
			SetFocus(0)
		.deselectCtrl()
		if .ctrls[i] is false
			.ConstructTab(i)
		else
			.ctrls[i].SetVisible(true)
		.ctrl = .ctrls[i]
		.selectTabResize()
		.tabs_changing = false
		.Send('TabsControl_SelectTab')
		if not keepFocus
			.FocusFirst(.ctrl.Hwnd)
		}

	allowSelectTab?(source, i)
		{
		if .Destroyed?()
			return false
		if .ctrls.Empty?()
			return false
		if source isnt false and source isnt .Tab // nested tabs
			return false
		return false isnt .Send('AllowSelectTab', i)
		}

	deselectCtrl()
		{
		if .ctrl is false
			return
		if .destroyOnSwitch
			{
			.ctrl.Destroy()
			.ctrls.Replace(.ctrl, false)
			}
		else
			.ctrl.SetVisible(false)
		}

	x: false
	selectTabResize()
		{
		if .ctrl isnt false and .x isnt false
			.resize(.x, .y, .w, .h)
		if .scrollable?(this) or .Send('TabsControl_RefreshRequired?') is true
			.WindowRefresh()
		else
			{
			parentRect = GetWindowRect(.Parent.WindowHwnd())
			ctrlRect = GetWindowRect(.ctrl.Hwnd)
			.ensureWindowSize(parentRect, ctrlRect)
			// If FlowControl is the parent control, then a WindowRefresh will be required
			// This will ensure the controls are positioned properly and "flow" when
			// .ctrl is changed
			if .Parent.Base?(FlowControl)
				.WindowRefresh()
			}
		}

	scrollable?(ctrl)
		{
		if false isnt .scrollWndProc(ctrl)
			return true
		// If this is a nested tab and we do not have a scroll bar in the immediate level
		// defer to the higher tab level
		else if ctrl.Controller.Base?(.Base())
			return .scrollable?(ctrl.Controller)
		return false
		}

	scrollWndProc(ctrl)
		{
		return ctrl.Member?(#WndProc) and ctrl.WndProc.Method?(#Adjust)
			? ctrl.WndProc
			: false
		}

	ensureWindowSize(contRect, ctrlRect)
		{
		windW = contRect.right - contRect.left
		windH = contRect.bottom - contRect.top
		relativeW = contRect.right - ctrlRect.left
		relativeH = contRect.bottom - ctrlRect.top
		windW += adjW = relativeW < .Xmin ? .Xmin - relativeW : 0
		windH += adjH = relativeH < .Ymin ? .Ymin - relativeH : 0
		if adjW is 0 and adjH is 0
			return
		SetWindowPos(.WindowHwnd(), 0, 0, 0, windW, windH, SWP.NOMOVE | SWP.NOZORDER)
		.resize(.x, .y, .w + adjW, .h + adjH)
		}

	ConstructTab(i)
		// WARNING: this does NOT fill in RecordControl data (use Select)
		{
		if .ctrls[i] isnt false
			return false
		.construct(i)
		DoStartup(.ctrls[i])
		return true
		}

	ConstructAndSetTab(i)
		{
		if .ConstructTab(i)
			.Send('TabsControl_SelectTab')
		}

	construct(i)
		{
		.being_constructed = i
		.ctrls[i] = .Construct(Object('WndPane',
			Object('Border', .controls[i], .border),
			windowClass: 'SuBtnfaceArrow',
			hidden: i isnt .Tab.GetSelected())
			)
		.being_constructed = false
		return .ctrls[i]
		}

	// NOTE: constructOnEdit tab option should only be used if necessary.
	// Avoid using this option on potentially slow tabs like attachments
	constructTabsOnEdit()
		{
		for (i= 0 ; i < .GetAllTabCount(); i++)
			// don't use GetDefault because sometimes .controls contains classes
			if .controls[i].Member?("constructOnEdit") and
				.controls[i].constructOnEdit is true
				.ConstructAndSetTab(i)
		}

	being_constructed: false
	TabBeingConstructed()
		{
		return .being_constructed isnt false
			? .controls[.being_constructed].Tab
			: false
		}

	TabGetConstructed(hwnd)
		{
		if false is i = .ctrls.FindIf({|x| x.Hwnd is hwnd})
			return false
		return .controls[i].Tab
		}

	TabSplit: ' \xbb '
	TabGetPath()
		{
		if false is tab = .TabBeingConstructed()
			tab = ''

		if 0 is parent_tab = .Send('TabBeingConstructed')
			parent_tab = ''
		else if parent_tab is false
			parent_tab = .Send('TabGetConstructed', .WndProc.Hwnd)

		if not String?(parent_tab)
			parent_tab = ''
		return Opt(parent_tab, .TabSplit) $ tab
		}

	// Use to interact / update tab data prior to control construction
	TabConstructData(i = false)
		{
		return .controls[i is false ? .Tab.GetSelected() : i]
		}

	TabGetSelectedName()
		{
		return .TabName(.GetSelected())
		}

	Constructed?(i)
		{
		return .ctrls[i] isnt false
		}

	TabName(i)
		{
		return .controls[i].Tab
		}

	GetAllTabNames()
		{
		return .controls.Map({ it.Tab })
		}

	FindTab(name)
		{
		return .controls.FindIf({ it.Tab is name })
		}

	FindTabBy(member, value)
		{
		.Tab.ForEachTab({
			| data, idx|
			if data[member] is value
				return idx
			})
		return false
		}

	destroyTab(i)
		{
		if i is .GetSelected()
			.ctrl = false
		if false isnt ctrl = .ctrls.Extract(i, false)
			ctrl.Destroy()
		}

	Resize(x, y, w, h)
		{
		if not .resizing and not .tabs_changing and .being_constructed is false
			.resize(x, y, w, h)
		}

	resizing: false
	resize(.x, .y, .w, .h)
		{
		if .resizing
			return false
		.resizing = true
		.resizeCtrl()
		(.resizeTabsMethod)()
		.Recalc()
		.resizing = false
		}

	resizeCtrl()
		{
		if .ctrl is false
			return
		x = .x + .offsets.x
		y = .y + .offsets.y
		w = Max(.w, .ctrl.Xmin) - .offsets.w
		h = .h - .offsets.h
		.ctrl.Resize(x, y, w, h)
		}

	getter_offsets()
		{ return .offsets = .vertical ? .vertOffsets() : .horzOffsets() }

	getter_resizeTabsMethod()
		{
		return .resizeTabsMethod = Sys.SuneidoJs?()
			? .resizeJsTabs
			: .vertical
				? .resizeVertTabs
				: .resizeHorzTabs
		}

	resizeJsTabs()
		{
		h = .ctrl is false ? 0 : .h
		ctrlHwnd = .ctrl is false ? 0 : .ctrl.Hwnd
		SetWindowPos(.Tab.Hwnd, ctrlHwnd, .x, .y, .w, h, 0)
		}

	resizeVertTabs()
		{
		x = .alternativePos
			? .x + .w - .offsets.w + .offsets.y
			: .x
		.Tab.Resize(x, .y, .Tab.Ymin + .offsets.h, .h)
		}

	resizeHorzTabs()
		{
		y = .alternativePos
			? .y + .offsets.y + .h - .offsets.h
			: .y
		.Tab.Resize(.x, y, .w, .Tab.Ymin)
		}

	tabOffset: 			6
	alternativeOffset: 	3
	vertOffsets()
		{
		yOffset = 4
		offsets = Object()
		offsets.x = .alternativePos
			? .alternativeOffset
			: .Tab.Ymin + 2
		offsets.y = yOffset
		offsets.w = .Tab.Ymin + yOffset
		offsets.h = .tabOffset - yOffset
		return offsets
		}

	horzOffsets()
		{
		offsets = Object()
		offsets.x = 2
		offsets.y = .alternativePos
			? .alternativeOffset
			: .Tab.Ymin + 1
		offsets.w = .tabOffset
		offsets.h = .Tab.Ymin + .alternativeOffset
		return offsets
		}

	SetEnabled(enabled)
		{
		Assert(Boolean?(enabled))
		if .ctrl isnt false
			.ctrl.SetEnabled(enabled)
		.Tab.SetEnabled(enabled)
		}

	readonly: false
	SetReadOnly(readonly)
		{
		if readonly is false
			{
			if .ctrl isnt false // Ctrl should be constructed...
				.constructTabsOnEdit()
			else // wait until constructed
				.Defer(.constructTabsOnEdit)
			}
		.readonly = readonly
		super.SetReadOnly(readonly)
		}

	TabControl_SelChanging(source)
		{
		if source is .Tab // nested tabs
			if .destroyOnSwitch and .stay_on_tab?()
				return true
		return .Send('TabControl_SelChanging', :source)
		}

	stay_on_tab?()
		{
		stay = false
		if (false isnt (items = .Window.GetValidationItems()))
			for (item in items)
				if (not item.ConfirmDestroy())
					if (item.Method?('CloseWindowConfirmation') and
						not item.CloseWindowConfirmation())
					return true
				else
					stay = true
		if (stay and CloseWindowConfirmation(.Window.Hwnd))
			stay = false
		return stay
		}

	// notification from AccessControl
	AccessChanged(@args)
		// pre:		event is an AccessControl event string
		// post:	performs Browse processing for event
		{
		// sometimes observer is removed in the process of Access
		// notifying the observers.  Since Access copies the observers
		// for the notify, this may get done even after removing the access observer
		if (.Empty?())
			return true
		if args[0] is 'before_setdata'
			.access_before_setdata()
		return true
		}

	access_before_setdata()
		{
		// Destroy all BUT the active tab
		inactiveTabs = Seq(.GetTabCount()).Remove(.GetSelected())
		for tab in inactiveTabs
			if .ctrls[tab] isnt false
				{
				.ctrls[tab].Destroy()
				.ctrls[tab] = false
				}
		}

	Getter_Customizable?()
		{
		return .customizable?
		}

	GetControl(i = false)
		// pre: ctrl i must already be constructed
		{
		if false is ctrl = i is false ? .ctrl : .ctrls[i]
			return false

		if ctrl.Destroyed?()
			return false

		return ctrl.GetControl().Ctrl
		}

	TabsControl_GetCurrentControl()
		{
		return .GetControl()
		}

	GetChildren()
		{
		return Object(.Tab).Append(.ctrls.Values().Remove(false))
		}

	GetTabCount()
		{
		return .ctrls.Size()
		}

	GetAllTabCount()
		{
		return .controls.Size()
		}

	SetImageList(images)
		{
		.Tab.SetImageList(images)
		}

	SetImage(i, img)
		{
		.Tab.SetImage(i, img)
		}

	MoveTab(tabIdx, at)
		{
		if tabIdx is at
			return
		if false isnt ctrl = .ctrls.Extract(tabIdx, false)
			.ctrls.Add(ctrl, :at)
		if false isnt tab = .controls.Extract(tabIdx, false)
			.controls.Add(tab, :at)
		.Tab.Move(tabIdx, at)
		}

	CollectFields(customizable = false)
		{
		layouts = .controls.Copy()
		if customizable isnt false
			.CustomizableTabs().Each({ layouts.Add(customizable.Form(it)) })
		return CollectFields(layouts)
		}

	CustomizableTabs()
		{
		return .customizableTabs is false ? Object() : .customizableTabs.Copy()
		}

	Destroy()
		{
		for tab in .. .GetTabCount()
			.destroyTab(tab)
		.Send('RemoveAccessObserver', .AccessChanged)
		super.Destroy()
		}
	}