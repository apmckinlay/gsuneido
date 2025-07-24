// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// option could be a book option (ie. '/Menu/Option'), a class name
	// (ie. Fake_Screen), or false
	CallClass(option, user = false)
		{
		.validateIfDefaultUser(option, user)
		option = .getOption(option)
		permFn = .getPermissionFn(option, user)
		return permFn(option, user)
		}

	validateIfDefaultUser(option, user)
		{
		if .defaultUser?(user) and option isnt false and not .customTableOption?(option)
			.getPermissionClass().VerifyOption(option)
		}

	defaultUser?(user)
		{
		permUser = user
		if user is false
			permUser = Suneido.User
		return permUser is 'default'
		}

	getOption(option)
		{
		return option isnt false ? option : Suneido.GetDefault(#CurrentBookOption, false)
		}

	getPermissionFn(option, user)
		{
		// automatically give access to default users, custom field accesses, and
		// screens that are not menu options
		if .defaultUser?(user) or .customTableOption?(option)
			return .yesPermission
		if option is false
			return .noPermission
		return option.Prefix?('/') ? .bookPermission : .accessGotoPermission
		}

	customTableOption?(option)
		{
		return String?(option) and option =~ "^Access_custom_\d\d\d\d\d\d$"
		}

	yesPermission(@unused)
		{
		return true
		}

	noPermission(@unused)
		{
		return false
		}

	bookPermission(option, user)
		{
		permissionClass = .getPermissionClass()
		return permissionClass.Permission(option, user)
		}

	accessGotoPermission(option, user)
		{
		permissionClass = .getPermissionClass()
		permission = false

		QueryApply(BookOptionQuery(FindCurrentBook(), option))
			{
			option = it.path $ "/" $ it.name
			if true is result = permissionClass.Permission(option, user)
				return true
			if permission is false
				permission = result
			}
		return permission
		}

	getPermissionClass()
		{
		contrib = Contributions('AccessPermissions').Last()
		return contrib.cl
		}

	// Used by FindCurrentBook
	GetDefaultBook()
		{
		return .getPermissionClass().DefaultBook
		}
	}
