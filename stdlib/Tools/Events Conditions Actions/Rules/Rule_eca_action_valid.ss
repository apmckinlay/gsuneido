// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if not Object?(.eca_actions) or .eca_actions.Empty?() or
		.eca_actions.Any?({ it.GetDefault('action_name', '') is '' })
		return 'Please select at least one action.'

	for a in .eca_actions
		Plugins().ForeachContribution('ECA', 'action')
			{ |c|
			if c.Member?('validFunc') and c.name is a.action_name and
				'' isnt result = (c.validFunc)(a)
				return result
			}
	return ''
	}