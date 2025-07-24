// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ReasonProtected(access, protectField, foreignKeyUsage)
		{
		reasonDeleteNotAllowed = .reasonDeleteNotAllowed(access)
		reasonDeleteNotAllowed $= foreignKeyUsage
		protect = access.GetData()[protectField]
		if .noProtectReason(protect)
			protect = ''
		if Object?(protect)
			protect = protect.reason
		if protect isnt '' or reasonDeleteNotAllowed isnt ''
			{
			msg = protect $ (reasonDeleteNotAllowed isnt ''
				? '\n\n' $ reasonDeleteNotAllowed
				: '')
			access.AlertInfo('Reason Protected', msg.Trim())
			}
		else
			access.AlertInfo('Reason Protected', 'No Information')
		}

	noProtectReason(protect)
		{
		return (protect is '' or protect is false or protect is true or
			(Object?(protect) and not protect.Member?('reason')) or
			(Object?(protect) and protect.Member?('reason') and protect.reason is ""))
		}

	AllowDelete?(access, record, protect, protectField, newrecord?)
		{
		msg = ''
		if (protect or not ProtectRuleAllowsDelete?(record, protectField, newrecord?))
			msg = .buildProtectedReason(access, protectField)

		if '' isnt reasonDeleteNotAllowed = .reasonDeleteNotAllowed(access)
			msg $= '\n\n' $ reasonDeleteNotAllowed

		if msg isnt ''
			{
			if not msg.Prefix?('This record can not be deleted.')
				msg = 'This record can not be deleted.' $ msg
			access.AlertInfo('Delete', msg)
			return false
			}
		return true
		}

	buildProtectedReason(access, protectField)
		{
		msg  = 'This record can not be deleted.'
		if protectField isnt false
			{
			protect = .protectResult(access, protectField)
			reason = Object?(protect) and protect.Member?('reason') and
				protect.reason isnt ''
				? protect.reason
				: (String?(protect) ? protect : '')
			if reason isnt ''
				msg $= '\n\n' $ reason
			}
		return msg
		}

	protectResult(access, protectField)
		{
		return access.GetRecordControl().GetField(protectField)
		}

	reasonDeleteNotAllowed(access)
		{
		RecordAllowDelete(access.GetQuery(), access.GetData())
		}

	AllowEdit?(data, protect, protectField)
		{
		return protect is false and
			(protectField is false or
			data.GetField(protectField) is false or
			String?(data.GetField(protectField)) or
			Object?(data.GetField(protectField)))
		}
	}