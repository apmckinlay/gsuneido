// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(recordData, minimizeOutput?, missingTest ='check_local')
		{
		missingTest? = .missingTest?(recordData, missingTest)
		desc = .getDescription(missingTest?, minimizeOutput?)
		maxRating = 5
		if missingTest?
			maxRating = 4
		return Object(warnings: Object(), :desc, :maxRating, size: missingTest? ? 1 : 0)
		}

	getDescription(missingTest?, minimizeOutput?)
		{
		if missingTest?
			return 'A test class was not found for this record. Please create one'

		if not minimizeOutput?
			return 'A test class was found for this record or is not required'
		return ''
		}

	missingTest?(recordData, missingTest)
		{
		if recordData.recordName.Suffix?("Test") or
			.testBaseClass?(recordData.recordName, recordData.code) or
			recordData.recordName.Prefix?("Field_") or
			LibRecordType(recordData.code) not in ('function', 'class')
			return false

		name = recordData.recordName.RightTrim('?')
		if Boolean?(missingTest)
			return missingTest
		else
			return .localMissingTest?(recordData.lib, name)
		}

	testBaseClass?(name, code)
		{
		if not name.Suffix?('Tests')
			return false
		scan = ScannerWithContext(code) // skips comments and whitespace
		token = scan.Next()
		return token is 'Test' or token.Suffix?('Tests')
		}

	localMissingTest?(lib, name)
		{
		return QueryEmpty?(lib $
			' where group is -1 and name in ("' $ name $ '_Test","' $ name $ 'Test")')
		}
	}