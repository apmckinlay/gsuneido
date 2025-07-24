// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(record, protect_rule, new_record = false, header_delete = false)
		{
		if protect_rule is false
			return true

		protect = .EvalProtectRule(record, protect_rule)

		if .isheaderDeleteAndIsAllowed?(header_delete, protect)
			return true

		// need to check for neverDelete even for new records
		// i.e. for system generated new records (in browses)
		// that you don't want the user to be able to delete
		if .neverDeleteSpecified?(protect)
			return false

		return new_record is true or .protectRuleShouldNotPreventDelete?(protect)
		}

	// This method exists only for the purpose of overriding in the test. The only
	// reason it is public is that it seems you cannot access the private name
	// externally when the class name ends in ?
	EvalProtectRule(record, protect_rule)
		{
		return record.Eval(Global('Rule_' $ protect_rule))
		}

	isheaderDeleteAndIsAllowed?(header_delete, protect)
		{
		if not header_delete
			return false

		return Object?(protect) and protect.GetDefault('allowHeaderDelete', false) is true
		}

	neverDeleteSpecified?(protect)
		{
		return Object?(protect) and protect.GetDefault("neverDelete", false) is true
		}

	protectRuleShouldNotPreventDelete?(protect)
		{
		if protect is false or protect is ""
			return true

		return Object?(protect) and protect.GetDefault("allowDelete", false) is true
		}
	}