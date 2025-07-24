// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Manage Channels"

	New()
		{
		super(.layout())
		.Window.HelpPage = "/res/IMManageChannels"
		}

	layout()
		{
		return Object('VirtualList',
			'im_channels rename imchannel_num to imchannel_num_default,
				imchannel_status to imchannel_status_default ',
			name: 'ListChannel',
			columns:
				#(imchannel_name,
				imchannel_abbrev,
				imchannel_desc,
				imchannel_status_default)
			defaultColumns:
				#(imchannel_name,
				imchannel_abbrev,
				imchannel_desc,
				imchannel_status_default)
			mandatoryFields: #(imchannel_name, imchannel_abbrev),
			validField: 'imchannel_valid',
			protectField: 'imchannel_protect',
			columnsSaveName: .Title,
			filtersOnTop:,
			title: .Title,
			select: #(('imchannel_status_default', '=', "active")),
			resetColumns:)
		}
	}






















