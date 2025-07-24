// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (view, titleLeftCtrl)
	{
	args = Object()
	pre = 'VirtualListViewControl_'
	view.Members().Filter({ it.Prefix?(pre) and not Function?(view[it]) }).
		Each({ args[it.RemovePrefix(pre)] = view[it] })
	layout = Object('Vert',
		args.title isnt false
			? Object('CenterTitle', args.title, :titleLeftCtrl,
				search: args.filtersOnTop)
			: '')
	layout.Append(view.Addons.Collect('ExtraLayout'))
	layout.Add(Object('VertSplit',
		#('Vert', name: 'select', ystretch: 1),
		Object('Vert'
		Object('VirtualListScroll',
			Object('Horz',
				Object("VirtualListExpandBar", args.preventCustomExpand?,
					args.enableDeleteBar, switchToForm: args.switchToForm)
				Object(#Vert
					Object('VirtualListHeader',
						headerSelectPrompt: args.headerSelectPrompt,
						checkBoxColumn: args.checkBoxColumn)
					Object('VirtualListGrid', protectField: args.protectField))
				),
			Object('VirtualListThumbBar',
				disableSelectFilter: args.disableSelectFilter or args.filtersOnTop),
			'VirtualListSelectButton', args.thinBorder, args.hdrCornerCtrl,
			expandExtra: Object('VirtualListExpandButtons', args.switchToForm)),
		args.filtersOnTop isnt false or args.validField isnt false ? 'Status' : '',
		ystretch: 2),
		splitter: #Vert /* no splitter*/))
	return layout
	}
