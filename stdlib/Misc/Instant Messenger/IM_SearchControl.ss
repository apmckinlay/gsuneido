// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Search'

	New()
		{
		super(.layout())
		}

	layout()
		{
		return Object('Vert'
				Object('VirtualList', 'im_history where im_from is ' $
					Display(Suneido.User) $ '
					rename imchannel_num to imchannel_num_search',
					name: 'listSearch',
					columns:
						#(im_num,
						im_from,
						im_to,
						im_message,
						imchannel_num_search),
					columnsSaveName: .Title,
					readonly:,
					filtersOnTop:,
					title: .Title,
					resetColumns:,
					enableMultiSelect:
					),
			)
		}

	VirtualList_ExtraSetupRecordFn()
		{
		return .BeforeRecord
		}

	BeforeRecord(rec)
		{
		rec.im_num = rec.im_num.FormatEn('dd/MM/yyyy')
		}

	}
