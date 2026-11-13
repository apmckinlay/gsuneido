// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// option could be a book option (ie. '/Menu/Option'), a class name
	// (ie. Fake_Screen), or false
	CallClass(option, user = false, userRoles = false)
		{
		Assert(user is false or userRoles is false)
		.validateIfDefaultUser(option, user, userRoles)
		option = .getOption(option)
		permFn = .getPermissionFn(option, user, userRoles)
		return permFn(option, user, :userRoles)
		}

	validateIfDefaultUser(option, user, userRoles)
		{
		if .defaultUser?(user, userRoles) and option isnt false and
			not .customTableOption?(option)
			.getPermissionClass().VerifyOption(option)
		}

	defaultUser?(user, userRoles)
		{
		if userRoles isnt false
			return false
		permUser = user
		if user is false
			permUser = Suneido.User
		return permUser is #default
		}

	getOption(option)
		{
		return option isnt false ? option : Suneido.GetDefault(#CurrentBookOption, false)
		}

	getPermissionFn(option, user, userRoles = false)
		{
		// automatically give access to default users, custom field accesses, and
		// screens that are not menu options
		if .defaultUser?(user, userRoles) or .customTableOption?(option)
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

	bookPermission(option, user, userRoles)
		{
		permissionClass = .getPermissionClass()
		return permissionClass.Permission(option, user, :userRoles)
		}

	accessGotoPermission(option, user, userRoles)
		{
		permissionClass = .getPermissionClass()
		permission = false

		QueryApply(BookOptionQuery(FindCurrentBook(), option))
			{
			option = it.path $ '/' $ it.name
			if true is result = permissionClass.Permission(option, user, :userRoles)
				return true
			if permission is false
				permission = result
			}
		return permission
		}

	getPermissionClass()
		{
		contrib = Contributions(#AccessPermissions).Last()
		return contrib.cl
		}

	// Used by FindCurrentBook
	GetDefaultBook()
		{
		return .getPermissionClass().DefaultBook
		}
	}
