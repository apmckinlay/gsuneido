// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function(view, titleLeftCtrl)
	{
	args = Object()
	pre = #VirtualListViewControl_
	view.
		Members().
		Filter({ it.Prefix?(pre) and not Function?(view[it]) }).
		Each({ args[it.RemovePrefix(pre)] = view[it] })
	layout = [#Vert,
		args.title isnt false
			? [#CenterTitle, args.title, :titleLeftCtrl,
				search: args.filtersOnTop, subTableSelect: args.subTableSelect]
			: ""]
	layout.Append(view.Addons.Collect(#ExtraLayout))
	layout.Add([#VertSplit,
		#(Vert, name: select, ystretch: 1),
		[#Vert,
			[#VirtualListScroll,
				[#Horz,
					[#VirtualListExpandBar, args.preventCustomExpand?,
						args.enableDeleteBar, switchToForm: args.switchToForm],
					[#Vert,
						[#VirtualListHeader,
							headerSelectPrompt: args.headerSelectPrompt,
							checkBoxColumn: args.checkBoxColumn],
						[#VirtualListGrid, protectField: args.protectField]]],
				[#VirtualListThumbBar,
					disableSelectFilter: args.disableSelectFilter or args.filtersOnTop],
				#VirtualListSelectButton, args.thinBorder, args.hdrCornerCtrl,
				expandExtra: [#VirtualListExpandButtons, args.switchToForm]],
			args.filtersOnTop isnt false or args.validField isnt false ? #Status : "",
			ystretch: 2],
		splitter: #Vert /* no splitter*/])
	return layout
	}
