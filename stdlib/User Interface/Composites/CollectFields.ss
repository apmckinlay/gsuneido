// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(layouts, path? = false)
		{
		inst = new this()
		return inst.Collect(layouts, path?)
		}

	New()
		{
		.dynamicCtrls = GetContributions(#DynamicControls)
		}

	Collect(layouts, path?)
		{
		path = path? ? Object() : false
		collect = Object()
		layouts.Each()
			{
			if Object?(it) // Can be a control class reference
				.collectFields(it, collect, path)
			}
		return collect.UniqueValues()
		}

	collectFields(ob, collect, path = false)
		{
		if path isnt false
			{
			path = path.DeepCopy()
			first = ob.GetDefault(0, '')
			if ob.GetDefault('Tab', '') isnt ''
				{
				if .skipTab?(ob)
					return
				path.Add(Object(section: ob.Tab, type: 'Tab'))
				collect.Add(path.DeepCopy())
				}
			if first is 'Accordion'
				path.Add(Object(section: '', type: 'Accordion'))
			if first is 'Letterheads'
				return
			}
		if true isnt ob.GetDefault(#Hide?, false)
			for member, ctrl in ob
				.checkCtrl(member, ctrl, collect, path is false ? path : path.DeepCopy())
		}

	skipTab?(ctrl)
		{
		tab = ctrl.Tab
		if tab is 'Hidden Options' and .user() not in ('default', 'axon')
			return true
		if tab is 'Developer Options' and .user() isnt 'default'
			return true
		if true is ctrl.GetDefault(#Hide?, false)
			return true
		return false
		}

	user()
		{
		return Suneido.User
		}

	checkCtrl(member, ctrl, collect, path)
		{
		.collectControls(ctrl, path, collect)
		if String?(member) and not member.Blank?() // named argument for a control
			return
		if String?(ctrl)
			.stringDef(ctrl, collect, path)
		else if Object?(ctrl) and ctrl.NotEmpty?()
			.objectDef(ctrl, collect, path)
		}

	collectControls(ctrl, path, collect)
		{
		if path isnt false and Object?(ctrl)
			{
			first = ctrl.GetDefault(0, '')
			second = ctrl.GetDefault(1, '')
			.collectAccordion(first, path)

			if first in ('Static', 'Heading', 'Heading1', 'Heading2', 'Heading3') and
				collect.Size() > 0
				{
				item = path.DeepCopy()
				item.Add(Object(section: second, type: 'Static', name: second))
				collect.Add(item)
				}
			}
		}

	collectAccordion(first, path)
		{
		if String?(first) and path.Size() > 0 and
			path.Last().type is 'Accordion' and path.Last().section is ''
			{
			path.Last().section = first
			}
		}

	stringDef(ctrlStr, collect, path)
		{
		if not .dynamicControl?(Object(ctrlStr), collect, path)
			.fieldDef(ctrlStr, collect, path)
		}

	fieldDef(ctrlStr, collect, path)
		{
		if Prompt(ctrlStr) is ctrlStr
			return

		dd = Datadict(ctrlStr)
		if path is false
			collect.Add(ctrlStr)
		else
			{
			item = path.DeepCopy()
			item.Add(Object(section: Prompt(ctrlStr), type: 'Field', name: ctrlStr))
			collect.Add(item)

			if dd.Base?(Field_config_cols) or
				GetControlClass.FromField(ctrlStr).Base?(ChooseTwoListAndResetControl)
				{
				item = path.DeepCopy()
				item.Add(Object(section: Prompt(ctrlStr) $ ' > Reset',
					type: 'ResetButton', name: ctrlStr))
				collect.Add(item)
				}
			}
		ctrlOb = dd.Control.Copy()
		ctrlOb.fieldName = ctrlStr
		.dynamicControl?(ctrlOb, collect, path)
		}

	objectDef(ctrlOb, collect, path)
		{
		if ctrlOb.Member?(#name)
			.stringDef(ctrlOb.name, collect, path)
		if ctrlOb.GetDefault(0, '') in ('Button', 'EnhancedButton') and path isnt false
			{
			item = path.DeepCopy()
			btnText = ctrlOb.GetDefault(1, '')
			item.Add(Object(section: btnText, type: 'Button',
				name: ctrlOb.GetDefault('name', ToIdentifier(btnText))))
			collect.Add(item)
			}
		if not .dynamicControl?(ctrlOb, collect, path)
			.collectFields(ctrlOb, collect, path)
		}

	dynamicControl?(ctrlOb, collect, path)
		{
		// false indicates a non control Object (IE: control arguments)
		ctrl = ctrlOb.GetDefault(0, false)
		if not String?(ctrl) or not ctrl.Capitalized?() or
			false isnt ctrl.Match('[[:space:]]')
			return false
		if dynamic? = false isnt ctrl = .dynamicCtrl(ctrl)
			.collectFields((.dynamicCtrls[ctrl])(@.dynamicArgs(ctrlOb)), collect, path)
		return dynamic?
		}

	dynamicCtrl(ctrl)
		{
		if .dynamicCtrls.Member?(ctrl)
			return ctrl
		if not ctrl.Suffix?(#Control)
			return false
		return .dynamicCtrls.Member?(ctrl = ctrl.BeforeLast(#Control))
			? ctrl
			: false
		}

	dynamicArgs(ctrlOb)
		{
		args = ctrlOb.Copy().Delete(0)
		if not args.Member?('fieldName')
			args.fieldName = args.GetDefault('name', '')
		return args
		}

	FindTab(field, collectedFields) // Relies on path?:
		{
		for ob in collectedFields
			{
			if ob.Size() is 1
				{
				if ob[0].type is 'Field' and ob[0].GetDefault('name', '') is field
					return false // Field is on the header, not a tab
				}
			else if ob[1].GetDefault('name', '') is field
				return ob[0].section
			}
		return false
		}
	}
