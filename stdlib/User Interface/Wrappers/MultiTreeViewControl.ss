// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TreeView control with custom, multiple selection capabilities
// REFACTOR: should not use private methods/members from TreeViewControl
TreeViewControl
	{
	New(.multi? = true, readonly = false, style = 0)
		{
		super(readonly, style)
		.selection = Object()
		.ResetTheme()
		}

	Tree_SelChanged(olditem, newitem)
		{
		super.Tree_SelChanged(olditem, newitem)
		.setselected(newitem, true)
		}

	// TVN_BEGINDRAG notification handler, overrides TreeViewControl
	TVN_BEGINDRAG(lParam)
		{
		tv = NMTREEVIEW(lParam)
		.settips(false) // Disable tooltips so painting isn't screwed up
		.TreeViewControl_destroyDragImage()	// Destroy the drag image list
		.TreeViewControl_dragimage = .createdragimage()
		ImageList_BeginDrag(.TreeViewControl_dragimage, 0, 0, 0)
		ImageList_DragEnter(.Hwnd, tv.ptDrag.x, tv.ptDrag.y)
		SetCapture(.Hwnd)
		.TreeViewControl_dragging = tv.itemNew.hItem
		.TreeViewControl_target = false
		return 0
		}
	TVN_BEGINLABELEDIT(lParam)
		{
		if .selection.Size() > 1
			{
			.UnselectAll(NMTVDISPINFO2(lParam).item.hItem)
			return true
			}
		else
			return super.TVN_BEGINLABELEDIT(lParam)
		}
	LBUTTONDOWN(wParam, lParam)	// allow multiple selection
		{
		if not .multi?
			return "callsuper"
		if .ctrlPressed?(wParam)
			{
			SetFocus(.Hwnd)
			return 0
			}
		x = .hittest(LOWORD(lParam), HIWORD(lParam))
		if x isnt NULL and .selection.Find(x) is false
			{
			.UnselectAll()
			// Because the selection notification comes on LBUTTONUP,
			// but item must be selected to be dragged...
			.setselected(x, true)
			// Next two lines are so that the selection is correct when dragging.
			// Without this, sometimes the wrong item gets saved when dragging.
			// SetFocus(0) is to prevent the treeview from going into label edit mode
			SetFocus(0)
			.SelectItem(x)
			}
		return "callsuper"	// Single-select,  Let TreeView handle it
		}
	ctrlPressed?(wParam)
		{
		return MK.CONTROL is (wParam & MK.CONTROL)
		}

	NM_KILLFOCUS(lParam /*unused*/)
		{
		.UnselectAll(.GetLastSelection(false))
		return 0
		}

	RBUTTONDOWN(lParam)
		{
		if .multi? and
			NULL isnt x = .hittest(LOWORD(lParam), HIWORD(lParam))
			{
			if .selection.Find(x) isnt false
				{
				GetCursorPos(pt = Object())
				.ShowContextMenu(.GetLastSelection(true), pt.x, pt.y)
				return -1
				}
			else // clicked on an item that is not selected
				{
				.UnselectAll()
				// Because the selection notification comes on RBUTTONUP,
				// but item must be selected to be dragged...
				.setselected(x, true)
				}
			}
		return "callsuper"	// Let TreeView handle other responsibilities
		}

	multiSelect?()
		{
		return .multi? and KeyPressed?(VK.CONTROL)
		}

	LBUTTONUP(wParam, lParam)
		{
		if not .multi?
			return super.LBUTTONUP(:wParam, :lParam)
		if .TreeViewControl_dragging is false
			{
			if NULL is hItem = .hittest(LOWORD(lParam), HIWORD(lParam))
				return 0
			.dragItem(hItem, wParam)
			return 0
			}
		.settips(true)					// Re-enabled tooltips disabled for drag operation
		.TreeViewControl_destroyScrollTimer()
		ImageList_DragLeave(.Hwnd)
		ImageList_EndDrag()
		ReleaseCapture()
		if .TreeViewControl_target isnt false
			.dropItem(wParam)
		.TreeViewControl_dragging = .TreeViewControl_target = false
		return 0
		}
	dragItem(hItem, wParam)
		{
		index = .selection.Find(hItem)
		if .ctrlPressed?(wParam)
			{
			tvi = Object(:hItem, mask: TVIF.STATE, stateMask: TVIS.SELECTED)
			SendMessageTreeItem(.Hwnd, TVM.GETITEM, 0, tvi)
			.setTviState(tvi, hItem)
			if index is 0 and .selection.Size() > 1 and tvi.state is 0
				{
				x = .selection.Copy()
				.selection = Object()
				SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.CARET, x[1])
				.selection = x
				}
			SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
			.setselected(hItem, tvi.state is TVIS.SELECTED)
			.Send("SelectTreeItem", false, false)
			}
		else if index isnt false
			{
			if hItem is .selection[0]
				.UnselectAll(hItem)
			else
				{
				.UnselectAll()
				SetFocus(0)
				.SelectItem(hItem)
				SetFocus(.Hwnd)
				}
			}
		}

	setTviState(tvi, hItem)
		{
		tvi.state = (tvi.state & TVIS.SELECTED) is TVIS.SELECTED
			? 0	: .CanMultiSelect?(hItem)
				? TVIS.SELECTED	: 0
		}

	dropItem(wParam)
		{
		SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
		for dragindex in (selection = .selection.Copy()).Members()
			if not .Static?(selection[dragindex])
				// Copy if Ctrl is pressed, otherwise, Move
				.AttemptSend(
					.ctrlPressed?(wParam) ? "DragCopy" : "DragMove",
					selection[dragindex],
					.TreeViewControl_target,
					last?: dragindex is (selection.Size() - 1))
		}
	RBUTTONUP(wParam, lParam)
		{
		if not .multi?
			return super.RBUTTONUP(:wParam, :lParam)
		// NOTE: Sequence is important in this function...
		// End right dragging...
		.TreeViewControl_rightbuttondrag = false
		if .TreeViewControl_dragging is false
			return 0
		.TreeViewControl_destroyScrollTimer()
		.settips(true) // Re-enabled tooltips disabled for drag operation
		ImageList_DragLeave(.Hwnd)
		ImageList_EndDrag()
		ReleaseCapture()
		if .TreeViewControl_target isnt false
			{
			if .Static?(.selection[0])
				{
				SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
				.TreeViewControl_dragging = .TreeViewControl_target = false
				return 0
				}
			menu = Object(
				Object(name: '&Move Here' id: 1),
				Object(name: '&Copy Here' id: 2),
				'', 'Cancel')
			if MK.CONTROL is (wParam & MK.CONTROL)
				menu[1].def = true	// If Ctrl is pressed, 'Copy' is ithe default item
			else
				menu[0].def = true	// otherwise it is 'Move'...
			ClientToScreen(.Hwnd, pt = Object(x: LOWORD(lParam) y: HIWORD(lParam)))
			// Display a popup menu with 'Move Here, Copy Here and Cancel commands'
			x = ContextMenu(menu).Show(.Hwnd, pt.x, pt.y)
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
			.dragAction(x)
			}
		.TreeViewControl_dragging = .TreeViewControl_target = false
		return 0
		}

	dragAction(x)
		{
		if x not in (1, 2)
			return
		cmd =  x is 1 ? "DragMove" : "DragCopy"
		for dragindex in (selection = .selection.Copy()).Members()
			.AttemptSend(cmd, selection[dragindex], .TreeViewControl_target,
				last?: dragindex is (selection.Size() - 1))
		}

	MOUSEMOVE(lParam)
		{
		if not .multi?
			return super.MOUSEMOVE(:lParam)
		if .TreeViewControl_dragging is false
			return "callsuper"
		ptMouse = Object(x: LOWORD(lParam) y: HIWORD(lParam))
		ImageList_DragMove(ptMouse.x, ptMouse.y)
		// hit test
		item = SendMessageTreeHitTest(.Hwnd, TVM.HITTEST, 0, Object(pt: ptMouse.Copy()))
		if not .Container?(item)
			item = SendMessage(.Hwnd, TVM.GETNEXTITEM, TVGN.PARENT, item)
		ImageList_DragLeave(.Hwnd)
		if .selection.Find(item) isnt false or item is NULL or not .Container?(item)
			{
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, NULL)
			.TreeViewControl_target = false
			}
		else
			{
			SendMessage(.Hwnd, TVM.SELECTITEM, TVGN.DROPHILITE, item)
			.TreeViewControl_target = item
			}
		ImageList_DragEnter(.Hwnd, ptMouse.x, ptMouse.y)
		// scroll timer
		.SetScrollTimer(item, ptMouse.x, ptMouse.y)
		// default processing
		return "callsuper"
		}
	UnselectAll(hSkipItem = false)
		{
		.selection.Copy().Remove(hSkipItem).Each(.UnselectItem)
		}
	Getter_Selection()
		{
		return .selection
		}
	GetLastSelection(ctrlPressed?)
		{
		if .selection.Empty?()
			return 0
		if ctrlPressed? and .multi?
			return .selection.Last()
		if .selection.Size() > 1
			.UnselectAll(.selection.Last())
		return .selection[0]
		}
	hittest(x, y)
		{
		// Return item handle if hittest succeeds and occurs on an item label or bitmap,
		// or NULL otherwise
		tvht = Object(pt: Object(:x, :y))
		item = SendMessageTreeHitTest(.Hwnd, TVM.HITTEST, 0, tvht)
		return ((tvht.flags & TVHT.ONITEM) isnt 0) ? item : NULL
		}
	setselected(item, select)
		{
		// Insert or remove the item handle 'item' into/from the array of selected items
		if select is true
			.selection.AddUnique(item)
		else
			.selection.Remove(item)
		}
	DeleteItem(item)
		{
		if not super.DeleteItem(item)
			return false
		.setselected(item, false)
		return true
		}
	SelectItem(item)
		{
		super.SelectItem(item)
		if .selection.Size() > 1 and not .multiSelect?()
			.UnselectAll(item)
		}
	UnselectItem(item)
		{
		tvi = Object(hItem: item mask: TVIF.STATE stateMask: TVIS.SELECTED state: 0)
		SendMessageTreeItem(.Hwnd, TVM.SETITEM, 0, tvi)
		.setselected(item, false)
		}
	getchildren(item, item_hierarchy)
		{
		// Add item's children to iteme hierarchy (for CanMultiSelect?)
		for (item = .SendMessage(TVM.GETNEXTITEM, TVGN.CHILD, item);
			item isnt 0;
			item = .SendMessage(TVM.GETNEXTITEM, TVGN.NEXT, item))
			{
			item_hierarchy.Add(item)
			.getchildren(item, item_hierarchy)
			}
		}
	CanMultiSelect?(item)
		{
		if .Static?(item)
			return false
		// Get item hierarchy
		item_hierarchy = Object()
		.getchildren(item, item_hierarchy)
		for (item = .SendMessage(TVM.GETNEXTITEM, TVGN.PARENT, item);
			item isnt 0;
			item = .SendMessage(TVM.GETNEXTITEM, TVGN.PARENT, item))
			item_hierarchy.Add(item)
		// Check for statics already in selection and redundancy
		for selected in .selection
			{
			if .Static?(selected)
				return false
			// Ensure that no redundant selection occurs
			if item_hierarchy.Find(selected) isnt false
				return false
			}
		return true
		}
	createdragimage()
		{
		if 1 > selcount = .selection.Size()	// Return nothing if nothing selected
			return 0
		if selcount is 1					// Return normal drag image for one item
			return SendMessage(.Hwnd, TVM.CREATEDRAGIMAGE, 0, .selection[0])
		// Generate image lists for each individual item, up to 4
		images = Object(info: Object())
		imageListLimit = 4
		if selcount > imageListLimit
			{
			x = .AddItem(0, "... [" $ String(selcount) $ " items total ] ...", -1)
			images.Add(SendMessage(.Hwnd, TVM.CREATEDRAGIMAGE, 0, x), at: 0)
			super.DeleteItem(x)
			}
		for (i = 0; i < imageListLimit and i < selcount; i++)
			images.Add(SendMessage(.Hwnd, TVM.CREATEDRAGIMAGE, 0, .selection[i]), at: i)
		// Get image data
		for (i = 0; i < images.Size(list:); i++)
			{
			ImageList_GetImageInfo(images[i], 0, imginf = Object())
			images.info.Add(imginf, at: i)
			}
		for (dragimage = images[0], i = 1, y = images.info[0].rcImage.bottom;
			i < images.Size(list:); i++)
			{
			oldimage = dragimage
			dragimage = ImageList_Merge(oldimage, 0, images[i], 0, 0, y)
			y += images.info[i].rcImage.bottom
			ImageList_Destroy(oldimage)
			ImageList_Destroy(images[i])
			}
		return dragimage
		}
	settips(show = true)
		{
		// Enable or disable tooltips (for drag and drop purposes)
		style = GetWindowLong(.Hwnd, GWL.STYLE)
		style = (show) ? style & ~TVS.NOTOOLTIPS : style | TVS.NOTOOLTIPS
		SetWindowLong(.Hwnd, GWL.STYLE, style)
		}
	}
