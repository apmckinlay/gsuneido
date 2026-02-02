// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// TODO: only add How Do I button if there are any items
Controller
	{
	Xmin: 700
	Ymin: 500
	login: false
	help_book: false
	book: false

	notificationSubscription: false
	dblockedSubscription: false
	userBookMessageSubscription: false

	New(book = false, title = false, start = "Cover",
		login = false, help_book = false, .noWikiNotes = false)
		{
		super(.controls(book, title, help_book))
		.InitializeBook(book, start, login, help_book)
		}

	InitializeBook(book, start, login = false, help_book = false)
		{
		if book is false
			return false
		.marksplit = .Vert.SplitMarks
		.treesplit = .marksplit.SplitTree
		.browser = .treesplit.Vert.BookBrowser.Browser
		.bookbrowser = .treesplit.Vert.BookBrowser
		.toolbar = .Vert.Horz.GetDefault(#Toolbar, .nullToolbar)
		.help_book = help_book
		.tabs = .treesplit.Vert.Tab
		.tree = .treesplit.TreeView
		.marks = .marksplit.BookMarkContainer
		.howdobutton = .FindControl('How_Do_I')
		.start = start
		.login = login
		.setbuttons()
		.redir_plugin_buttons()
		.Redir('On_Find')
		.Redir('On_Find_Next')
		.Redir('On_Find_Previous')
		.startPage = Object(path: "", name: start, loaded: false)
		if QueryEmpty?(book, path: "", name: start)
			.startPage.name = .model.First().name
		if not help_book
			Suneido.CurrentBook = book
		.book = book
		BookLog(.book $ ':open')
		if not Suneido.Member?("OpenBooks")
			Suneido.OpenBooks = Object()
		Suneido.OpenBooks[book] = this
		.subscriptions()
		.Defer(.startUp, uniqueID: "BookSetTitle") // for applications to override title
		.icon = .SetIcon(.book, .Window.Hwnd)
		return true
		}

	subscriptions()
		{
		.subs = Object()
		if .help_book is false
			{
			TimerManager.Start()
			.subs.Add(
				PubSub.Subscribe('book notification',
					{ .Defer(.bookNotification, uniqueID: 'BookNotification') }),
				PubSub.Subscribe('book dblocked',
					{ .Defer(.dblocked, uniqueID: 'BookDBlocked') }),
				PubSub.Subscribe('book userBookMessage',
					{ .Defer(.userBookMessage, uniqueID: 'BookUserBookMessage') })
				)
			}
		}

	startUp()
		{
		.goto(@.startPage)
		.startPage.loaded = true

		.app_SetTitle()
		if false isnt c = .FindControl('search')
			c.SetFocus()

		menuOptions = GetContributions(.book $ '_WindowSubMenu')
		menuOptions.Add(Object(root: 'Show/Hide', order: 0, options: options = Object()))
		.splitOptions.Each({ options[it.text] = it.cmd })
		.Window.AddWindowMenuOptions(menuOptions.Sort!({|x, y| x.order <= y.order }))
		}

	getter_splitOptions()
		{
		return .splitOptions = Object(
			Object(text: 'Bookmarks', cmd: .toggleSplit(.marksplit))
			Object(text: 'Tree View', cmd: .toggleSplit(.treesplit)))
		}

	toggleSplit(split)
		{
		return { split.Open? ? split.Close() : split.Open() }
		}

	bookNotification()
		{
		if false is (vert = .FindControl('notifyVert'))
			return
		if vert.Tally() > 0 // new button is already shown
			return
		vert.Insert(0, #('EnhancedButton', 'NEW !', command: 'NewMessages',
			tip: 'View New Messages', buttonStyle:, mouseEffect:,
			weight: 'bold', pad: 25, textColor: 0x0000ff))
		}

	setPageHeadMessage(msg)
		{
		if false is c = .FindControl('pagehead')
			return

		if msg isnt ''
			{
			c.Set(msg)
			c.SetColor(CLR.RED)
			return
			}
		c.Set(PageHeadName())
		c.SetColor(CLR.BLACK)
		}

	dblocked()
		{ .setPageHeadMessage('DATABASE LOCKED - NO CHANGES ALLOWED') }

	userBookMessage()
		{
		msg = Suneido.CurrentUserBookMessage
		if not msg.Empty?()
			.setPageHeadMessage(msg[0].subject $ '\r\n' $ msg[0].notes)
		else
			.setPageHeadMessage('')
		}

	UseIcon(book)
		{
		iconFns = GetContributions('BookIcon').Reverse!()
		for iconFn in iconFns
			if '' isnt icon = iconFn(book)
				return icon
		throw "BookIcon contribution should be defined."
		}

	SetIcon(book, hwnd)
		{
		iconfile = .UseIcon(book)
		icon = LoadImage(NULL, iconfile, IMAGE.ICON, 0, 0, LR.SHARED | LR.LOADFROMFILE |
			LR.DEFAULTSIZE)
		if icon isnt 0
			SendMessage(hwnd, WM.SETICON, ICON.BIG, icon)
		return icon
		}

	GetIcon()
		{ return .icon }

	nullToolbar: class
		{
		EnableButton(@unused)
			{ }

		GetState()
			{ return [] }

		SetState(unused)
			{ }
		}

	On_NewMessages()
		{
		ShowBookNotifications()
		.FindControl('notifyVert').Remove(0)
		}

	Startup()
		{
		login = .login
		if (String?(login))
			login = login.Eval() // Eval - not sure what has been used for login arg
		if (Function?(login))
			if (false is login(hwnd: .Window.Hwnd))
				.Window.Destroy()
		}

	controls(book, title, help_book)
		{
		.book = book
		.wikiNotes = .noWikiNotes ? false : OptContribution(.book $ 'WikiNotes', false)
		.Title = title isnt false ? title : book isnt false ? book : ""
		if book is false
			return #(Fill)
		.model = BookModel(book)
		.chapters = .model.Children()
		tabs = .buildTabs()
		horz = 	Object('Horz', ystretch: 0)
		if help_book is true
			horz.Add(Object('Toolbar', Object('Back' drop:), Object('Forward' drop:),
				'PreviousBookPage', 'NextBookPage', Object('Print' drop:)))
		horz2 = Object('Horz')
		if help_book isnt true
			horz2.Add('Skip', ['Vert', ['Static' PageHeadName()
				font: '@serif', weight: 'bold', size: '+6', name: 'pagehead']])
		.addBookButtons(horz2, help_book, book)
		horz.Add(Object('Vert'
			#(EtchedLine before: 0, after: 1 xstretch: .001)
			#(Skip 3),
			horz2,
			name: 'Vert2'))
		return .layout(horz, tabs)
		}

	buildTabs()
		{
		tabs = Object('Tab2', orientation: 'top')
		for (i = .chapters.Size() - 1; i >= 0; --i)
			{
			c = .chapters[i]
			if (c.name is "Cover" or c.name is "res")
				.chapters.Delete(i)
			else
				tabs.Add(c.name, at: 1)
			}
		return tabs
		}

	addBookButtons(horz2, help_book, book)
		{
		horz2.Add('Fill')
		if Suneido.User is 'default'
			horz2.Add(
				Object('EnhancedButton', 'Refresh', textColor: CLR.Highlight,
					buttonStyle:, mouseEffect:, tip: 'Reload the current page (F5)'),
				'Skip',
				Object('EnhancedButton', 'Go To', pad: 20, textColor: CLR.Highlight,
					buttonStyle:, mouseEffect:, tip: 'Go To Definition'),
				'Skip')
		if .wikiNotes isnt false and
			not Libraries().Has?(.wikiNotes.GetDefault('ExcludeForLib', false))
			horz2.Add(Object('Button', 'Wiki', tip: book $ ' Wiki Notes'), 'Skip')
		// "NEW" button will be inserted into "notifyVert" (in BookNotification)
		horz2.Add(#(Vert name: 'notifyVert'), 'Skip')

		if help_book is true
			horz2.Add(['BookSearch' book])
		else if help_book isnt true
			{
			horz2.Add(#IM_Button)
			.add_plugin_buttons(horz2)
			if TableExists?(.book $ 'Help')
				{
				if TableExists?(.book $ 'HelpHowToIndex')
					horz2.Add(Object('MenuButton', 'How Do I?'), 'Skip')
				horz2.Add(Object('EnhancedButton', text: 'Help',
					image: 'questionMark_black.emf',
					weight: 'bold', pad: 20, buttonStyle:, mouseEffect:
					imageColor: CLR.Highlight, mouseOverImageColor: CLR.Highlight,
					textColor: CLR.Highlight, imagePadding: 0.15))
				}
			}
		horz2.Add(#(Skip 5))

		horz2.Add(Object('EnhancedButton', command: 'ExtraMenu', image: 'menu',
			mouseEffect:, tip: 'More options', imagePadding: 0.15), 'Skip')
		}

	On_ExtraMenu(source)
		{
		if 0 < i = MenuButtonControl.PopupMenu(.extraMenuDetails.menu, source.Hwnd)
			(.extraMenuDetails.cmds[.extraMenuDetails.cmdMap[i - 1]])(book: this)
		}

	getter_extraMenuDetails()
		{
		groups = Object().Set_default(Object())
		submenus = Object()
		cmds = Object()
		GetContributions(.book $ '_ExtraBookMenu').
			Add([text: 'Show/Hide', submenu: .splitOptions]).
			Each({ .processMenuOption(it, groups, submenus, cmds) })

		menu = Object()
		groups.Members().Sort!().Each()
			{|groupIdx|
			if menu.NotEmpty?()
				menu.Add('')
			groups[groupIdx].Sort!().Each()
				{|option|
				menu.Add(option)
				if submenus.Member?(option)
					menu.Add(submenus[option])
				}
			}
		return .extraMenuDetails = [:menu, :cmds, cmdMap: menu.Flatten()]
		}

	processMenuOption(ob, groups, submenus, cmds)
		{
		group = ob.GetDefault('group', 5 /*= default middle group*/)
		groups[group].Add(ob.text)
		if ob.Member?('submenu')
			{
			submenu = Object()
			ob.submenu.Each()
				{
				submenu.Add(it.text)
				cmds[it.text] = it.cmd
				}
			submenus[ob.text] = submenu.Sort!()
			}
		else
			cmds[ob.text] = ob.cmd
		}

	layout(horz, tabs)
		{
		Object('Vert',
			horz,
			#(Skip 3),
			Object("BookSplit"
				Object("BookSplit"
					Object("BookTree" model: .model xstretch: 1)
					Object("Vert"
						tabs
						#(BookBrowser, Browser)
						xstretch: 3)
					"west"
					splitname:	"Tree"
					name:		"SplitTree"
					)
				#(BookMarkContainer)
				"east"
				splitname:	"Bookmarks"
				name:		"SplitMarks"
				)
			#(Skip 5)
			)
		}

	url(path, name)
		{ return "suneido:/" $ .book $ path $ "/" $ name }

	plugin_buttons: #()
	add_plugin_buttons(horz)
		{
		.plugin_buttons = Object()
		Plugins().ForeachContribution('BookButtons', 'addButton')
			{|c|
			if c.Member?(#control)
				horz.Add(c.control)
			else
				{
				horz.Add(Object('Button', c.name, tip: c.tip), 'Skip')
				.plugin_buttons.Add(c)
				}
			}
		}

	redir_plugin_buttons()
		{
		for c in .plugin_buttons
			.Redir("On_" $ ToIdentifier(c.name), c.target)
		}

	Commands:
		(
		(Back			"Alt+Left"	"Move back to previous page in history list", Left)
		(Forward		"Alt+Right"	"Move forward to next page in history list", Right)
		(PreviousBookPage	"Alt+Up"	"Move backward to previous page in book", Up)
		// cannot use Next, it conflicts with AccessControl Next command
		// WindowBase.Mapcmd only allows one cmdnum (required by windows accelerator) per method name
		(NextBookPage	"Alt+Down"	"Move forward to next page in book", Down)
		(Find			"Ctrl+F"	"Find pages containing text")
		(Find_Next,		"F3", 		"Find the next occurrence in the current item")
		(Find_Previous,	"Shift+F3",	"Find the previous occurrence in the current item")
		(Copy, 			"Ctrl+C")
		(Go_To, 		"F12")
		(Help			"F1")
		(Refresh		"F5")
		(Screenshot		"Ctrl+Alt+S")
		)

	Goto(address) // from tree & bookmarks
		{
		if (address is .CurrentPage())
			return
		name = address.Extract("[^/]*$")
		path = address[.. -name.Size() - 1]

		if .startPage.loaded is false
			{
			.startPage.path = path
			.startPage.name = name
			}
		else
			.goto(path, name)
		}

	goto(path, name)
		{
		dest = .url(path, name)
		if (not .bookbrowser.ProgramPage?() and
			dest is Url.Decode(.browser.Current()))
			return
		.browser.Goto(dest)
		}

	prevUrl: false
	Going(url) // sent by Browser
		{
		for contrib in Contributions('BookGoing')
			contrib(this, url)
		if url is "about:startup"
			return false
		address = .url_to_name(url)
		if false isnt pos = address.FindLast('?term')
			address = address[..pos]
		if false is current = .model.Get(address)
			return false
		.current = current
		if .prevUrl isnt url
			BookLog(url.Replace('^suneido:/'))
		.prevUrl = url
		for (c in .chapters.Members())		// Select the correct tab
			if (address $ "/" =~ "^/?" $ .chapters[c].name $ "/")
				{
				.tabs.Select(c)
				break
				}
		.marks.GotoPath(address)		// Select the correct book mark (if one exists)
		.tree.GotoPath(address)		// Select the correct tree item
		.setbuttons()					// Enable/disable navigation buttons
		.app_SetTitle(address)
		return 0
		}

	current: false
	CurrentPage() // used by bookmarks
		{ return Object?(.current) ? .current.path $ "/" $ .current.name : "" }

	setbuttons()
		{
		// Enable/disable navigation buttons
		.toolbar.EnableButton("Back", .browser.CanBack?())
		.toolbar.EnableButton("Forward", .browser.CanForward?())
		.toolbar.EnableButton("Print", true)

		// don't show "How Do I?" button if there's no index
		if .howdobutton isnt false and .MenuButton_How_Do_I().Empty?()
			.howdobutton.SetEnabled(false)
		else if .howdobutton isnt false and not .MenuButton_How_Do_I().Empty?()
			.howdobutton.SetEnabled(true)
		}

	ProgramPage()
		{
		.toolbar.EnableButton("Print", false)
		Suneido.CurrentBookOption = .current.path $ '/' $ .current.name
		.app_SetTitle(Suneido.CurrentBookOption)
		}

	app_SetTitle(option = "")
		{
		if false isnt title = BookWindowTitle(.Book, option)
			.Window.SetTitle(title)
		}

	TabControl_SelChanging()
		{ return .bookbrowser.PageFrozen?() }

	SelectTab(i, source)
		{
		if (source is .tabs)
			.goto("", .chapters[i].name)
		}

	TabClick(i, source)
		{
		if source is .tabs
			{
			if Object?(.current) and .current.name isnt '' and
				.current.path.Prefix?('/' $ .chapters[i].name $ '/')
				.goto(.current.path.BeforeLast('/'), .current.path.AfterLast('/'))
			else
				.goto("", .chapters[i].name)
			}
		}

	On_NextBookPage()
		{
		if not .help_book
			return
		x = .model.Next(.current)
		if (x is false)
			{ Beep(); return }
		.goto(x.path, x.name)
		}

	On_PreviousBookPage()
		{
		if not .help_book
			return
		x = .model.Prev(.current)
		if (x is false)
			{ Beep(); return }
		.goto(x.path, x.name)
		}

	On_Back()
		{
		if not .help_book
			return
		.Back()
		}
	Back()
		{
		.browser.GoBack()
		}

	url_to_name(url)
		{
		url = url.RemovePrefix('suneido:/')
		return url[url.Find('/')..]
		}

	urls_to_names(urls)
		{
		names = Object()
		for (url in urls)
			names.Add(.url_to_name(url))
		return names
		}

	Drop_Back(r, source)
		{
		list = .urls_to_names(.browser.BackList())
		i = ContextMenu(list).Show(source.Hwnd, r.left, r.bottom)
		if (i > 0)
			.browser.GoBack(i)
		}

	On_Forward()
		{
		if not .help_book
			return
		.browser.GoForward()
		}

	Drop_Forward(r, source)
		{
		list = .urls_to_names(.browser.ForwardList())
		i = ContextMenu(list).Show(source.Hwnd, r.left, r.bottom)
		if (i > 0)
			.browser.GoForward(i)
		}

	Drop_Print(r, source)
		{
		page = .CurrentPage()
		if (page is "/Cover" or page is "/Contents")
			page = ""
		children? = not QueryEmpty?(.book, path: page)
		disabled = children? ? #() : #("Print Section", "Preview Section")

		printOptions = Sys.SuneidoJs?()
			? #("Print Page", "Print Section")
			: #("Print Page", "Print Section", "Preview Section", "Page Setup")
		options = printOptions.Map()
			{
			Object(name: it, state: disabled.Has?(it) ? MFS.DISABLED : MFS.ENABLED)
			}
		i = ContextMenu(options).Show(source.Hwnd, r.left, r.bottom)
		if (options.Member?(i - 1))
			this['On_' $ options[i - 1].name.Tr(' ', '_')]()
		}

	On_Print()
		{ .On_Print_Page() }

	On_Print_Page()
		{ .browser.Print() }

	On_Print_Section()
		{ .print(OLECMDID.PRINT) }

	On_Preview_Section()
		{ .print(OLECMDID.PRINTPREVIEW) }

	print(cmd)
		{
		start = .CurrentPage()
		if (start is "/Cover" or start is "/Contents")
			start = "" // whole book

		if Sys.SuneidoJs?()
			{
			DoWithSelectedFileName('print.html')
				{ |filename|
				BookPrint(.book, start, filename, cmd)
				}
			}
		else
			{
			filename = GetAppTempFullFileName("su")
			htmlFilename = filename.BeforeLast('.') $ '.html'
			BookPrint(.book, start, htmlFilename, cmd)
			.browser.Goto("file:///" $ htmlFilename)
			.Delay(1000) /*= 1 sec */
				{
				.delete_temp(filename)
				.delete_temp(htmlFilename)
				}
			}
		}

	delete_temp(filename)
		{ DeleteFile(filename) }

	On_Page_Setup()
		{ .browser.PageSetup() }

	searchHwnd: false
	On_Find()
		{
		BookLog('Help Find')
		if (.searchHwnd is false)
			.searchHwnd = Window(Object("BookSearch", this, .book),
				w: 220, h: 400, keep_placement:).Hwnd
		else
			WindowActivate(.searchHwnd)
		}

	On_Refresh()
		{
		if .current isnt false
			.browser.Goto(.url(.current.path, .current.name))
		}

	On_Wiki()
		{
		if .wikiNotes is false or .current is false
			return

		url = `http://` $ .wikiNotes.Server $ ':' $ .wikiNotes.Port $ '/Wiki?' $
			.wikiNotes.PagePrefix $ .buildWikiNotesURLName(.current.path, .current.name)
		ShellExecute(.WindowHwnd(), 'open', Url.Encode(url))
		}

	buildWikiNotesURLName(path, name)
		{
		nonAlphaNumeric = '^a-zA-Z0-9'
		name = name.Tr(nonAlphaNumeric)
		path = path.Tr(nonAlphaNumeric)
		return path $ name
		}

	BookRefresh()
		{ .On_Refresh() }

	On_Screenshot()
		{ OptContribution('BookScreenshot', {})(hwnd: .Window.Hwnd) }

	On_Go_To()
		{
		if Suneido.User isnt 'default'
			return

		bookedit_rec = Query1(.book, num: .current.num)
		html_page? = bookedit_rec.text =~ '^Biz_Menu|^BookMenuPage|^GetBookPage' or
			BookContent.Match(.book, bookedit_rec.text)
		if .Title is TranslateLanguage("User's Manual")
			{
			.gotoFromUserManual(bookedit_rec, html_page?)
			return
			}

		name = (bookedit_rec.text.Prefix?('#(')
			? bookedit_rec.text.Tr('#()').BeforeFirst(',')
			: bookedit_rec.text).BeforeFirst('(')
		.gotoRecord(html_page?, name)
		}

	gotoFromUserManual(bookedit_rec, html_page?)
		{
		name = line = false
		libview_rec = Gotofind(.current.name, exact:)
		if not libview_rec.Empty?() and bookedit_rec isnt false
			{
			if false is (goto = (.goto_dialog)(.Window.Hwnd))
				return
			html_page? = html_page? and goto is 'BookEdit'
			name = .current.name
			line = 1
			}
		.gotoRecord(html_page?, name, line, libview_rec)
		}

	goto_dialog: Controller
		{
		Title: "Go To Definition"
		CallClass(hwnd)
			{ return ToolDialog(hwnd, Object(this)) }

		Controls: #(Horz
			(Button LibraryView) Skip (Button BookEdit) Skip (Button Cancel))

		On_LibraryView()
			{ .Window.Result('LibraryView') }

		On_BookEdit()
			{ .Window.Result('BookEdit') }
		}

	gotoRecord(html_page?, name, line = false, libview_rec = false)
		{
		if html_page?
			.goto_bookedit()
		else
			GotoLibView(:name, :line, list: libview_rec)
		}

	goto_bookedit()
		{
		if false isnt page = Query1(.book, num: .current.num)
			OpenBook(.book, .book $ page.path $ '/' $ page.name, bookedit?:)
		}

	BookSearchClosed()
		{ .searchHwnd = false }

	logfonts: #()
	setFontState(state)
		{
		if not .help_book and state.Member?("logfont") and
			state.logfont isnt false
			{
			if state.logfont.Member?("fontPtSize")
				state.logfont.lfHeight = StdFonts.LfSize(state.logfont.fontPtSize)
			else
				state.logfont.fontPtSize = StdFonts.PtSize(state.logfont.lfHeight)
			SetGuiFont(.LogFont = state.logfont)
			}
		.logfonts = Object()
		for m in state.Members()
			if m.Suffix?("logfont")
				.logfonts[m] = state[m]
		}

	SetState(state)
		{
		// need to set font before opening book
		.setFontState(state)
		if (state.Member?("book"))
			{
			_parent = .Parent
			_ctrlspec = #()
			// should destroy Fill
			// do NOT call .Destroy since that means it's called twice
			// which messes up other stuff
			start = state.Member?("start") ? state.start : "Contents"
			title = .Title isnt "" ? .Title : state.Member?("title") ? state.title : false
			help_book = state.Member?("help_book") ? state.help_book : false
			if .book is false
				{
				.Initialize(.controls(state.book, title, help_book))
				.InitializeBook(state.book, start, :help_book)
				}
			SetWindowText(.Window.Hwnd, .Title)
			.sizeWindow()
			}
		else
			return
		if (state.Member?("tree"))
			.tree.SetState(state.tree)				// Restore tree state
		if (state.Member?("marks"))
			{
			if state.marks.Size() > 0
				{
				// remove bookmarks that don't match with paths in the book
				for (i = state.marks.Size() - 1; i >= 0; --i)
					{
					bookmark = state.marks[i]
					if not .ValidBookmark?(String?(bookmark) ? bookmark
						: bookmark.path)
						state.marks.Delete(i)
					}
				}
			.marks.SetState(state.marks)			// Restore bookmark state
			}
		if (state.Member?("toolbar"))
			.toolbar.SetState(state.toolbar)		// Restore toolbar state
		if (state.Member?("marksplit"))
			.marksplit.SetSplit(state.marksplit)	// Restore bookmark splitter state
		if (state.Member?("treesplit"))
			.treesplit.SetSplit(state.treesplit)	// Restore tree splitter state
		if (state.Member?("login"))
			{ .login = state.login; .Startup() }
		}

	sizeWindow()
		{
		if Sys.SuneidoJs?()
			return

		GetClientRect(.Window.Hwnd, rc = Object())
		SendMessage(.Window.Hwnd, WM.SIZE, WMSIZE.RESTORED,
			rc.right | rc.bottom << 16) /* = bit shift */
		}

	ValidBookmark?(bookmark)
		{
		return not QueryEmpty?(.book, path: bookmark.BeforeLast('/'),
			name: bookmark.AfterLast('/'))
		}

	GetState()
		{
		state = Object()
		if this.Member?('BookControl_book')
			state.book = .book
		if this.Member?('Title')
			state.title = .Title
		if this.Member?('BookControl_tree')
			state.tree = .tree.GetState()
		if this.Member?('BookControl_marks')
			state.marks = .marks.GetState()
		if this.Member?('BookControl_toolbar')
			state.toolbar = .toolbar.GetState()
		if this.Member?('BookControl_start')
			state.start = .start
		if this.Member?('BookControl_marksplit')
			state.marksplit = .marksplit.GetSplit()
		if this.Member?('BookControl_treesplit')
			state.treesplit = .treesplit.GetSplit()
		if this.Member?('BookControl_login')
			state.login = .login
		if this.Member?('BookControl_help_book')
			state.help_book = .help_book
		if not .logfonts.Empty?()
			state.Merge(.logfonts)
		if not .help_book
			state[GetCustomFontSaveName()] = .LogFont
		return state
		}

	Getter_LogFont()
		{ .LogFont = Suneido.logfont }

	Getter_Browser()
		{ return .browser }

	Getter_Book()
		{
		// occasionally app_SetTitle will call this
		// when the book is not initialized or is destroyed
		if not this.Member?('BookControl_book')
			{
			SuneidoLog('INFO: .book reference not in BookControl')
			return ''
			}
		return .book
		}

	On_Help()
		{ OpenBook(.book $ 'Help', .current) }

	lrucache: false
	MenuButton_How_Do_I()
		{
		if .lrucache is false
			.lrucache = LruCache(.getfn)
		return .lrucache.Get(.current).Members().Sort!()
		}

	getfn(current)
		{
		if current is false or
			not TableExists?(.book $ 'HelpHowToIndex')
			return #()
		query = {|name| Query1(.book $ "HelpHowToIndex
			where name = " $ Display(name)) }
		if false is x = query(current.name)
			x = query(current.path $ '/' $ current.name)
		if x is false
			return #()
		x.howtos = x.howtos.RemoveIf(
			{ 'hidden' is BookEnabled(.book $ 'Help', it) })
		ob = Object()
		for ht in x.howtos
			ob.Add(ht.BeforeLast('/') at: ht.AfterLast('/'))
		return ob
		}

	On_How_Do_I(name)
		{
		BookLog("How Do I > " $ name)
		path = .lrucache.Get(.current)[name]
		OpenBook(.book $ 'Help', Object(:path, :name))
		}

	Ok_to_CloseWindow?()
		{
		// .current may be false if start page can not be found and user
		// closes window before going to any other pages
		if .current is false or .current.path is "" or
			.bookbrowser.ProgramPage?() is false
			return true
		GetWindowPlacement(.Window.Hwnd,
			wpPlace = Object(length: GetWindowPlacementSize()))
		if wpPlace.showCmd is SW.SHOWMINIMIZED
			return true
		.Goto(.current.path)
		return false
		}

	MenuPage?()
		{
		return Query1(.book, num: .current.num).text.Prefix?('Biz_Menu(')
		}

	HelpBook?()
		{
		return .help_book
		}

	Destroy()
		{
		BookLog(.book $ ':close')
		.subs.Each(#Unsubscribe)

		Plugins().ForeachContribution('BookOptions', 'close')
			{|c|
			if .book is c.book
				(c.func)()
			}

		if (.searchHwnd isnt false)
			DestroyWindow(.searchHwnd)

		// Suneido.OpenBooks will not be initialized if book was false
		// (see 'return' in New)
		if Suneido.Member?("OpenBooks")
			Suneido.OpenBooks.Delete(.book)
		if .help_book is false
			TimerManager.Stop()
		super.Destroy()
		}
	}
