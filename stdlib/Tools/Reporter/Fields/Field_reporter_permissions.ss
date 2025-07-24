// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
Field_string
	{
	Prompt: 'Available to'
	Control: (ChooseMany, listField: 'reporter_permission_list', saveNone:,
		text: 'If no groups are selected, only the user who created the report\n' $
			'will have access to it. Groups must also have permission to\n' $
			'the data source in order to use reports.')
	}