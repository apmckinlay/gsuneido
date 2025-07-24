// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// TODO: handle dragging records in the browse
PassthruController
	{
	curAccels: false
	New(@args)
		{
		super(.controls(@args))
		.flip = .FindControl('flip')
		.browse = .flip.GetControl(.pos.browse)
		if .linkField is false
			.browse.ResetCustomKey(.customKey)
		customFields = .browse.GetCustomFields()
		formWrapper = .FindControl('formWrapper')
		layout = .formLayout
		if Function?(.formLayout)
			layout = (.formLayout)()
		formWrapper.Append(Object('Record', layout, custom: customFields))
		.form = .flip.FindControl('Data')
		}
	Startup()
		{
		if true is .Send('Multiview_Handles_Flip?')
			return
		// make accelerators work inside book
		.curAccels = .Window.SetupAccels(.Commands)
		if .Window.Member?('Ctrl')
			.Window.Ctrl.Redir('On_Flip', this)
		else
			.Defer({ .Window.Ctrl.Redir('On_Flip', this) })
		}
	Data()
		{
		.Send('Data')
		}
	Getter_Commands()
		{
		return Object(#('Flip', 'Alt+A'))
		}
	customKey: false
	controls(.query, form, .columns, validField = false, .protectField = false,
		.linkField = false, .headerFields = #(), .dataMember = 'browseflipform_data',
		name = 'BrowseFlipForm', primary_accessobserver = false, .keyField = false,
		statusBar = false, mandatoryFields = #(), columnsSaveName = '', title = '',
		extraCtrls = #Skip, expandLayout = false, extraFmts = false, customDelete = #(),
		preventCustomExpand? = false, stretchColumn = false)
		{
		.formLayout = form
		if linkField is false
			{
			.customKey = ListCustomize.BuildCustomKeyFromQueryTitle(query, title)
			if columnsSaveName is '' and .customKey isnt false
				columnsSaveName = .customKey
			}
		browse = Object(.linkField isnt false ? 'LineItem' : 'Browse', query, columns,
			:validField, :protectField, :primary_accessobserver, :linkField,
			:headerFields, :dataMember, :statusBar, :mandatoryFields, :columnsSaveName,
			:expandLayout,
			extraMenu: .extraMenu(preventCustomExpand?),
			:extraFmts, :keyField, :customDelete, :preventCustomExpand?,
			switchToForm:,
			:stretchColumn)
		if name isnt 'BrowseFlipForm'
			browse.name = name
		titleCtrl = title isnt ""
			? Object('CenterTitle', title, custom_screen: .customKey isnt false)
			: #(Skip 0)
		return Object('Vert',
			titleCtrl,
			.linkField isnt false
				? ''
				: Object(#Horz #(Button 'Flip', pad: 40, tip: 'Alt+A'), #Fill, extraCtrls)
			Object('Flip',
				browse,
				Object('Scroll',
					#('Vert', name: 'formWrapper'))
					name: 'flip')
			)
		}

	extraMenu(preventCustomExpand?)
		{
		extraMenu = Object()
		if not preventCustomExpand?
			extraMenu.Add('Expand All', 'Contract All')
		extraMenu.Add('Switch to form view\tAlt+A')
		return extraMenu
		}

	pos: (browse: 0, form: 1)
	Flipped?()
		{
		return .flip.GetCurrent() is .pos.form
		}

	On_Flip()
		{
		if .Destroyed?()
			return
		rec = false
		cur = .flip.GetCurrent()
		if cur is .pos.browse
			{
			if .linkField isnt false
				.browse.SelectRecordByFocus(GetFocus())
			if false is (rec = .getSelectedRow()) or
				false is .Send('BrowseFlipForm_Flip_To_Form', rec)
				return
			}
		else if false is .Send('BrowseFlipForm_Flip_To_Browse')
			return

		if .linkField is false and cur is .pos.form and .form.Valid(forceCheck:) isnt true
			{
			.AlertInfo('Flip', 'Please correct the information on the current screen')
			return
			}

		.flip.SetCurrent(cur = cur is .pos.browse ? .pos.form : .pos.browse)
		if cur is .pos.form
			.flip_to_form(rec, .browse.GetReadOnly())
		else if .linkField isnt false
			SetFocus(.browse.GetGridHwnd())
		}
	getSelectedRow(alertTitle = 'Flip')
		{
		if false is rec = .getCurrentRecord()
			{
			.AlertInfo(alertTitle, 'Please select a line')
			return false
			}
		if rec.listrow_deleted is true or rec.vl_deleted is true
			{
			.AlertInfo(alertTitle, 'Line is marked as deleted')
			return false
			}
		return rec
		}

	getCurrentRecord()
		{
		return .linkField isnt false
			? .browse.GetSelectedRecord()
			: .browse.GetCurrentRecord()
		}

	flip_to_form(rec, readonly)
		{
		.form.Set(rec)
		.form.SetProtectField(.protectField)
		.form.SetReadOnly(readonly)
		.FocusFirst(.FindControl('formWrapper').Parent.Hwnd)
		}

	On_Context_Expand_All()
		{
		.browse.ExpandByField(.getExpandableItems(), .keyField)
		}

	On_Context_Contract_All()
		{
		.browse.ExpandByField(.getExpandableItems(), .keyField, collapse?:)
		}

	getExpandableItems()
		{
		.browse.GetAllLineItems(includeAll?:).
			Filter({ it.vl_deleted isnt true }).
			Map({ it[.keyField] })
		}

	On_Context_Switch_to_form_view()
		{
		.On_Flip()
		}

	Browse_BeforeValid()
		{
		.form.HandleFocus() // to commit current editing before valid and saving
		}
	Browse_ExtraValid()
		{
		return .extraValid("Browse_ExtraValid")
		}
	LineItem_ExtraValid()
		{
		return .extraValid("LineItem_ExtraValid")
		}
	extraValid(msg)
		{
		if false is .Send(msg) // resend for Access
			return false
		if .flip.GetCurrent() is .pos.browse
			return true
		return .form.Valid() is true
		}
	Browse_AfterSet()
		{
		if .flip.GetCurrent() is .pos.form
			{
			if .keyField isnt false // reset record in form
				{
				form_rec = .form.Get()
				browsedata = .getAllData()
				if false isnt (i = browsedata.FindIf(
					{ |x| x[.keyField] is form_rec[.keyField]}))
					{
					rec = browsedata[i]
					.flip_to_form(rec, .browse.GetReadOnly())
					return
					}
				}
			.On_Flip() // have to flip back if record not found
			}
		}

	LineItem_AfterSet()
		{
		.Browse_AfterSet()
		.setCornerButtonColor(gray?:)
		}

	LineItem_ItemSelected(rec)
		{
		if .linkField is false
			return 0
		.setCornerButtonColor(rec.vl_deleted is true)
		return 0
		}

	LineItem_DeleteRecord(record, source = false)
		{
		.setCornerButtonColor(record.New?() or record.vl_deleted is true)
		.Send('LineItem_DeleteRecord', record, :source)
		}

	setCornerButtonColor(gray? = false)
		{
		if false is hdrCornerCtrl = .browse.GetHdrCornerControl()
			return
		if gray?
			hdrCornerCtrl.SetImageColor(CLR.Inactive, CLR.Inactive)
		else
			hdrCornerCtrl.SetImageColor(CLR.Highlight, CLR.Highlight)
		}

	LineItem_SwitchToForm()
		{
		.On_Flip()
		}

	getAllData()
		{
		return .linkField isnt false
			? .browse.GetAllLineItems()
			: .browse.GetAllBrowseData()
		}

	Browse_AllowDelete(record)
		{
		return .Send("BrowseFlipForm_AllowDelete", record)
		}
	List_AllowMove()
		{
		if (.flip.GetCurrent() is .pos.form or
			not .browse.GetList().GetHighlighted().Empty?() or
			false is .Send("BrowseFlipForm_AllowMove"))
			return false
		return true
		}
	GetRecord()
		{
		return .flip.GetCurrent() is .pos.form
			? .form.Get()
			: .getCurrentRecord()
		}
	GetSelectedRecord(alertTitle = 'Flip')
		{
		return .flip.GetCurrent() is .pos.form
			? .form.Get()
			: .getSelectedRow(alertTitle)
		}
	GetDeleted()
		{
		if .linkField isnt false
			return .browse.GetDeleted()
		return .browse.GetAllBrowseData().Difference(.browse.GetBrowseData())
		}
	GetTransQuery()
		{
		if .linkField isnt false
			return .query
		return .browse.GetTransQuery()
		}
	SetMainRecordField(field, value) // so custom fill-in works
		{
		.form.SetField(field, value)
		}

	// needed for custom fields
	GetQuery()
		{
		return .query
		}

	GetAccessCustomKey()
		{
		if .linkField isnt false
			return .Send('GetAccessCustomKey')
		return .customKey
		}

	Default(@args)
		{
		return .browse[args[0]](@+1 args)
		}

	Valid?()
		{
		cur = .flip.GetCurrent()
		if cur is .pos.form and .form.Valid() isnt true
			return false

		return .browse.Valid?()
		}

	readonly: false
	GetReadOnly()
		{
		return .readonly
		}

	SetReadOnly(readOnly)
		{
		.readonly = readOnly
		.browse.SetReadOnly(readOnly)
		.form.SetReadOnly(readOnly)
		}

	NoData()
		{
		.Send('NoData')
		}

	Destroy()
		{
		if .curAccels isnt false
			.Window.RestoreAccels(.curAccels)
		super.Destroy()
		}
	}
