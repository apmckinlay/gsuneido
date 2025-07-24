// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Diff to Overridden'
	CallClass(name, orig_lib, orig_text)
		{
		lib_rec_map = .fetch_other_defintions(name, orig_lib)
		if lib_rec_map.Empty?()
			.AlertInfo(.Title, orig_lib $ ':' $ name $ ' - no other definitions found')
		else
			ToolDialog(0, Object(this, name, orig_lib, orig_text, lib_rec_map))
		}

	New(.name, .orig_lib, .orig_text, .lib_rec_map)
		{
		.diffPane = .FindControl('diffPane')
		.list = .FindControl('ChooseList')
		if .list isnt false
			.list.SelectItem(0)
		else
			.NewValue(.lib_rec_map[0].label)
		}

	fetch_other_defintions(name, orig_lib)
		{
		libs = LibraryTables()
		used = Libraries().Intersect(libs) // handle LibraryTables overridden
		libs = used.MergeUnion(libs) // put used libs first
		found_libs = Object()
		pureName = LibraryTags.RemoveTagFromName(name)
		libs.Each()
			{
			Gotofind.ForEach(it, pureName, exact: false, includeText?:)
				{ |rec|
				if not (it is orig_lib and rec.name is name)
					{
					tag = LibraryTags.GetTagFromName(rec.name)
					found_libs.Add(Object(lib: it, name: rec.name,
						label: it $ Opt(': ', tag.RemovePrefix('__')), :rec))
					}
				}
			}
		return found_libs
		}

	Controls()
		{
		ctrls = Object('Vert')
		if .lib_rec_map.Size() > 1
			ctrls.Add(#(Skip small:),
				Object('Horz', #('Static', 'Library: '),
					Object('ChooseList', list: .lib_rec_map.Map({ it.label }))),
				#(Skip small:))
		ctrls.Add(#(Horz (Pane Fill) name: 'diffPane' ystretch: 3))
		return ctrls
		}

	lastValue: false
	NewValue(value)
		{
		idx = .lib_rec_map.FindIf({ it.label is value })
		if idx is false or (value is .lastValue and .lastValue isnt false)
			return
		.lastValue = value

		.diffPane.RemoveAll()
		listOld = .lib_rec_map[idx].rec.lib_current_text
		listNew = .orig_text
		titleLength = 20
		titleOld = (.lib_rec_map[idx].lib $ ':' $ .lib_rec_map[idx].name).
			RightFill(titleLength, '\t')
		titleNew = (.orig_lib $ ':' $ .name ).RightFill(titleLength, '\t')
		.diffPane.Insert(0, Object('Diff2', listNew, listOld, .orig_lib, .name,
			:titleOld, :titleNew))
		}
	}
