// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// e.g. AccessGoTo('Ar_Invoices', 'arivc_num_new', false)
class
	{
	CallClass(access, goto_field, goto_value, hwnd = 0, newRecord = false,
		defaultSelect = false, window = 'Modal', onDestroy = function () {},
		editMode = false)
		{
		_accessGoTo? = true
		result = .CheckPermission(access)
		accessClass = .getAccessClass(result)
		permission = result.permission

		if permission is false
			{
			Alert("You do not have permission to access this option.", "Access",
				hwnd, MB.ICONWARNING)
			return false
			}
		if not Suneido.Member?("AccessGoToCount")
			Suneido.AccessGoToCount = 0
		if Suneido.AccessGoToCount >= 2
			return false
		.logBookAccessGoTo(access, ':open')
		++Suneido.AccessGoToCount
		book_option = false
		module = CurrentModule()
		acc = Class?(accessClass[0]) ? accessClass[0] : accessClass
		if acc.Member?('AccessPermission')
			book_option = .findPermittedOption(acc)
		else if false isnt rec = QueryMatchBookOption(Suneido.CurrentBook, access, module)
			book_option = rec.path $ '/' $ rec.name

		accessOb = Object(.wrapper, accessClass, readOnly: permission is 'readOnly',
			:goto_field, :goto_value, :newRecord, :book_option, :defaultSelect, :editMode)

		result = .openAccess(window, accessOb, hwnd, access, onDestroy)

		if window isnt 'Modal'
			{
			--Suneido.AccessGoToCount
			.logBookAccessGoTo(access, ':close')
			}
		return result
		}

	getAccessClass(result)
		{
		accessClass = result.accessCl
		if Object?(accessClass)
			{
			// accessCl is read only, need to make it modifiable to set accessGoTo? flag
			accessClass = accessClass.Copy()
			accessClass.accessGoTo? = true
			}
		else
			accessClass = Object(accessClass, accessGoTo?:)

		return accessClass
		}

	findPermittedOption(accessClass)
		{
		permittedOption = false
		for option in accessClass.AccessPermission
			{
			if option is 'All'
				continue
			if true is perm = .getPermission(option)
				return option
			else if perm is 'readOnly'
				permittedOption = option
			}
		return permittedOption
		}

	modalOnDestroy(access, onDestroy)
		{
		--Suneido.AccessGoToCount
		.logBookAccessGoTo(access, ':close')
		onDestroy()
		}

	CheckPermission(access)
		{
		accessCl = false
		permission = false
		try
			{
			accessCl = .getAccess(access)
			if accessCl.Member?('AccessPermission')
				{
				for option in accessCl.AccessPermission
					{
					if option is 'All'
						{
						permission = true
						break
						}
					if true is permission = .getPermission(option)
						break
					}
				}
			else
				permission = .getPermission(access)
			}
		catch (err)
			{
			SuneidoLog('ERROR: ' $ err)
			permission = false
			}
		return Object(:accessCl, :permission)
		}

	// overridden for test
	getAccess(access)
		{
		return access.Eval() // Eval is okay here, would have to determine if Global or ob
		}

	// overridden for test
	getPermission(option)
		{
		AccessPermissions(option)
		}

	logBookAccessGoTo(access, text)
		{
		// do not attempt to BookLog access
		// if there are no books currently open (app not running in a book)
		if not Suneido.Member?("OpenBooks") or Suneido.OpenBooks.Empty?()
			return
		BookLog(access $ text)
		}

	openAccess(window, accessOb, hwnd, access, onDestroy)
		{
		result = false
		switch window
			{
		case 'Modeless' :
			result = Window(accessOb, title: 'Access', keep_placement:, :onDestroy)
		case 'Dialog' :
			result = ToolDialog(hwnd, accessOb, title: 'Access', keep_size: access)
		default :
			result = ModalWindow(accessOb, title: 'Access', keep_size:, useDefaultSize:,
				onDestroy: { .modalOnDestroy(access, onDestroy) })
			}
		return result
		}

	wrapper: Controller
		{
		AccessGoToBookOption: false

		New(@args)
			{
			super(@.setup(args))

			if false is ctrl = .FindControl('Access')
				return

			if args.readOnly is true
				ctrl.SetReadOnly(true)
			if Object?(args.defaultSelect)
				ctrl.SetSelect(args.defaultSelect)
			if (args.goto_value isnt false)
				.AccessGoto(args.goto_field, args.goto_value, wrapper: this)
			if args.newRecord isnt false and ctrl.Method?('On_New')
				{
				child = .GetChild()
				if child.Member?('AccessGotoNew')
					child.AccessGotoNew()
				fillinData = Object?(args.newRecord)
					? args.newRecord
					: false
				.Defer({ ctrl.On_New(:fillinData) })
				}
			if args.editMode isnt false and ctrl.Method?('SetEditMode')
				.Defer(ctrl.SetEditMode)
			}

		setup(args)
			{
			.AccessGoToBookOption = args.GetDefault('book_option', false)
			return args
			}

		// this is called by Send message
		AccessGoTo_CurrentBookOption()
			{
			return .AccessGoToBookOption
			}

		On_Cancel()
			{
			// need this to override the default method in Control
			// otherwise "close" doesn't save or validate
			// however, it also disables ESC key to close
			}
		}
	}
