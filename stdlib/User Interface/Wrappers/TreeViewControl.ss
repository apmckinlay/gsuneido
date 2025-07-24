// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* TreeView
provides:
	AddItem(parent, name, image, container?, param)
	InsertItem(parent, name, image, children, param, insertAfter, state, stateMask)
	GetName(item) => name
	GetParam(item) => number
	EditLabel(item)
	EnsureVisible(item)
	ExpandItem(item, collapse = false)
	Expanded?(item)
	DeleteItem(item)
	Children?(item)
	HasChildren?(item)
	ItemExists?(item) => bool
	GetChildren(item) => list
	GetParent(item)
	GetImage(item)
	GetItemState(item, state, stateMask)
	SetName(item, name)
	SetImageList(images)
	SetImage(item, image)
	SetItemState(item, state, stateMask)
	SelectItem(item)
	SortChildren(hParent, lpfnCompare)
sends:
	SelectTreeItem(olditem, newitem)
	Expanding(item)
	Collapsed(item)
	Rename(item, newname)
	Move(dragging, target)
	Container?
	Children?
*/
WndProc
	{
	Name: 		'TreeView'
	Xmin:		50
	Ymin:		100
	Xstretch:	1
	Ystretch:	1
	editwndproc_old: 0	// Address of old wndproc of edit control for label editing
	scrolldir: 	false	// Scroll direction for scroll timer
	hScrollTimer: 0		// Handle to scroll timer
	dragimage:	0		// Handle to drag image list

	New(.readonly = false, style = 0)
		{
		style |= WS.CHILD | WS.VISIBLE | TVS.HASLINES |
			TVS.LINESATROOT | TVS.HASBUTTONS | TVS.SHOWSELALWAYS
		if .readonly
			style |= TVS.DISABLEDRAGDROP
		else
			style |= TVS.EDITLABELS
		.createControls(style)
		.dragging = false
		.target = false
		.rightbuttondrag = true
		// Message map
		.Map = Object()
		.Map[TVN.SELCHANGED]	= "TVN_SELCHANGED"
		.Map[TVN.SELCHANGING]	= "TVN_SELCHANGING"
		.Map[TVN.GETDISPINFO]	= "TVN_GETDISPINFO"
		.Map[TVN.ITEMEXPANDING] = "TVN_ITEMEXPANDING"
		.Map[TVN.ITEMEXPANDED]	= "TVN_ITEMEXPANDED"
		.Map[TVN.BEGINDRAG]		= "TVN_BEGINDRAG"
		.Map[TVN.BEGINRDRAG]	= "TVN_BEGINRDRAG"
		.Map[TVN.BEGINLABELEDIT]	= "TVN_BEGINLABELEDIT"
		.Map[TVN.ENDLABELEDIT]	= "TVN_ENDLABELEDIT"
		.Map[TVN.KEYDOWN]		= "TVN_KEYDOWN"
		.Map[NM.CLICK] 			= "NM_CLICK"
		.Map[NM.KILLFOCUS] 		= "NM_KILLFOCUS"
		.Map[TTN.SHOW]			= "TTN_SHOW"
		}

	toolTips: false
	createControls(style)
		{
		.CreateWindow(WC_TREEVIEW, '', style)
		.SubClass()
		.destroyToolTips()
		.toolTips = .Construct(#(ToolTip))
		.toolTips.AddTool(.Hwnd, LPSTR_TEXTCALLBACK)
		// When style: TVS.NOTOOLTIPS is not used, TreeView creates its own tooltips.
		// In order to properly handle the tooltips, we replace the default handler with
		// our own ToolTipControl. As TVM.SETTOOLTIPS returns the handle to the
		// default handler, we can destroy it to avoid any potential conflicts.
		if NULL isnt oldToolTips = .SendMessage(TVM.SETTOOLTIPS, .toolTips.Hwnd)
			DestroyWindow(oldToolTips)
		}

	destroyToolTips()
		{
		if .toolTips isnt false
			.toolTips.Destroy()
		.toolTips = false
		}

	TTN_SHOW(lParam /*unused*/)
		{
		return true // Prevents the TreeView tooltips from adjusting the window z-order
		}

	TVN_SELCHANGING(lParam)
		{
		nmtv = NMTREEVIEW(lParam)
		// direct call via code: 0, mouse: 1, keyboard: 2
		return nmtv.action in (0, 1, 2)
			? nmtv.itemOld.hItem is nmtv.itemNew.hItem
			: true
		}

	TVN_SELCHANGED(lParam)
		{
		nmtv = NMTREEVIEW(lParam)
		olditem = nmtv.itemOld.hItem
		newitem = nmtv.itemNew.hItem
		if newitem isnt 0
			.Tree_SelChanged(olditem, newitem)
		return 0
		}

	Tree_SelChanged(olditem, newitem)
		{
		.Send("SelectTreeItem", olditem, newitem)
		}

	TVN_KEYDOWN(lParam)
		{
		ctrlPressed? = NMTVKEYDOWN(lParam).wVKey is VK.CONTROL
		hSelectedItem = .GetLastSelection(ctrlPressed?)
		if NMTVKEYDOWN(lParam).wVKey is VK.F2
			{
			if hSelectedItem isnt 0 and not .Static?(hSelectedItem)
				{
				.SendMessage(TVM.ENSUREVISIBLE, 0, hSelectedItem)
				.EditLabel(hSelectedItem)
				}
			}
		else if not ctrlPressed?
			.Send(#TreeView_KeyboardNavigation, hSelectedItem)
		return 'callsuper'
		}
	// overriden in MultiTreeViewControl
	GetLastSelection(ctrlPressed? /*unused*/)
		{
		return .GetSelectedItem()
		}

	// overriden in MultiTreeViewControl
	NM_KILLFOCUS(lParam /*unused*/)
		{
		return 0
		}

	NM_CLICK()
		{
		GetCursorPos(pt = Object())
		ScreenToClient(.Hwnd, pt)
		if 0 is item = SendMessageTreeHitTest(.Hwnd, TVM.HITTEST, 0, info = Object(:pt))
			return 0
		if info.flags is TVHT.ONITEMLABEL
			.Send('TreeView_ItemClicked', item)
		else
			.Send('TreeView_Clicked', item, info.flags)
		return 0
		}

	TVN_GETDISPINFO(lParam)
		{
		// Set the cChildren member for container items that have children...
		// Note: it must be set to some value, because Windows does not initialize it
		di = NMTVDISPINFO2(lParam)
		if ((di.item.mask & TVIF.CHILDREN) is TVIF.CHILDREN)
			StructModify(NMTVDISPINFO2, lParam,
				{ it.item.cChildren = .Children?(it.item.hItem) ? 1 : 0 })
		return 0
		}

	TVN_ITEMEXPANDING(lParam)
		{
		tv = NMTREEVIEW(lParam)
		action = tv.action is TVE.COLLAPSE ? #Collapsing : #Expanding
		return .AttemptSend(action, tv.itemNew.hItem)
		}

	AttemptSend(@args)
		{
		if .Send(@args) is 0 and .Method?(args[0])
			(this[args[0]])(@+1args)
		return 0
		}

	TVN_ITEMEXPANDED(lParam)
		{
		tv = NMTREEVIEW(lParam)
		action = tv.action is TVE.COLLAPSE ? #Collapsed : #Expanded
		return .AttemptSend(action, tv.itemNew.hItem)
		}

	TVN_BEGINLABELEDIT(lParam)
		{
		.setupEdit()
		.editwndproc_old = GetWindowLongPtr(.hedit, GWL.WNDPROC)
		if .editwndproc_old isnt 0
			SetWindowProc(.hedit, GWL.WNDPROC, .EditWndProc)
		// Return true to immediately cancel label editing for non-editable items
		if (noedit = .Static?(NMTVDISPINFO2(lParam).item.hItem))
			Beep()	// Beep if non-editable
		return noedit
		}

	TVN_ENDLABELEDIT(lParam)
		{
		.clearEdit()
		tvdi = NMTVDISPINFO(lParam)
		if tvdi.item.pszText is false
			return 0
		item = tvdi.item.hItem
		name = tvdi.item.pszText
		return .Send("Rename", item, name)
		}

	hedit: false
	oldAccels: false

	setupEdit()
		{
		.hedit = .SendMessage(TVM.GETEDITCONTROL, 0, 0)
		.AddHwnd(.hedit)
		}

	clearEdit()
		{
		if .hedit is false
			return
		.DelHwnd(.hedit)
		.hedit = false
		.oldAccels = false
		}

	On_Copy()
		{
		.sendEditMessage(WM.COPY, altMsg: "On_Context_Copy")
		}

	On_Paste()
		{
		.sendEditMessage(WM.PASTE, altMsg: "On_Context_Paste")
		}

	On_Cut()
		{
		.sendEditMessage(WM.CUT, altMsg: "On_Context_Cut")
		}

	On_Delete()
		{
		.sendEditMessage(WM.CLEAR)
		}

	On_Undo()
		{
		.sendEditMessage(WM.UNDO)
		}

	On_Select_All()
		{
		.sendEditMessage(EM.SETSEL, 0, -1)
		}

	sendEditMessage(msg, wParam = 0, lParam = 0, altMsg = false)
		{
		if .hedit isnt false
			SendMessage(.hedit, msg, wParam, lParam)
		else if altMsg isnt false
			.Send(altMsg)
		}

	// drag and drop support
	TVN_BEGINDRAG(lParam)
		{
		tv = NMTREEVIEW(lParam)
		item = tv.itemNew.hItem
		p = tv.ptDrag
		SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.CARET, item)
		.destroyDragImage()	// Destroy the drag image list
		.dragimage = SendMessage(.Hwnd, TVM.CREATEDRAGIMAGE, 0, item)
		ImageList_BeginDrag(.dragimage, 0, 0, 0)
		ImageList_DragEnter(.Hwnd, p.x, p.y)
		SetCapture(.Hwnd)
		.dragging = item
		.target = false
		return 0
		}
	TVN_BEGINRDRAG(lParam)
		{
		// The style TVS.DISABLEDRAGDROP only suppresses TVN.BEGINDRAG (left mouse button)
		// As a result, we have to manually suppress the right mouse button drag event
		if .readonly
			return 0
		.rightbuttondrag = true
		.TVN_BEGINDRAG(lParam)
		}
	MOUSEMOVE(lParam)
		{
		if .dragging is false
			return "callsuper"
		x = LOWORD(lParam)
		y = HIWORD(lParam)
		ImageList_DragMove(x, y)
		// hit test
		tvht = Object(pt: Object(:x, :y))
		item = SendMessageTreeHitTest(.Hwnd, TVM.HITTEST, 0, tvht)
		if not .Container?(item)
			item = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.PARENT, item)
		if item is NULL or item is .dragging or not .Container?(item)
			{
			ImageList_DragLeave(.Hwnd)
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
			ImageList_DragEnter(.Hwnd, x, y)
			.target = false
			}
		else
			{
			ImageList_DragLeave(.Hwnd)
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, item)
			ImageList_DragEnter(.Hwnd, x, y)
			.target = item
			}
		// scroll timer
		.SetScrollTimer(item, x, y)
		return "callsuper"
		}
	LBUTTONUP(wParam)
		{
		if .dragging is false
			return 0
		.destroyScrollTimer()
		ImageList_DragLeave(.Hwnd)
		ImageList_EndDrag()
		ReleaseCapture()
		if .target isnt false
			{
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
			if not .Static?(.dragging)
				{
				if ((wParam & MK.CONTROL) is MK.CONTROL)	// Copy if Ctrl is pressed
					.AttemptSend(#DragCopy, .dragging, .target)
				else
					.AttemptSend(#DragMove, .dragging, .target)	// Otherwise, Move
				}
			}
		.dragging = .target = false
		return 0
		}
	RBUTTONUP(wParam, lParam)
		{
		// NOTE: sequence is important in this function...
		// End right dragging...
		.rightbuttondrag = false
		if .dragging is false
			return 0
		.destroyScrollTimer()
		ImageList_DragLeave(.Hwnd)
		ImageList_EndDrag()
		ReleaseCapture()
		if .target isnt false
			{
			if .Static?(.dragging)
				{
				SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
				.dragging = .target = false
				return 0
				}
			menu = Object(
					 Object(name: '&Move Here' id: 1),
					 Object(name: '&Copy Here' id: 2),
					 '', 'Cancel'
					 )
			// If Ctrl is pressed, 'Copy' is ithe default item; otherwise it is 'Move'
			if ((wParam & MK.CONTROL) is MK.CONTROL)
				menu[1].def = true
			else
				menu[0].def = true
			// Create a popup menu with 'Move Here, Copy Here and Cancel commands'
			ClientToScreen(.Hwnd, pt = Object(x: LOWORD(lParam) y: HIWORD(lParam)))
			x = ContextMenu(menu).Show(.Hwnd, pt.x, pt.y)
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
			switch (x)
				{
			case 1:	// Move
				.AttemptSend(#DragMove, .dragging, .target)
			case 2:	// Copy
				.AttemptSend(#DragCopy, .dragging, .target)
			default :
				}
			}
		.dragging = .target = false
		return 0
		}
	// end of drag & drop support

	AddItem(parent, name, image = 0, container? = false, param = 0)
		{
		tvi = Object(
			mask: TVIF.TEXT | TVIF.IMAGE | TVIF.SELECTEDIMAGE | TVIF.CHILDREN |
				TVIF.PARAM,
			pszText: name,
			iImage: image,
			iSelectedImage: image,
			cChildren: container? ? I.CHILDRENCALLBACK : 0,
			lParam: param
			)
		tvins = Object(
			item: tvi,
			hInsertAfter: TVI.LAST,
			hParent: parent
			)
		return SendMessageTreeInsert(.Hwnd, TVM.INSERTITEM, 0, tvins)
		}
	InsertItem(parent, name, image = 0, children = 0, param = 0,
		insertAfter = false, state = 0, stateMask = 0)
		{
		tvi = Object(
			mask: TVIF.TEXT | TVIF.IMAGE | TVIF.SELECTEDIMAGE | TVIF.CHILDREN |
				TVIF.PARAM | TVIF.STATE,
			:state,
			:stateMask,
			pszText: name,
			iImage: image,
			iSelectedImage: image,
			cChildren: children,
			lParam: param
			)
		tvins = Object(
			item: tvi,
			hInsertAfter: insertAfter is false ? TVI.LAST : insertAfter,
			hParent: parent
			)
		return SendMessageTreeInsert(.Hwnd, TVM.INSERTITEM, 0, tvins)
		}
	maxNameSize: 128
	GetName(item)
		{
		tvi = Object(
			hItem: item,
			mask: TVIF.TEXT,
			cchTextMax: .maxNameSize)
		SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi)
		return tvi.pszText
		}
	Highlight(parent, child)
		{
		parentItem = .findItemBy(TVI.ROOT)
			{ |item| parent is .GetName(item).AfterFirst(':') ? item : false }
		if parentItem is false
			return
		childItem = .findItemBy(parentItem)
			{ |item| child is .GetName(item).BeforeFirst(':') ? item : false }
		if childItem isnt false
			.SelectItem(childItem)
		}
	findItemBy(item, block)
		{
		which = TVGN.CHILD // Get parents first child
		while 0 isnt item = .SendMessage(TVM.GETNEXTITEM, which, item)
			{
			if false isnt block(item)
				return item
			which = TVGN.NEXT // Cycle through parents children
			}
		return false
		}
	GetParam(item)
		{
		tvi = Object(hItem: item, mask: TVIF.PARAM)
		return (SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi) is 0) ?
			false : Number(tvi.lParam)
		}
	EditLabel(item)
		{
		SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.CARET, item)
		SendMessage(.Hwnd, TVM.EDITLABEL, 0, item)
		}
	EnsureVisible(item)
		{ .SendMessage(TVM.ENSUREVISIBLE, 0, item) }
	ExpandItem(item, collapse = false)
		{
		if collapse isnt true
			{
			// reset EXPANDEDONCE flag or else expand won't work
			tvi = Object(
				hItem: item,
				mask: TVIF.STATE,
				stateMask: TVIS.EXPANDEDONCE,
				state: 0
				)
			SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
			// Expand the item
			.SendMessage(TVM.EXPAND, TVE.EXPAND, item)
			}
		else
			{
			// Collapse the item
			.SendMessage(TVM.EXPAND, TVE.COLLAPSE, item)
			.AttemptSend(#Collapsed, item)
			}
		}
	Expanded?(item)
		{
		// Return true if the item is expanded
		tvi = Object(
			hItem: item,
			mask: TVIF.STATE,
			stateMask: TVIS.EXPANDED,
			)
		SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi)
		return (tvi.state & TVIS.EXPANDED) is TVIS.EXPANDED
		}
	DeleteItem(item)
		{
		return (SendMessage(.Hwnd, TVM.DELETEITEM, 0, item) isnt 0)
		}
	DeleteAllItems()
		{
		return (SendMessage(.Hwnd, TVM.DELETEITEM, 0, TVI.ROOT) isnt 0)
		}
	Children?(item)
		{
		return Boolean?(x = .Send("Children?", item)) ? x : false
		}
	Container?(item)
		{
		return Boolean?(x = .Send("Container?", item)) ? x : false
		}
	Static?(item)
		{
		return Boolean?(x = .Send("Static?", item)) ? x : false
		}
	HasChildren?(item)
		{
		child = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.CHILD, item)
		return child isnt NULL
		}
	ItemExists?(item)
		{ return 1 is SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, Object(hItem: item)) }
	GetChildren(item = false)
		{
		if item is false // standard control GetChildren
			return #() // no child controls
		children = Object()
		for (item = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.CHILD, item);
			item isnt 0;
			item = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.NEXT, item))
			{
			children.Add(item)
			}
		return children
		}
	GetParent(item)
		{
		return SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.PARENT, item)
		}
	GetImage(item)
		{
		tvi = Object(mask: TVIF.IMAGE, hItem: item)
		SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi)
		return tvi.iImage
		}
	GetItemState(item, stateMask = 0xffffffff)
		{
		tvi = Object(mask: TVIF.STATE, hItem: item, :stateMask)
		SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi)
		return tvi.state
		}
	GetSelectedItem()
		{
		return .SendMessage(TVM.GETNEXTITEM, TVGN.CARET, 0)
		}
	SetName(item, name)
		{
		tvi = Object(
			hItem: item,
			mask: TVIF.TEXT,
			pszText: name)
		SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
		}
	SetImageList(images)
		{
		SendMessage(.Hwnd, TVM.SETIMAGELIST, TVSIL.NORMAL, images)
		}
	SetImage(item, image)
		{
		tvi = Object(
			hItem: item,
			mask: TVIF.IMAGE + TVIF.SELECTEDIMAGE,
			iImage: image,
			iSelectedImage: image)
		SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
		}
	SetItemState(item, state, stateMask = 0xffffffff)
		{
		tvi = Object(
				  mask: TVIF.STATE,
				  hItem: item,
				  :state,
				  :stateMask
			  )
		SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
		}
	SelectItem(item)
		{ SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.CARET, item) }
	ScrollItemIntoView(item)
		{ SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.FIRSTVISIBLE, item) }
	ContextMenu(x, y)
		{
		pt = Object(:x, :y)
		if pt.x is 0 and pt.y is 0
			// If both coordinates are 0, the message stemmed from a keyboard event.
			// This means the context menu should be displayed at the current selection
			{
			hSelectedItem = .GetSelectedItem()
			if hSelectedItem isnt 0
				.SendMessage(TVM.ENSUREVISIBLE, 0, hSelectedItem)
			SendMessageRect(.Hwnd, TVM.GETITEMRECT, true,
				rcItem = Object(left: hSelectedItem))
			pt.x = rcItem.left
			pt.y = rcItem.top + ((rcItem.bottom - rcItem.top) / 2).Int()
			ClientToScreen(.Hwnd, pt)
			return .ShowContextMenu(hSelectedItem, pt.x, pt.y)
			}
		ScreenToClient(.Hwnd, pt)
		tvht = Object(:pt)
		if 0 isnt item = SendMessageTreeHitTest(.Hwnd, TVM.HITTEST, 0, tvht)
			.SelectItem(item)
		return .ShowContextMenu(item, x, y)
		}

	// Send "ShowContextMenu" event to controller first argument is handle to item
	// that was clicked or NULL if no item was clicked.
	// The other arguments are screen coordinates of point clicked
	ShowContextMenu(item, x, y)
		{ return .Send(#ShowContextMenu, item, x, y) }

	EditWndProc(hwnd, msg, wParam, lParam)
		{
		_hwnd = .WindowHwnd()
		// Handle messages for the edit control that edits item labels.
		// This is so that return and escape will stop label editing
		switch (msg)
			{
		case WM.KEYDOWN:
			if wParam is VK.RETURN
				.SendMessage(TVM.ENDEDITLABELNOW, 0, 0)
		case WM.GETDLGCODE:
			return DLGC.WANTALLKEYS
		case WM.DESTROY:
			SetWindowProc(hwnd, GWL.WNDPROC, .editwndproc_old)
			ClearCallback(.EditWndProc)
		default:
			}
		return CallWindowProc(.editwndproc_old, hwnd, msg, wParam, lParam)
		}
	SetScrollTimer(item, x, y)
		{
		if item isnt NULL
			.handleItemScrollTimer(x, y)
		else
			.destroyScrollTimer()
		}
	handleItemScrollTimer(x, y)
		{
		// Check if capture scroll timer should be enabled
		newdir = false
		GetClientRect(.Hwnd, rcClient = Object())
		item_size = SendMessage(.Hwnd, TVM.GETITEMHEIGHT, 0, 0)
		if (y < SendMessage(.Hwnd, TVM.GETITEMHEIGHT, 0, 0))
			newdir = 'up'
		else if (y > rcClient.bottom - item_size)
			newdir = 'down'
		else if (x < item_size)
			newdir = 'left'
		else if (x > rcClient.right - item_size)
			newdir = 'right'
		if newdir isnt .scrolldir
			{
			.destroyScrollTimer()
			.scrolldir = newdir
			}
		if .scrolldir isnt false and .hScrollTimer is 0
			{
			.dragx = x
			.dragy = y
			.hScrollTimer = SetTimer(NULL, 0, 300 /* =ms */, .ScrollTimer)
			.ScrollTimer()
			}
		}
	ScrollTimer(@unused)
		{
		// Scrolls TreeView (in response to WM_TIMER) when mouse is captured
		ImageList_DragLeave(.Hwnd)
		switch (.scrolldir)
			{
		case 'up':		SendMessage(.Hwnd, WM.VSCROLL, SB.LINEUP, 0)
		case 'down':	SendMessage(.Hwnd, WM.VSCROLL, SB.LINEDOWN, 0)
		case 'left':	SendMessage(.Hwnd, WM.HSCROLL, SB.LINELEFT, 0)
		case 'right':	SendMessage(.Hwnd, WM.HSCROLL, SB.LINERIGHT, 0)
		case false :	// ignore
			}
		ImageList_DragEnter(.Hwnd, .dragx, .dragy)
		.MOUSEMOVE(.dragx | .dragy << 16)
		return
		}
	destroyScrollTimer()
		{
		if .hScrollTimer isnt NULL
			{
			KillTimer(NULL, .hScrollTimer)
			ClearCallback(.ScrollTimer)
			.hScrollTimer = NULL
			}
		}
	destroyDragImage()
		{
		// Ensure that the drag image list is destroyed
		if .dragimage isnt NULL
			{
			ImageList_Destroy(.dragimage)
			.dragimage = NULL
			}
		}

	// lpfnCompare(lParam1, lParam2, lParamSort) returns:
	// 	-# if lParam1 comes first
	//	+# if lParam2 comes first
	// 	NOTE: lParamSort, corresponds to the optional lParam member in TVSORTCB.
	//		  Currently unused
	SortChildren(hParent, lpfnCompare)
		{
		tvitem = Object(:hParent, :lpfnCompare)
		SendMessageTreeSort(.Hwnd, TVM.SORTCHILDRENCB, 0, tvitem)
		ClearCallback(lpfnCompare)
		}

	ForEachChild(parent, block)
		{
		for (x = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.CHILD, parent);
			x isnt 0;
			x = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.NEXT, x))
			{
			block(x)
			}
		}

	Reset()
		{
		.destroyScrollTimer()
		.destroyDragImage()
		style = GetWindowLong(.Hwnd, GWL.STYLE)
		DestroyWindow(.Hwnd)
		.createControls(style)
		.ResetTheme()
		super.Resize(.x, .y, .w, .h)
		}
	Resize(.x, .y, .w, .h)
		{
		super.Resize(x, y, w, h)
		}
	ResetTheme()
		{
		theme = IDE_ColorScheme.GetTheme()
		SendMessage(.Hwnd, TVM.SETBKCOLOR, 	 0, theme.defaultBack)
		SendMessage(.Hwnd, TVM.SETTEXTCOLOR, 0, theme.defaultFore)
		SendMessage(.Hwnd, TVM.SETLINECOLOR, 0, theme.defaultFore)
		}
	Refresh()
		{
		.destroyScrollTimer()
		.destroyDragImage()
		.DeleteAllItems()
		.ResetTheme()
		}
	DESTROY()
		{
		.destroyScrollTimer()
		.destroyDragImage()
		.destroyToolTips()
		return "callsuper"
		}
	}
