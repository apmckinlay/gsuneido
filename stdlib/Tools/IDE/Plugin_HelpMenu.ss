// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
#(
Contributions:
	(
	(UI, menu, menu: Help, order: 9)

	(UI, action, menu: Help, name: "&Users Manual",
		target: function (@args) { GotoSelectedUserManual(@args) })
	(UI, action, menu: Help, name: "&About Suneido",
		target: function (@unused) { AboutSuneido() })

	(UI, attach, menu: Help, to: WorkSpace)
	(UI, attach, menu: Help, to: LibraryView)
	(UI, attach, menu: Help, to: QueryView)
	(UI, attach, menu: Help, to: SchemaView)
	(UI, attach, menu: Help, to: ClassView)
	(UI, attach, menu: Help, to: TestRunner)
	(UI, attach, menu: Help, to: BookEdit)
	(UI, attach, menu: Help, to: 'Version Control')
	)
)