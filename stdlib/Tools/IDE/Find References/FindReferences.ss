// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// used by: FindReferencesControl
// uses: Gotofind, FindReferencesPos
class
	{
	CallClass(orig_name, referencesOnly = false,
		excludeTests = false, excludeUpdates = false, excludeLibs = #())
		{
		orig_name = LibraryTags.RemoveTagFromName(orig_name)
		libs = LibraryTables()
		used = Libraries().Intersect(libs) // handle LibraryTables overridden
		libs = used.MergeUnion(libs) // put used first
		list = Object()
		if not referencesOnly
			.definitions(libs, used, orig_name, list)
		libs = libs.Difference(excludeLibs)
		if list.NotEmpty?()
			list.Add(Object())
		basename = .basename(orig_name)
		.references(libs, used, orig_name, basename, list, excludeTests, excludeUpdates)
		.book_references(basename, list)
		return Object(:list, :basename)
		}

	definitions(libs, used, orig_name, list)
		{
		for item in Gotofind.FindItems(orig_name)
			for lib in libs
				Gotofind.ForEach(lib, item, exact: false)
					{ |rec|
					list.AddUnique(Object(
						Location: rec.name,
						Table: used.Has?(lib) ? lib : '(' $ lib $ ')',
						Found: '',
						Folder: .libRecFolder(rec, lib)))
					}
		}

	max_refs: 50
	references(libs, used, orig_name, base_name, list,
		excludeTests = false, excludeUpdates = false)
		{
		exclude_str = ""
		if excludeTests
			exclude_str $= ' and not name.Suffix?("Test") '
		if excludeUpdates
			exclude_str $= ' and not name.Prefix?("Update_20") '
		for lib in libs
			QueryApply(.referencesQuery(lib, orig_name, base_name, exclude_str))
				{|x|
				list.Add(Object(
					Location: x.name,
					Table: used.Has?(lib) ? lib : '(' $ lib $ ')',
					Found: x.text.LineAtPosition(x.findpos).Trim(),
					Pos: x.findpos,
					Folder: .libRecFolder(x, lib)))
				if list.Size() > .max_refs
					{
					list.Add(Object(Location: "TOO MANY"))
					return
					}
				}
		}

	referencesQuery(lib, origName, baseName, extraWhere = '')
		{
		return lib $ '
			where group is -1 and name isnt ' $ Display(origName) $ extraWhere $ '
			extend findpos = FindReferencesPos(text, name, ' $ Display(baseName) $
				',' $ Display(origName) $ ', ' $ Display(.textSearch?(origName)) $ ')
			where findpos isnt false
			sort name'
		}

	libRecFolder(rec, lib)
		{
		return LibRecGetPath(rec, lib).RemovePrefix(lib).RemoveSuffix(rec.name)
		}

	basename(name)
		{
		if name.Lower().Suffix?('.js')
			return name.BeforeLast('.js')
		if name.Has?('_') and
			LibraryTables().Map(#Capitalize).Has?(name.BeforeFirst('_'))
			return name.AfterFirst('_') // contribution
		if name =~ '^Rule_'
			return name.AfterFirst('_').Replace('__protect$', '')
		if name in ('Control', 'Component', 'Format')
			return name
		for p in #('^Field_', '^Trigger_', 'Control$', 'Component', 'Format$', '^Plugin_',
			'^Table_')
			if name =~ p
				return name.Replace(p, "")
		return name
		}

	textSearch?(name)
		{
		return name.Prefix?('Addon_')
		}

	book_references(name, list)
		{
		sep? = list.NotEmpty?() and list.Last().NotEmpty?()
		for book in BookTables()
			QueryApply(book $
				' where ' $ .notBinaryRecord() $ ' and text.Has?(' $ Display(name) $ ')')
				{|x|
				if sep?
					{
					list.Add(Object())
					sep? = false
					}
				list.Add(Object(Table: book, Location: x.path $ '/' $ x.name))
				}
		}

	// assumes binary record names have three letter extensions
	notBinaryRecord()
		{
		return 'name !~ "\.[a-zA-Z][a-zA-Z][a-zA-Z]$"'
		}

	AllOccurrences(lib, recName, basename, name)
		{
		if lib is '' or not LibraryTables().Has?(lib)
			return ''

		if false is rec = Query1(lib, name: recName, group: -1)
			return ''
		text = rec.lib_current_text
		locs = Object()
		for pos in FindReferencesPos.FindAllPos(text, recName, basename, name,
			.textSearch?(name))
			{
			line = text.LineFromPosition(pos) + 1 /*= to align*/
			line = (line $ ':').RightFill(8) /*= two tab indents*/
			locs.Add(line $ text.LineAtPosition(pos).Trim())
			}
		return locs.Join('\r\n')
		}

	DefinitionExists?(libs, origName)
		{
		libs = LibraryTables().Intersect(libs)
		for item in Gotofind.FindItems(origName)
			for lib in libs
				Gotofind.ForEach(lib, item, exact: false)
					{ |unused|
					return true
					}
		return false
		}

	ReferenceExists?(origName, excludeLibs = #(), extraWhere = '')
		{
		baseName = .basename(origName)
		return LibraryTables().Remove(@excludeLibs).Any?()
			{ |lib|
			not QueryEmpty?(.referencesQuery(lib, origName, baseName, extraWhere).
				RemoveSuffix('sort name'))
			}
		}
	}
