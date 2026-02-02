// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
AccessBase
	{
	Name: "Access1"
	New(.query, control = false, title = false,
		validField = false, .protectField = false, .option = 'Access')
		{
		super(.makecontrols(query, control, title))
		.data = .Vert.TitleScroll.Control
		.ValidField = validField
		.Data.SetProtectField(.protectField)
		.edit_button = .Vert.Horz.Edit
		.get_plugins()
		.lock = AccessLock(this, false)
		.load_record()
		.attachmentsManager = AttachmentsManager(.query, Object())
		}
	Commands:
		(
		("Edit",	"Alt+E")
		("Restore",	"Alt+R")
		("NextTab",	"Ctrl+Tab")
		("PrevTab",	"Shift+Ctrl+Tab")
		)
	makecontrols(query, control, title)
		{
		if title is false
			title = query
		.Title = title $ " - Access"
		.query_columns = QueryColumns(query)
		if control is false
			{
			control = Object('Vert')
			for field in .query_columns
				control.Add(field)
			}
		return Object('Vert'
			Object('TitleScroll', title, Object('Record', control)),
			#(Horz
				(EnhancedButton text: 'Edit' tip: 'Alt+E' xmin: 80
					buttonStyle:, mouseEffect:)
				(Button 'Restore' tip: 'Alt+R' xmin: 80))
			'Status')
		}

	get_plugins()
		{
		.afterfield_plugins = Object().Set_default(#())
		Plugins().ForeachContribution('Access1', false)
			{ |x|
			if x[2] isnt .option
				continue
			x = x.Copy() // so plugins can store stuff in it
			if String?(x.func)
				x.func = x.func.Compile()
			fields = Object?(x.fields) ? x.fields : (x.fields)()
			for f in fields
				{
				if not .query_columns.Has?(f)
					throw "invalid Access plugin field: " $ f
				if x[1] is 'AfterField'
					.afterfield_plugins[f].Add(x)
				else
					throw "invalid Access plugin"
				}
			}
		}

	afterfield_plugins: ()
	Record_NewValue(field, value)
		{
		if .query_columns.Has?(field)
			{
			if .afterfield_plugins.Member?(field)
				for x in .afterfield_plugins[field]
					(x.func)(field, value, .Window.Hwnd, .query, .GetData())

			.Send("Access1_AfterField", field, value)
			}
		}

	// control like SearchConfigControl needs to know if it is from Access1
	GetAccess1RecordControl()
		{
		return .GetRecordControl()
		}

	GetRecordControl()
		{
		return .data
		}

	GetData()
		{
		return .Data.Get()
		}

	original_record: false
	GetOriginal()
		{
		return .original_record
		}

	nextnumControl: false
	NextNumberControl_Register(source)
		{
		.nextnumControl = source
		}

	On_Restore()
		{
		.NotifyObservers("restore")
		.load_record()
		.lock.Unlock()
		if .nextnumControl isnt false
			.nextnumControl.RestoreNextNumber()
		.Send("Access1_Restore")
		.attachmentsManager.ProcessQueue(restore?:)
		}
	load_record()
		{
		.NotifyObservers('before_setdata')
		x = Query1(.query)
		.Send('Access_BeforeRecord', x)
		.setdata(x)
		.NotifyObservers('setdata')
		.Data.SetReadOnly(true)
		.edit_button.Pushed?(false)
		}
	setdata(x)
		{
		.new_record? = x is false
		if x is false
			x = Record()
		.original_record = x.Copy()
		if .ValidField isnt false
			x[.ValidField]
		if .protectField isnt false
			x[.protectField]
		.Data.Set(x)
		.Status.SetValid()
		.Status.Set('')
		}
	On_Edit()
		{
		if ReadOnlyAccess(this) is true
			return

		// saving will switch the edit mode off if successful,
		// otherwise it will be left in edit mode
		if .Data.Dirty?() and .EditMode?()
			{
			.Save()
			return
			}

		if not .EditMode?()
			.editMode()
		else
			.viewMode()
		}

	viewMode()
		{
		.edit_button.Pushed?(false)
		.Data.SetReadOnly(true)
		.lock.Unlock()
		}

	editMode()
		{
		.load_record() // re-read record in case it has changed since loaded
		if not .lock.Trylock()
			return
		.edit_button.Pushed?(true)
		.Data.SetReadOnly(false)
		.FocusFirst(.Vert.TitleScroll.ScrollHwnd)
		}

	ScrollToView(field)
		{
		if not .selectPromptText(field)
			{
			field.SetFocus()
			field.TopDown('On_Select_All')
			}

		if Sys.SuneidoJs?() // js scrolls to the focused field automatically
			return

		.Delay(500, .scrollToView) /*= 1/2 second */
		}
	scrollToView()
		{
		scroll = .Vert.TitleScroll.Vert.Scroll
		fieldRect = GetWindowRect(GetFocus())
		scrollRect = GetWindowRect(.Vert.TitleScroll.ScrollHwnd)
		if scrollRect.bottom < fieldRect.top and fieldRect.top < scrollRect.bottom and
			scrollRect.bottom < fieldRect.bottom and fieldRect.bottom < scrollRect.bottom
			return
		dis = fieldRect.bottom - scrollRect.bottom + 30 /*= bottom scroll bar */
		scroll.Scroll(0, -dis)
		}

	selectPromptText(field)
		{
		if not field.Parent.Base?(PairControl)
			return false
		prompt = field.Parent.GetChildren()[0]
		if not prompt.Base?(StaticControl)
			return false
		prompt.SetFocus()
		// static control on javascript does not support focus,
		// and Component.FocusFirst steals focus after the text is selected
		.Delay(100, prompt.On_Select_All) /*= 1/10 sec */
		return true
		}

	GetLockKey()
		{
		// there will be no key value in the table since most Access1Controls are
		// on a table/query with an empty key (table only has one record).
		// Using the query as the key value for now unless we can find something better
		return .query
		}

	GetTitle()
		{
		return .Title
		}

	EditMode?()
		{
		return not .Data.GetReadOnly()
		}
	SetEditMode()
		{
		if not .EditMode?()
			.On_Edit()
		return .EditMode?()
		}
	Access1_Save()
		{
		.Save()
		}
	Save()
		{
		if not .Data.Dirty?()
			return true

		if false is .CheckValid(evalRule?:)
			return false

		if .save() is false
			return false

		.NotifyObservers('after_save')
		.original_record = Query1(.query)
		.Data.SetReadOnly(true)
		.lock.Unlock()
		.edit_button.Pushed?(false)
		.Data.Dirty?(false)
		.Send('Access1AfterSaving')
		.new_record? = false
		.attachmentsManager.ProcessQueue()
		return true
		}
	save()
		{
		try
			return KeyExceptionTransaction()
				{ |t|
				if false is .Send('AccessBeforeSave', :t) or
					false is .NotifyObservers('save', t)
					{
					if not t.Ended?()
						t.Rollback()
					return false
					}
				return .new_record? ? .save_output(t) : .save_update(t)
				}
		catch (unused, "interrupt: KeyException")
			return false
		}
	save_output(t)
		{
		t.QueryOutput(.query, .Data.Get())
		return true
		}
	save_update(t)
		{
		x = t.Query1(.query)
		if RecordConflict?(.original_record, x, .query_columns, .Window.Hwnd)
			return false
		x.Update(.Data.Get())
		return true
		}
	ChangeQuery(.query)
		{
		.query_columns = QueryColumns(query)
		.load_record()
		}
	QueueDeleteAttachmentFile(newFile, oldFile, name, action)
		{
		return .attachmentsManager.QueueDeleteFile(
			newFile, oldFile, .GetData(), name, action)
		}
	Destroy()
		{
		.lock.Unlock()
		super.Destroy()
		}
	}
