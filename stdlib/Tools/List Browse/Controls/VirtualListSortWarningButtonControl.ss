// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'VirtualListSortWarningButton'
	New(.view)
		{
		super(.layout())
		.horz = .FindControl('sortWarning')
		}

	layout()
		{
		Assert(.view.GetModel().CheckAboveSortLimit?())
		return Object('Horz', name: 'sortWarning')
		}

	InsertWarning()
		{
		if .horz.GetChildren().Empty?()
			.horz.Append(Object('Horz',
				#(Skip, small:)
				Object('EnhancedButton', image: 'triangle-warning',
					command: 'FixDisabledSort', imageColor: CLR.orange,
					mouseOverImageColor: CLR.EnhancedButtonFace, alignTop:,
					tip: 'Sorting is disabled with your current Select. ' $
						'Click for more details'
					)
				#(Skip, small:),
				name: 'sortWarning'
				))
		}

	RemoveWarning()
		{
		.horz.RemoveAll()
		}

	On_FixDisabledSort()
		{
		.view.Addons.Send('FixDisabledSort')
		}
	}
