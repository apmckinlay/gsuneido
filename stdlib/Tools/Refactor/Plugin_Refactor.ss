// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
#(
Contributions:
	(
	(UI, action, menu: Refactor,
		name: "Rename Variable...", accel: 'Ctrl+Alt+R',
		target: function (@args) { Refactor_Rename_Variable(@args) })
	(UI, action, menu: Refactor,
		name: "Rename Member...",
		target: function (@args) { Refactor_Rename_Member(@args) })
	(UI, action, menu: Refactor,
		name: "Extract Table...",
		target: function (@args) { Refactor_Table_Class(@args) })
	(UI, action, menu: Refactor,
		name: "Extract Method...", accel: 'Ctrl+Alt+M',
		target: function (@args) { Refactor_Extract_Method(@args) })
//	(UI, action, menu: Refactor,
//		name: "Extract Member...",
//		target: function (@args) { Refactor_Extract_Member(@args) })
	(UI, action, menu: Refactor,
		name: "Extract Variable...", accel: 'Ctrl+Alt+L'
		target: function (@args) { Refactor_Extract_Variable(@args) })
	(UI, action, menu: Refactor,
		name: "Inline Variable...", accel: 'Ctrl+Alt+I'
		target: function (@args) { Refactor_Inline_Variable(@args) })
//	(UI, action, menu: Refactor,
//		name: "Extract Function...",
//		target: function (@args) { Refactor_Extract_Function(@args) })
	(UI, action, menu: Refactor,
		name: "Convert Function To Class...",
		target: function (@args) { Refactor_Convert_Function_To_Class(@args) })
	(UI, action, menu: Refactor,
		name: "Format Code",
		target: function (@args) { Refactor_Format_Code(@args) })

	(UI, attach, menu: Refactor, to: LibraryView)
	)
)
