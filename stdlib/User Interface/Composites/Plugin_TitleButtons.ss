// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
#(
ExtensionPoints:
	(
	('button')
	)
Contributions:
	(
	(TitleButtons, 'button', control: FilterButtonControl,
		condition: function (args)
			{
			return args.GetDefault('search', false) is true
			}
		seq: 10)

	(TitleButtons, 'button', control: CustomizeButtonControl,
		condition: function (args) { return args.GetDefault('custom_screen', false) },
		seq: 21)

	(TitleButtons, 'button', control: HelpButtonControl,
		postCondition: function (ctrl)
			{ return ctrl.Window.Base?(Dialog) or ctrl.Window.Base?(ModalWindow) },
		seq: 22)

	(TitleButtons, 'button', control: NotesControl
		seq: 23)

	(TitleButtons, 'button', control: SearchConfigControl,
		condition: function (args) {
			return args.GetDefault(0, "").Has?('Setup Options') },
		seq: 20.5)
	)
)
