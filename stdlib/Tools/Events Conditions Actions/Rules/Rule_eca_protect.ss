// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	if '' isnt msg = EventConditionActions.PermissionToChange()
		return msg

	if Object?(.eca_conditions) and not .eca_conditions.Empty?() and
		.eca_conditions.Any?({ it.HasNonEmptyMember?(it.Members()) })
		return #(eca_event:, allowDelete: )

	if Object?(.eca_actions) and not .eca_actions.Empty?() and
		.eca_actions.Any?({ it.GetDefault('action_name', '') isnt '' })
		return #(eca_event:, allowDelete: )

	if .eca_event is ''
		return #(eca_actions:, eca_conditions:, allowDelete: )

	return false
	}