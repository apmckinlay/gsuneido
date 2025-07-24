// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// used by GotoLibView and FindReferences
class
	{
	CallClass(name, libs = false, exact = false)
		{
		if name is '' or libs is #()
			return #()

		libs = .libs(libs).Remove(#Builtin)
		list = Object()
		for item in .FindItems(name, exact)
			for lib in libs
				.find(lib, item, list, exact)
		return list
		}

	find(lib, item, list, exact)
		{
		.ForEach(lib, item, exact)
			{ |x|
			try
				{
				path = LibRecGetPath(x, lib)
				if not Libraries().Has?(lib)
					path = path.Replace(lib, '(' $ lib $ ')', 1)
				list.AddUnique(path)
				}
			catch (e)
				Print('ERROR: Gotofind', x.name, lib, e)
			}
		}

	ForEach(lib, item, exact, block, includeText? = false)
		{
		name = item.Replace('^<lib>', lib.Capitalize())
		.query(name, lib, exact, block, includeText?)
		}

	query(name, lib, exact, block, includeText? = false)
		{
		project = includeText? is false ? ' project num, parent, name, group' : ''
		if exact is true
			{
			if false isnt x = Query1(lib $ project,
				group: -1, :name)
				block(x)
			}
		else
			QueryApply(lib $
				' where name >= ' $ Display(name) $
				' and name < ' $ Display(name $ '_z') $ project,
				group: -1)
				{ |x|
				if name is LibraryTags.RemoveTagsFromName(x.name)
					block(x)
				}
		}

	libs(libs)
		{
		return libs is false
			? Libraries().Reverse!().MergeUnion(.libraryTables())
			: libs.Copy().Reverse!()
		}
	libraryTables() // split out for tests
		{
		return LibraryTables()
		}

	FindItems(name, exact = false)
		{
		if not exact
			name = LibraryTags.RemoveTagsFromName(name)

		list = Object(name)
		if .addProtectDefinitions?(name, exact)
			.addFieldRecords(list, name.AfterFirst('_').Replace('__protect$', ''))

		if .onlyExactMatches?(name)
			exact = true

		if .libraryPrefix?(name)
			name = name.AfterFirst('_')

		if not exact
			.addGenericChecks(name, list)

		return list
		}

	addProtectDefinitions?(name, exact)
		{
		return name =~ '^(Rule_|Field_)' and not exact
		}

	onlyExactMatches?(name)
		{
		return name =~ '^(Rule_|Field_|Trigger|Table_)'
		}

	libraryPrefix?(name)
		{
		return name.Has?('_') and
			.libraryTables().Map(#Capitalize).Has?(name.BeforeFirst('_'))
		}

	addGenericChecks(name, list)
		{
		if name =~ '(Control|Format)$'
			list.Add('<lib>_' $ name)
		else
			{
			list.Add(name $ 'Control', name $ 'Component', name $ 'Format',
				'Trigger_' $ name, '<lib>_' $ name, 'Table_' $ name)
			.addFieldRecords(list, name.UnCapitalize())
			}
		}

	addFieldRecords(list, fieldName)
		{
		list.AddUnique('Rule_' $ fieldName)
		list.AddUnique('Field_' $ fieldName)
		list.AddUnique('Rule_' $ fieldName $ '__protect')
		}
	}