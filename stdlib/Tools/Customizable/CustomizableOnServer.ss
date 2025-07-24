// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	NotifyCustomizableFieldsChanges(key)
		{
		.Synchronized()
			{
			if not Suneido.Member?('CustomizableFieldsChanges')
				Suneido.CustomizableFieldsChanges = Object().Set_default(Object())

			changes = Suneido.CustomizableFieldsChanges.custom_fields
			if changes.NotEmpty?() and changes.Last().key is key
				changes.Last().asof = Date()
			else
				changes.Add(Object(:key, asof: Date()))
			}
		}

	GetChangesOnServer()
		{
		.Synchronized()
			{
			if not Suneido.Member?('CustomizableFieldsChanges')
				return #()

			return .getAndRemoveExpiredChanges(Suneido.CustomizableFieldsChanges)
			}
		}

	getAndRemoveExpiredChanges(allChanges)
		{
		return Object(custom_fields: .removeExpired(allChanges.custom_fields, 'key')
			custom_screen: .removeExpired(allChanges.custom_screen, 'field'))
		}

	serverTimeout: 300 // 5 minutes
	removeExpired(changes, keyField)
		{
		date = .now()
		return changes.
			RemoveIf( { date.MinusSeconds(it.asof) >= .serverTimeout } ).
			Map({ it[keyField] }).
			UniqueValues()
		}

	now()
		{
		return Date()
		}

	NotifyCustomScreenChanges(field)
		{
		.Synchronized()
			{
			if not Suneido.Member?('CustomizableFieldsChanges')
				Suneido.CustomizableFieldsChanges = Object().Set_default(Object())

			Suneido.CustomizableFieldsChanges.custom_screen.Add(
				Object(:field, asof: Date()))
			}
		}
	}