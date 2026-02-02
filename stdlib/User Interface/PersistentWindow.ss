// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Window
	{
	hSaveTimer: 0
	Master?: false
	option: ''

	New(control = false, master? = false, stateobject = #(), newset = false,
		.option = '')
		{
		super(control: .Create(control, stateobject, newset), show: false, :newset)
		// Perform persistence functions
		.ensureSuneidoPersistentMember()
		// Note that stateobject.master overrides this
		.Master? = (newset is false) ? master? : true
		// Set state, or show window normally
		if .IsStateData?(stateobject)
			.SetState(stateobject)
		else if false isnt state = .loadState()
			.SetState(state)
		else
			// Show normally
			.Show(SW.NORMAL)
		// Add 'this' to global persistent object list
		Suneido.Persistent.Windows.Add(this)
		}
	ENDSESSION()
		{
		if .Master?
			{
			.saveAllStates()
			Exit(true) // immediate forced exit
			}
		return 0
		}

	destroying?: false
	DESTROY()
		{
		if .destroying?
			return 0

		.destroying? = true
		// If the master is being closed, save all windows
		if .Master?
			{
			// Exit, saving persistence state:
			if .hSaveTimer isnt 0
				KillTimer(.Hwnd, .hSaveTimer)
			.saveAllStates()
			.runExitContribs()
			Suneido.Persistent.Exit = true
			.exit()
			}
		// Remove from global persistent object list
		.ensureSuneidoPersistentMember()
		c = Suneido.Persistent.Windows.Find(this)
		if c isnt false
			Suneido.Persistent.Windows.Delete(c)
		if not Suneido.Persistent.Exit
			.saveState()
		// Inherited DESTROY
		return super.DESTROY()
		}
	exit()
		{
		params = .collectWindows()
		params.calls = FormatCallStack(GetCallStack(limit: 10), levels: 10)
		BookLog('exit', :params)
		Exit()
		}

	collectWindows()
		{
		list = Object()
		try
			{
			block = {|hwnd, unused|
				title = GetWindowText(hwnd)
				list.Add([title, hwnd, IsWindowVisible(hwnd), IsWindowEnabled(hwnd)])
				true // continue enumerating
				}
			EnumThreadWindows(GetCurrentThreadId(), block, NULL)
			ClearCallback(block)
			}
		catch (e)
			list.e = e
		return list
		}

	runExitContribs()
		{
		try
			for func in Contributions('PersistentWindowExit')
				func()
		catch (e)
			SuneidoLog("ERROR: PersistentWindowExit" $ e)
		}
	ensureSuneidoPersistentMember(set = false)
		{
		if not Suneido.Member?('Persistent')
			Suneido.Persistent = Object(Exit: false Windows: Object())
		if set isnt false
			Suneido.Persistent.Set = set
		if not Suneido.Persistent.Member?('Set')
			SuneidoLog('ERROR: Should not missing persistent set', calls:)
		}
	EnsureTableOK()
		{
		// Ensure persistent table exists and has proper format
		Database('ensure persistent (keynum, classname, option, classstate, pos,
			last, set, user)
			index(user, set) key(keynum)
			index(set)')
		}
	IsStateData?(x)
		{
		return x.Member?('classname') and String?(x.classname) and
				x.Member?('classstate') and Object?(x.classstate) and
				x.Member?('pos') and Object?(x.pos)
		}
	IsPersistentWindow?(x)
		{
		return BaseClass(x) is this.Base()
		}
	state_pos: ()
	GetState()
		{
		// Get window placement state
		GetWindowPlacement(.Hwnd, wpPlace = Object(length: GetWindowPlacementSize()))
		// Return window rectangle, show command and control's custom state data
		// Note: pos.master contains either false (not the master)
		// or an object containing libraries to use
		state = Object(
			classname: Display(.Ctrl.Base()),
			classstate: .Ctrl.GetState(),
			pos: Object(
				rect: wpPlace.rcNormalPosition,
				showcmd: .showCmd(wpPlace),
				master: (.Master?) ? Libraries() : false
				).MergeNew(.state_pos), // keep other resolutions
			last: false,
			set: Suneido.Persistent.Set,
			option: .option
			user: Suneido.User_Loaded
			)
		state.pos[.resolution()] =  state.pos.rect.Copy()
		state.pos[.resolution()].showcmd = state.pos.showcmd
		return state
		}
	showCmd(wpPlace)
		{
		// don't restore "minimized" state - confusing to users
		return wpPlace.showCmd isnt SW.SHOWMINIMIZED ? wpPlace.showCmd : SW.NORMAL
		}
	SetState(stateobject)
		{
		.state_pos = stateobject.pos // keep other resolutions
		res = .resolution()
		if stateobject.pos.Member?(res)
			{
			stateobject.pos.rect = stateobject.pos[res]
			stateobject.pos.showcmd =
				stateobject.pos[res].GetDefault(#showcmd, stateobject.pos.showcmd)
			}
		.loadFont(stateobject)
		// Restore the control's internal state
		.Ctrl.SetState(stateobject.classstate)
		// Ensure that only one master exists
		// Note: pos.master contains either false (not the master)
		// or an object containing libraries to use
		if stateobject.pos.master isnt false
			{
			.Master? = true
			for i in Suneido.Persistent.Windows.Members()
				if .IsPersistentWindow?(Suneido.Persistent.Windows[i]) and
					Suneido.Persistent.Windows[i].Master? is true
					{
					.Master? = false
					break
					}
			}
		// Restore the window's position and placement state
		GetWindowPlacement(.Hwnd, wpPlace = Object(length: GetWindowPlacementSize()))
		wpPlace.rcNormalPosition = stateobject.pos.rect
		wpPlace.showCmd = stateobject.pos.showcmd
		// if size > screen size then maximize
		r = stateobject.pos.rect
		w = r.right - r.left
		h = r.bottom - r.top
		xmax = GetSystemMetrics(SM.CXMAXIMIZED)
		ymax = GetSystemMetrics(SM.CYMAXIMIZED)
		if w > xmax or h > ymax
			wpPlace.showCmd = SW.MAXIMIZE
		SetWindowPlacement(.Hwnd, wpPlace)
		}
	loadFont(stateobject)
		{
		mLogName = GetCustomFontSaveName()
		if stateobject.classstate.Member?(mLogName)
			stateobject.classstate.logfont = stateobject.classstate[mLogName]
		}
	resolution()
		{
		wa = GetWorkArea()
		return (wa.right - wa.left) $ 'x' $ (wa.bottom - wa.top) $ '@' $ GetDpiFactor()
		}
	saveAllStates()
		{
		.ensureSuneidoPersistentMember()
		.EnsureTableOK()
		RetryTransaction()
			{|t|
			// Delete previous open persistent windows from this set
			t.QueryDo('delete ' $ .query(set: Suneido.Persistent.Set, last: "false"))
			keynum = .Next(t)
			// Write each control to persistent table
			q = t.Query('persistent')
			for w in Suneido.Persistent.Windows
				{
				new_record = w.GetState()
				new_record.keynum = keynum++
				q.Output(new_record)
				}
			}
		}
	Load(set = "IDE")
		{
		n = 0
		// Save persistent set name
		.ensureSuneidoPersistentMember(set)
		// Load all Persistent Windows from database table 'persistent'
		.EnsureTableOK()
		Transaction(read:)
			{|t|
			n = .load(t, .query(:set, last: "false"))
			if n is 0
				n = .load(t, .query(defaultuser:, :set, last: "false"))
			// check for a workspace in the set where last is "true" IF no windows loaded
			if n is 0
				n = .loadHiddenWorkspace(t, set)
			}
		// If no windows were loaded, create a master WorkSpace control:
		if n is 0
			PersistentWindow(#(WorkSpace) newset: set)
		// Return number of windows loaded from database
		return n
		}
	load(t, query)
		{
		n = 0
		t.QueryApply(query)
			{|restore_data|
			if .IsStateData?(restore_data)
				{
				PersistentWindow(stateobject: restore_data,
					option: restore_data.option)
				n++
				}
			}
		return n
		}

	loadHiddenWorkspace(t, set)
		{
		x =	t.QueryFirst(.query(:set, last: "true",
			classname: "WorkSpaceControl /* stdlib class : Controller */", sort?:))
		if x is false
			x = t.QueryFirst(.query(defaultuser:, :set, last: "true",
				classname: "WorkSpaceControl /* stdlib class : Controller */", sort?:))
		if x isnt false and .IsStateData?(x)
			{
			PersistentWindow(master?:, stateobject: x, option: x.option)
			return 1
			}
		return 0
		}

	newSet(set)
		{
		.EnsureTableOK()
		if not QueryEmpty?(.query(:set))
			throw 'persistent set already exists'
		// Save current persistent set
		.CloseSet()
		// Create a new set if a string for the "newset" argument of .New was given
		.ensureSuneidoPersistentMember(set)
		// Set current user name:
		if Suneido.User_Loaded isnt "default"
			{
			Alert("You were previously logged-in as \"" $ Suneido.User_Loaded $
				"\".\n\nYou are now logged-in as the default user.", "New Set",
				flags: MB.ICONINFORMATION)
			Suneido.User = Suneido.User_Loaded = "default"
			}
		}
	CloseSet()
		{
		if not Suneido.Member?('Persistent')
			return
		.saveAllStates()
		// Close all current persistent windows...
		Suneido.Persistent.Exit = true
		for w in Suneido.Persistent.Windows.Copy()
			{
			if w.Master?
				w.Master? = false	// So master doesn't close Suneido
			w.Destroy()
			}
		Suneido.Persistent.Exit = false
		}
	Create(control, stateobject, newset)
		{
		.EnsureTableOK()
		// Check if a new set should be started...
		if newset isnt false
			.newSet(newset)
		// Return a value for the "control" argument of Window.New(...)
		if control isnt false
			ctrl = control
		else if .IsStateData?(stateobject)
			{
			.restoreLibs(stateobject)
			ctrl = stateobject.classname.Eval() // needs Eval - has class info comment
			}
		else
			ctrl = WorkSpaceControl
		.runStartupContribs()
		return ctrl
		}
	restoreLibs(stateobject)
		{
		if not Object?(stateobject.pos.master)
			return
		for library in stateobject.pos.master
			try
				Use(library)
			catch (err)
				Alert(err, "", flags: MB.ICONWARNING)
		LibraryTags.Reset()
		}
	runStartupContribs()
		{
		try
			for func in Contributions('PersistentWindowStartup')
				func()
		catch (e)
			SuneidoLog('ERROR: PersistentWindowStartup ' $ e)
		}
	Next(t) // public since it's used by Persistent_DemoData
		{
		return t.QueryMax('persistent', 'keynum', 0) + 1
		}
	saveState()
		{
		.EnsureTableOK()
		Transaction(update:)
			{|t|
			t.QueryDo('delete ' $ .query(set: Suneido.Persistent.Set,
				classname: Display(.Ctrl.Base()), option: .option))
			new_record = .GetState()
			new_record.last = true
			new_record.keynum = .Next(t)
			t.QueryOutput('persistent', new_record)
			}
		}
	loadState()
		{
		.EnsureTableOK()
		Transaction(update:)
			{|t|
			if false isnt x = t.QueryFirst(.query(set: Suneido.Persistent.Set,
				last: "true", classname: Display(.Ctrl.Base()), option: .option, sort?:))
				x.Delete()
			}
		return x
		}
	query(defaultuser = false, set = false, last = false, classname = false,
		option = false, sort? = false)
		{
		query = 'persistent where user = ' $
			Display(defaultuser ? "default" : Suneido.User_Loaded)
		if set isnt false
			query $= ' and set = "' $ set $ '"'
		// Last MUST be the STRING "true" or the string "false",
		// or the boolean value false
		if last isnt false
			query $= ' and last = ' $ last
		if classname isnt false
			query $= ' and classname = "' $ classname $ '"'
		if option isnt false
			query $= ' and option = ' $ Display(option)
		if sort?
			query $= ' sort keynum'
		return query
		}
	}
