// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/*
parent should call:
	Inactivate()
	Destroy()
model must support:
	Nextnum() -> integer
	NewItem
	Update
	DeleteItem
	Children
	Children?
	Container?
	Static?
	EnsureUnique
	Get
model can support:
	DisplayName(table, name) 	: Allows for customizing tab names (IE: Unused libraries)
	Modified?(data) 		 	: Used to determine if the selected item has been modified
	Valid?(data)			 	: Used to set images accordingly
	Synced?(tabData, savedData)	: Used to determine if a tabs data needs to be refreshed
	Editable?(item)				: Used to control if a tab's view control is editable
view must support:
	Get() => object
	Set(object)
	GetState()
	SetState(stateobject)
	Dirty?()
	SetReadOnly(readonly = true)
	GetFirstVisibleLine()		-> Line position
	SetFirstVisibleLine(pos)
	Invalidate()
	AfterSave()
view can support:
	Reset(newData = false, oldData = false) : Allows for customizable reset criteria
*/
PassthruController
	{
	Name: "Explorer"
	redirs: #(On_Cut, On_Copy, On_Paste, On_Delete, On_Select_All, On_Undo, On_Redo)
	New(.modelClass, view, .extraTabMenu = false, .besideTabs = false, .treeArgs = false)
		{
		.viewCtrl = view.Copy()
		.tree = .FindControl(#TreeView)
		.tabsCtrl = .FindControl(#Tabs)
		.tabsCtrl.SetImageList(.imageHandler.ImageResources)
		.Search_vals = Record()
		.sub = PubSub.Subscribe('Redir_SendToEditors', .sendToEditors)
		}

	New2()
		{
		super.New2()
		.model = Construct(.modelClass)
		.imageHandler = ExplorerMultiImageHandler()
		}

	// Receive / echo commands from / to owned editors
	sendToEditors(@args)
		{
		.ForeachTab({ it.Editor.SendToAddons(@args) })
		}

	ResetTheme()
		{ .imageHandler.ResetTheme() }

	viewClass: ExplorerAdapterControl
		{
		Name: 'PlaceholderCtrl'
		Editor: false
		Dummy?: true
		Default(@unused) { return false }
		}

	Getter_View()
		{
		return not .TabConstructed?(.tabsCtrl.GetSelected())
			? .viewClass
			: .tabsCtrl.GetControl()
		}

	TabConstructed?(idx)
		{ return .tabsCtrl.Constructed?(idx) }

	Getter_Tabs()
		{ return .tabsCtrl.Tab }

	Getter_Tree()
		{ return .tree }

	Getter_Model()
		{ return .model }

	Getter_CurItem()
		{ return .curitem }

	Getter_ImageHandler()
		{ return .imageHandler }

	Menu:
		(
		("&File",
			("New &Folder", "Create a new folder")
			("New &Item", "Create a new document")
			("&Delete Item", "Delete the selected item or folder")
			""
			("Print...", "Print the current item")
			""
			("&Close", "Close this window")
			)
		("&Edit",
			("&Undo\tCtrl+Z", "Undo the last action")
			("&Redo", "Redo the last action")
			"",
			("Cu&t\tCtrl+X", "Cut the selected text to the clipboard")
			("&Copy\tCtrl+C", "Copy the selected text to the clipboard")
			("&Paste\tCtrl+V", "Insert the contents of the clipboard")
			)
		("&Help",
			"&Users Manual\tF1",
			"",
			"&About Suneido")
		)
	Controls()
		{
		selectedTabColor = IDESettings.Get('ide_selected_tab_color', false)
		selectedTabBold = IDESettings.Get('ide_selected_tab_bold', true)
		return Object(
			'HorzSplit',
			Object('Vert',
				Object('ExplorerMultiTree').MergeNew(.treeArgs),
				#(Skip, 4)
				xmin: 150, xstretch: 1),
			Object('Vert',
				Object('Tabs',
					close_button: .imageHandler.CloseButton,
					scrollTabs: IDESettings.Get('ide_scroll_tabs', true),
					border: 0, extraControl: .besideTabs, :selectedTabColor,
					:selectedTabBold),
			xstretch: 6))
		}

	On_New_Folder()
		{ .NewItem(true) }

	On_New_Item()
		{ .NewItem(false) }

	NewItem(container?, name = #New, text = '')
		{
		if .curitem is false
			return

		parent = .Container?(.curitem)
			? .curitem
			: .tabsCtrl.GetTabData().parent
		if false is item = .tree.AddNewItem(parent, container?, :name, :text)
			return

		if not .editable?(.getnum(item))
			return

		.View.Dirty?(true)
		.update(item)
		}

	editable?(item)
		{ return .model.Method?(#Editable?) ? .model.Editable?(item) : false }

	On_Delete_Item(confirm = true, allowLibraryDelete? = false)
		{
		if not .deleteItem?(allowLibraryDelete?)
			return
		.updateTabs()
		libraryDeleted? = false
		for selitem in .tree.Selection.Copy()
			{
			if not .deleteChildren?(confirm, selitem)
				continue
			.model.DeleteItem(num: .getnum(selitem), name: .getname(selitem),
				group: .Container?(selitem), path: .tree.Path(selitem))
			if .Static?(selitem)
				libraryDeleted? = true
			else
				.delitem(selitem)
			}
		.CloseTabs({ |i| .keepTab?(i) })
		if libraryDeleted?
			.Reset()
		.invalidateViews()
		}

	deleteItem?(allowLibraryDelete?)
		{
		if .curitem is false
			return false
		else if .Static?(.curitem) and not allowLibraryDelete?
			{
			.AlertInfo('Delete Library', 'To delete libraries:\n\n' $
				'Right-Click a library folder, and select Delete > Delete Library')
			return false
			}
		return true
		}

	deleteChildren?(confirm, selitem)
		{
		return .Children?(selitem) and confirm is true
			? .confirmDeleteChildren(selitem)
			: true // no children to worry about, contiunue
		}

	confirmDeleteChildren(item)
		{
		deleting = .Static?(item) ? #library : #folder
		msg = 'If you delete the ' $ deleting $ ' "' $ .getname(item).Tr('()') $
			'", all the records it contains will also be deleted.\r\n\r\n'
		if deleting is #library
			msg $= 'Library deletions are not staged for Version Control\r\n\r\n'
		return YesNo(msg $ 'Continue?', 'Confirm Delete', .Window.Hwnd, MB.ICONWARNING)
		}

	On_Save()
		{ .updateTabs() }

	deleting: false
	delitem(item)
		{
		.deleting = true
		.tree.DeleteItem(item)
		.deleting = false
		}

	getname(item)
		{ return .tree.GetName(item) }

	getnum(item)
		{ return .tree.GetParam(item) }

	updateTabs()
		{ .forEachConstructedTab({ .update(.tabsCtrl.GetTabData(it).item) }) }

	forEachConstructedTab(block)
		{
		for idx in .. .tabsCtrl.GetTabCount()
			if .TabConstructed?(idx)
				block(idx)
		}

	SaveCode_AfterChange()
		{ .update(.curitem) }

	update(item)
		{
		if false is view = .verifyView(item)
			return

		if not view.Dirty?() or false is num = .getnum(item)
			return
		// Ensure we do not trigger multiple saves per update call
		view.Dirty?(false)

		// Build the record for saving
		x = [].Merge(view.Get())
		x.num = num
		x.name = .getname(item)
		.model.Save(x)

		// Sync the saved changes with the tab data and invalidate the controls
		.syncTabDetails(idx = .tabsCtrl.FindTabBy(#item, item), .itemData(item), false)
		view.AfterSave()
		.invalidateViews(idx)
		}

	verifyView(item)
		{
		if item is false or .deleting
			return false
		if false is idx = .tabsCtrl.FindTabBy(#item, item)
			return false
		if not .TabConstructed?(idx)
			return false
		view = .tabsCtrl.GetControl(idx)
		return not view.GetDefault(#Dummy?, false)
			? view
			: false
		}

	syncTabDetails(idx, data, view)
		{
		.Tabs.SetData(idx, data)
		.Tabs.SetText(idx, data.tabname)
		.tree.SetImage(data.item, data.theme)
		.tabsCtrl.SetImage(idx, data.image)
		if view isnt false
			view.Invalidate()
		}

	invalidateViews(skip = false)
		{
		.forEachConstructedTab()
			{
			if it isnt skip
				.tabsCtrl.GetControl(it).Invalidate()
			}
		}

	TreeView_KeyboardNavigation(olditem)
		{ .selectTreeItemDelayed(olditem) }

	delay: 			500 /* = 1/2 of a second */
	timer: 			false
	lastNewitem: 	false
	firstOlditem: 	false
	selectTreeItemDelayed(olditem)
		{
		if .timer is false
			.firstOlditem = olditem
		.timer = .Delay(.delay, uniqueID: 'keyboardNavigation')
			{
			.selectTreeItem(.firstOlditem, .lastNewitem, focusTree?:)
			.timer = .lastNewitem = .firstOlditem = false
			}
		}

	SelectTreeItem(olditem, newitem)	// From TreeView control
		{
		.lastNewitem = newitem
		if not .resetting and not .tabsChanging? and .timer is false
			.selectTreeItem(olditem, newitem)
		}

	selectTreeItem(olditem, newitem, focusTree? = false)
		{
		if olditem isnt false and olditem isnt 0
			.update(olditem)
		if false isnt idx = .tabsCtrl.FindTabBy(#item, newitem)
			.tabsCtrl.Select(idx)
		else
			.insertNewTab(newitem)
		if focusTree? // Ensure tree does not lose focus when navigating via arrow keys
			.tree.SetFocus()
		}

	noSelect: false
	insertNewTab(newitem)
		{
		if false is data = .itemData(newitem)
			return
		viewCtrl = .viewCtrl.Copy()
		viewCtrl.data = data
		viewCtrl.Tab = tabname = data.tabname
		viewCtrl.readonly = .readonly?(newitem, viewCtrl)
		.tabsCtrl.Insert(tabname, viewCtrl, :data, image: data.image, noSelect: .noSelect)
		.tree.SetImage(data.item, data.theme)
		}

	moveTab?()
		{ return IDESettings.Get('ide_move_tab', true) }

	itemData(item)
		{
		if item is false or false is data = .model.Get(num = .getnum(item))
			return false

		/* Sudo tree items won't have real "nums" provided by the model, for example:
		- LibTreeModel > Library root folders (IE: stdlib)
		- SchemaModel > Never provides "num" as it does not coincide to anything
		*/
		data.num = data.GetDefault('num', num)
		data.item = item
		data.parent = .tree.GetParent(item)
		data.tooltip = data.path = .Getpath(item)
		data.tabname = .displayName(data.table, data.name)
		.imageDetails(data)
		return data
		}

	displayName(table, name)
		{ return .model.Method?(#DisplayName) ? .model.DisplayName(table, name) : name }

	imageDetails(data)
		{
		data.theme = data.GetDefault(#group, true)
			? .folderImage(data, .recordModified?(data))
			: .valid?(data)
				? .recordModified?(data)
					? #Modified
					: #Document
				: #Invalid
		data.image = .imageHandler[data.theme]
		.imageHandler.SetTheme(data, data.theme)
		}

	recordModified?(data)
		{ return .model.Method?(#Modified?) ? .model.Modified?(data) : false }

	folderImage(data, modified?)
		{
		folderState = .tree.Expanded?(data.item) ? #Open : #Closed
		return modified? ? #Folder $ folderState $ #Modified : #Folder $ folderState
		}

	valid?(data)
		{ return .model.Method?(#Valid?) ? .model.Valid?(data) : true }

	readonly?(item, viewCtrl)
		{ return viewCtrl.GetDefault(#readonly, false) or not .editable?(.getnum(item)) }

	TabControl_SelChanging()
		{
		.update(.curitem)
		return false
		}

	prevView: false
	tabsChanging?: false
	TabsControl_SelectTab()
		{
		.tabsChanging? = true
		if .prevView isnt false and .prevView isnt .viewClass
			{
			.RemoveRedir(.prevView)
			.Send(#Explorer_DeselectTab, .prevView)
			}

		if .moveTab?()
			.tabsCtrl.MoveTab(.tabsCtrl.GetSelected(), 0)
		if false isnt .curitem and .tree.GetSelectedItem() isnt .curitem
			.tree.SelectItem(.curitem)

		.Send(#ResetAddons)
		.Send(#SetRedirs)

		if .View.GetDefault(#Editor, false) isnt false
			.View.Editor.UPDATEUI()

		if not .resetting
			.CloseTabs({ |i| .keepTab?(i) })
		.prevView = .View.Copy()
		.tabsChanging? = false
		}

	keepTab?(i)
		{
		data = .tabsCtrl.GetTabData(i)
		return i isnt .tabsCtrl.GetSelected() and .Send(#CloseTab?, data) is true
			? ''
			: .tree.ItemExists?(data.item)
		}

	Refresh(records)
	// warning - if dirty and called prior to save, changes are lost. This is on purpose
		{
		refreshed? = false
		// Copy to ensure each open ExplorerMultiControl refreshes properly
		recs = records.Copy()
		.Tabs.ForEachTab()
			{ |data, idx|
			if recs.Empty?()
				break
			recIdx = recs.FindIf({ it.name is data.name and it.table is data.table })
			if recIdx isnt false
				{
				if .RefreshTab(idx, recs[recIdx].force)
					refreshed? = true
				recs.Delete(recIdx)
				}
			}
		return refreshed?
		}

	RefreshTab(idx, force = false)
		{
		tabData = .tabsCtrl.GetTabData(idx)
		if false is savedData = .itemData(tabData.item)
			return false
		if refresh? = force or not .synced?(tabData, savedData)
			.refreshTab(idx, savedData)
		return refresh?
		}

	synced?(tabData, savedData)
		{ return .model.Method?(#Synced?) ? .model.Synced?(tabData, savedData) : true }

	refreshTab(idx, data)
		{
		if false isnt view = .tabsCtrl.GetControl(idx)
			{
			pos = view.GetFirstVisibleLine()
			view.Set(data) // This may clear out the undo / redo queue
			view.SetFirstVisibleLine(pos)
			}
		else
			{
			ctrlData = .tabsCtrl.TabConstructData(idx)
			ctrlData.data = data
			}
		.syncTabDetails(idx, data, view)
		}

	ResetControls(force = false)
		{
		if .WindowActive?() and not force
			return false
		state = .Send(#GetState)
		.Reset()
		if state isnt 0
			.Send(#SetState, state, resetting:)
		return true
		}

	GotoPath(path, skipFolder? = false)
		{
		if false is item = .getItem(path.Split("/"), skipFolder?)
			return false
		if .curitem is item
			return true
		.update(.curitem)
		if false isnt idx = .tabsCtrl.FindTabBy(#item, item)
			.tabsCtrl.Select(idx)
		else
			.tree.SelectItem(item)
		return true
		}

	getItem(path, skipFolder? = false)
		{
		item = false
		list = .tree.GetChildren(TVI.ROOT)
		for (i = 0; i < path.Size(); i++)
			if false isnt item = .gotoPathItem(list, path, i, skipFolder?)
				{
				.tree.ExpandItem(item)
				list = .tree.GetChildren(item)
				}
			else
				break
		return item
		}

	gotoPathItem(list, path, pathIdx, skipFolder?)
		{
		for item in list
			{
			// don't go to folder
			if path[pathIdx] is .tree.GetName(item) and skipFolder? is true and
				.tree.Container?(item) is true and pathIdx is (path.Size() - 1)
				continue
			if path[pathIdx] is .tree.GetName(item)
				return item
			}
		return false
		}

	Tree_UpdateTabData(item, findBy)
		{
		if false isnt idx = .tabsCtrl.FindTabBy(findBy.field, findBy.value)
			if false isnt data = .itemData(item)
				.syncTabDetails(idx, data, .tabsCtrl.GetControl(idx))
		}

	Getpath(item)
		{ return .tree.Path(item) }

	Rename(item, name)
		{
		if 0 is .tree.Rename(item, name)
			return 0
		if false isnt tab = .tabsCtrl.FindTabBy(#item, item)
			{
			path = .Getpath(item)
			tabname = .displayName(.tree.RootName(item), name)
			tabData = Object(tooltip: path, :path, :tabname, :name)
			.tabsCtrl.SetTabData(tab, tabData, name: tabname)
			.View.SetState(.View.GetState())
			.View.Set(.model.Get(.getnum(item)))
			}
		.invalidateViews()
		return 0
		}

	Inactivate()
		{
		if not .Destroyed?()
			.updateTabs()
		}

	Explorer_Get()
		{ return .Get() }

	Get() // gets the current view
		{ return .model.Get(.getnum(.curitem)) }

	Children?(item)
		{
		if false is num = .getnum(item)
			return false
		return .model.Children?(num)
		}

	Container?(item)
		{ return .model.Container?(.getnum(item)) }

	Static?(item)
		{ return .model.Static?(.getnum(item)) }

	GetTreeModel(num)	// Recursive
		{
		obj = .model.Get(num, origText?:).Copy()
		if obj.group
			{
			// If it is a parent...
			obj.subitems = Object()
			for child in .model.Children(num)
				obj.subitems.Add(.GetTreeModel(child.num))
			}
		return obj
		}

	TabContextMenu(x, y, hover = false)
		{
		tab = hover
		if tab is false
			{
			ScreenToClient(.tabsCtrl.Tab.Hwnd, pt = Object(:x, :y))
			tab = SendMessageTabHitTest(.tabsCtrl.Tab.Hwnd, TCM.HITTEST, 0, Object(:pt))
			}
		if false is tabMenuOb = .tabMenuOb(tab)
			return 0
		menu = .translateMenuOb(tabMenuOb.menu, tabMenuOb.name)
		// Note: i is obtained from a flattened object
		// so menu needs to be flattened for .selectTabMenuItem as well
		if 0 isnt i = ContextMenu(menu).Show(.tabsCtrl.Tab.Hwnd, x, y)
			.selectTabMenuItem(tabMenuOb.menu.Flatten(), i, tab)
		return 0
		}

	tabMenuOb(tab)
		{
		name = ''
		menu = Object()
		if .tabsCtrl.GetTabCount() > 0
			menu.Add('Close &All')
		if .reopen?()
			menu.Add('&Reopen %1')
		if menu.Size() is 0
			return false
		if tab isnt false
			{
			name = .tabsCtrl.TabName(tab)
			menu.Add('&Close %1', 'Close Others', 'Close to the Right', at: 0)
			menu.Add('')
			if Object?(.extraTabMenu)
				menu.Add(@.extraTabMenu)
			}
		return Object(:name, :menu)
		}

	translateMenuOb(menuOb, name)
		{
		return menuOb.Map()
			{
			Object?(it)
				? .translateMenuOb(it, name)
				: TranslateLanguage(it, it.Prefix?('&Reopen')
					? Paths.Basename(.last_closed.Tr('()'))
					: name)
			}
		}

	selectTabMenuItem(tabmenu, i, tab)
		{
		// Ensure tab changes are committed before carrying out context option
		.updateTabs()
		switch option = tabmenu[i-1].Tr('&').Replace(' %1', '')
			{
		case 'Close' :
			.Tab_Close(tab)
		case 'Close Others' :
			.CloseTabs({ |i| i is tab })
		case 'Close to the Right' :
			.CloseTabs({ |i| i is tab }, closeToRight?:)
		case 'Close All' :
			.CloseTabs()
		case 'Reopen' :
			.GotoPath(.last_closed)
			.last_closed = false
		default:
			tabData = .tabsCtrl.GetTabData(tab, '')
			tabData.idx = tab
			.Send('TabMenu_' $ option.Tr('&|" '), tabData)
			}
		}

	reopen?()
		{
		if .last_closed is false
			return false
		for (i = 0; i < .tabsCtrl.GetTabCount(); ++i)
			if .last_closed is .tabsCtrl.GetTabData(i, '').path
				return false
		return true
		}

	Tab_AllowDrag()
		{ return not .moveTab?() }

	Tab_Close(tab)
		{ .closeTab(tab) }

	last_closed: false
	closeTab(i, skipCollapse? = false)
		{
		data = .tabsCtrl.GetTabData(i)
		.last_closed = data.path is false ? '' : data.path
		.update(data.item)
		.tabsCtrl.Remove(.adjustTabsForClosing(i))
		if .tabsCtrl.NoCtrls?()
			{
			.Send(#DeleteRedirs)
			.prevView = false
			.Send(#ResetAddons)
			.tree.SelectItem(NULL)
			}
		if not skipCollapse?
			.tree.ExpandItem(data.item, true)
		}

	adjustTabsForClosing(i)
		{
		origTab = i
		if 1 isnt .tabsCtrl.GetTabCount() and i is .tabsCtrl.GetSelected()
			.tabsCtrl.Select(i = i + (i is 0 ? 1 : -1))
		return .moveTab?() ? i : origTab
		}

	collapse()
		{
		.collapsenode(0)
		if not .tabsCtrl.NoCtrls?()
			.tree.EnsureVisible(.curitem)
		}

	CloseTabs(ignoreTab = function(i /*unused*/) { false }, closeToRight? = false)
		{
		// ignoreTab:
		//		returns true: 	Do not close tab,
		//		returns false: 	Close tab
		//		returns '': 	Close tab, do not collapse folders
		tabClosed? = false
		for (i = .tabsCtrl.GetTabCount() - 1; i >= 0; i--)
			{
			ignore = ignoreTab(:i)
			if ignore is true
				{
				if closeToRight?
					break
				continue
				}
			.closeTab(i, skipCollapse?:)
			if ignore is false
				tabClosed? = true
			}
		if tabClosed?
			.collapse()
		}

	collapsenode(parent)
		{
		// Collapse all nodes whose hierarchies don't involve open tabs
		collapsethis? = false is .tabsCtrl.FindTabBy(#parent, parent)
		.tree.ForEachChild(parent)
			{ |x|
			collapse = .collapsenode(x)
			if collapse and .Container?(x)
				.tree.ExpandItem(x, true)
			collapsethis? = collapsethis? and collapse
			}
		return collapsethis?
		}

	getter_curitem()
		{ return .tabsCtrl.GetTabData().item }

	CloseActiveTab()
		{ .closeTab(.tabsCtrl.GetSelected()) }

	GetSelected()
		{ return .curitem }

	RootSelected?()
		{ return .tree.GetParent(.curitem) is 0 }

	resetting: false
	Reset(model = false)
		{
		.resetting = true
		if model isnt false
			.modelClass = model
		.tree.Reset(.model = Construct(.modelClass))
		.resetTabs()
		.resetting = false
		}

	resetTabs()
		{
		.syncTabs(pathLookup?:)
		.collapse()
		.tree.SelectItem(.curitem)
		}

	syncTabs(skip = false, pathLookup? = false)
		{
		// Must go in the reverse order to avoid issues with closing tabs
		for (idx = .tabsCtrl.GetTabCount() - 1; idx >= 0; --idx)
			if idx isnt skip
				.syncTab(idx, .tabsCtrl.GetTabData(idx), :pathLookup?)
		}

	syncTab(idx, tabData, pathLookup? = false)
		{
		item = not pathLookup?
			? tabData.item
			: .findTabItem(tabData.path)
		if false is savedData = .itemData(item) // Record has been deleted
			.closeTab(idx, skipCollapse?:)
		else if not .synced?(tabData, savedData) // Record has been modified
			.refreshTab(idx, savedData)
		else
			.syncTabDetails(idx, savedData, .tabsCtrl.GetControl(idx))
		}

	// As .Reset reconstructs the model, each tab has to update their "item" member.
	// To do this, we do a reverse lookup using the path stored in the tabs data.
	findTabItem(path)
		{
		pathOb = path.Split('/')
		if pathOb.Empty?()
			return false
		if false is item = .getItem(pathOb)
			{
			pathOb[0] = pathOb[0].Has?('(') ? pathOb[0].Tr('()') : '(' $ pathOb[0] $ ')'
			if false is item = .getItem(pathOb)
				return false
			}
		return item
		}

	RestoreState(state)
		{
		if false isnt tabs = state.GetDefault('tabs', false)
			.restoreTabs(tabs, state.GetDefault('activeTabPath', false))
		if state.Member?('splitterpos')
			.HorzSplit.SetSplit(state.splitterpos)
		}

	restoreTabs(tabs, activeTabPath)
		{
		if .moveTab?() and activeTabPath isnt false
			{
			tabs.Remove(activeTabPath)
			tabs.Add(activeTabPath, at: 0)
			}
		tabs.Each()
			{
			.noSelect = it isnt activeTabPath
			if false is .GotoPath(it)
				.Send('Explorer_RestoreTab', it)
			}
		.noSelect = false
		// Handling for when the "activeTabPath" no longer associates with a record
		if .tabsCtrl.GetControl() is false and .tabsCtrl.GetTabCount() isnt 0
			.tabsCtrl.Select(0)
		.tree.SelectItem(.curitem)
		}

	GetTabsPaths(all? = false, skipFolder? = false)
		{
		tabs = Object()
		.Tabs.ForEachTab()
			{ |tab, idx|
			if not all? and idx is 10 /*= 10 most recent tabs*/
				break
			if skipFolder? and tab.group is true
				continue
			tabs.Add(tab.path)
			}
		return tabs
		}

	ForeachTab(block)
		{ .forEachConstructedTab({ block(.tabsCtrl.GetControl(it)) }) }

	Destroy()
		{
		.sub.Unsubscribe()
		.imageHandler.Destroy()
		super.Destroy()
		}
	}