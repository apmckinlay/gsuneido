// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
/* READ ME:
	This control is designed to be a flexible code editor.
	This asks the various CodeViewAddons provided by the editor and constructs the
	required containers in order for them to function. Additionally, it automatically
	redirects any methods the addons require to them.
	All CodeViewAddons should:
		1. Specify the container its controls will be added to (if any)
		2. Specify any methods which should be redirected to it

	If both are defined properly, the controls it defines should be added automatically.
	Additionally, any methods defined should automatically reach the addon, allowing for
	complete encapsulation of the addon code.

	By default, this control has a basic set of addons it will use for the editor.
	If you want to add / remove / replace, define the addons argument as such:
		addons: #(
			// The below will add this addon to the editor
			Addon_name1: true,
			// The below will remove this addon from the editor (if it was present)
			// 	IE: Debugger specifies: #(Addon_status: false, Addon_overview_bar: false)
			Addon_name2: false,
			// The below will replace the specified addon with the new one
			// 	IE: LibViewViewControl replaces Addon_status with Addon_editor_status
			Addon_name3: #(addon: Addon_name4, [OPTIONAL args: ...])
			// The below will overwrite the addon arguments.
			// This will add the addon if it is not already specified.
			Addon_name5: #(<args>)
			)
*/
ExplorerAdapterControl
	{
	Name: 	CodeView
	BaseAddons: #(
		Addon_suneido_style:,
		Addon_indent_guides:,
		Addon_brace_match:,
		Addon_calltips:,
		Addon_class_outline:,
		Addon_folding:,
		Addon_go_to_definition:,
		Addon_go_to_line:,
		Addon_highlight_cursor_line:,
		Addon_highlight_occurrences:,
		Addon_show_line_numbers:,
		Addon_show_margin:,
		Addon_status:,
		Addon_overview_bar:,
		Addon_flag:,
		Addon_scroll_zoom:
		)
	redirMethods: #()
	New(.data = #(table: '', name: '', text: ''), addons = #()
		divider = #HorzSplit, xstretch = 3, ystretch = 1, ide = true, .readonly = false)
		{
		super(.controls(divider, addons, ide, xstretch, ystretch), #text)
		}

	Exterior:			codeViewExterior
	TopLeft:			codeViewTopLeft
	TopRight:			codeViewTopRight
	BottomLeft:			codeViewBotLeft
	BottomRight:		codeViewBotRight
	Getter_InjectPoints()
		{
		return .InjectPoints = [
			editor: 	.TopLeft,
			exterior: 	.Exterior,
			topLeft: 	.TopLeft,
			topRight: 	.TopRight,
			bottomLeft: .BottomLeft
			bottomRight: .BottomRight
			]
		}

	/*vvvvvvvvvv Control Building Start vvvvvvvvv*/
	controls(divider, addons, ide, xstretch, ystretch)
		{
		_xstretch = xstretch
		_ystretch = ystretch
		.subSplitA = .Name $ #_SubSplitA
		.subSplitB = .Name $ #_SubSplitB
		.primarySplit = .Name $ #_DividerSplit
		.editor = .buildEditor(addons = .editorAddons(addons), ide)
		_injectAddons = .injectAddons(addons)
		interior = .buildControls(divider)
		codeView = [#Vert, #(Skip, small:), interior, name: .Exterior]
		return .addAddonControls(#exterior, codeView)
		}

	editorAddons(addons)
		{
		editorAddons = .BaseAddons.Copy()
		for addon, state in addons
			{
			if state is true
				editorAddons[addon] = state
			else if state is false
				editorAddons.Delete(addon)
			else if Object?(state)
				.overwriteAddon(addon, state, editorAddons)
			}
		return editorAddons
		}

	overwriteAddon(addon, state, editorAddons)
		{
		if state.Member?(#addon)
			{
			editorAddons.Delete(addon)
			editorAddons[state.addon] = state.GetDefault(#args, true)
			}
		else
			editorAddons[addon] = state
		}

	injectAddons(editorAddons)
		{
		injectAddons = Object()
		editorAddons.Members().Each()
			{
			addon = Global(it)
			if addon.Member?(#Inject)
				injectAddons.Add(addon)
			}
		return injectAddons
		}

	buildEditor(addons, ide)
		{ return Object(ScintillaAddonsControl, IDE: ide).MergeNew(addons) }

	buildControls(divider)
		{
		pairs = .buildPairs(divider,
			.buildGroup(#topLeft), .buildGroup(#topRight),
			.buildGroup(#bottomLeft), .buildGroup(#bottomRight))
		pairA = pairs.pairA
		pairB = pairs.pairB
		interior = pairA.NotEmpty?() and pairB.NotEmpty?()
			? [divider, pairA, pairB, name: .primarySplit]
			: pairA.NotEmpty?()
				? pairA
				: pairB
		return interior
		}

	buildGroup(injectPoint)
		{
		controlName = .InjectPoints[injectPoint]
		baseControl = controlName is .InjectPoints.editor
			? .editorControl(controlName)
			: [#Vert, name: controlName]
		return .addAddonControls(injectPoint, baseControl)
		}

	addAddonControls(injectPoint, baseControl, _injectAddons)
		{
		addonControls = Object().Set_default('')
		addons = injectAddons.Filter({ it.Inject is injectPoint })
		addons.Each({ it.Controls(addonControls) })
		if .skipGroup?(baseControl, addonControls, addons)
			return false
		addonControls.Members().Sort!().Each({ baseControl.Add(addonControls[it]) })
		return baseControl
		}

	editorControl(name, _xstretch, _ystretch)
		{ return [#Horz, .editor, :xstretch, :ystretch, :name, editorGroup?:] }

	skipGroup?(baseControl, addonControls, addons)
		{
		// Always add the editor/exterior groups, even if empty
		if .RequiredGroups().Has?(baseControl.name)
			return false
		return addonControls.Empty?() and not addons.Any?({ it.Method?(#InjectControls) })
		}

	RequiredGroups()
		{ return Object(.InjectPoints.editor, .InjectPoints.exterior) }

	buildPairs(divider, topLeft, topRight, botLeft, botRight)
		{
		if divider is #HorzSplit
			{
			pairA = [topLeft, botLeft]
			pairB = [topRight, botRight]
			subDivider = #VertSplit
			}
		else
			{
			pairA = [topLeft, topRight]
			pairB = [botLeft, botRight]
			subDivider = #HorzSplit
			}
		pairA = .pairLayout(subDivider, pairA.Filter({ it isnt false }), .subSplitA)
		pairB = .pairLayout(subDivider, pairB.Filter({ it isnt false }), .subSplitB)
		return [:pairA, :pairB]
		}

	pairLayout(subDivider, pair, name)
		{
		if pair.Empty?()
			return #()
		editorPair? = pair.Any?({ it.GetDefault(#editorGroup?, false) })
		control = pair.Size() is 2
			? [subDivider, pair.PopFirst(), pair.PopFirst(), :name]
			: [#Vert, pair.PopFirst()]
		control.editorPair? = editorPair?
		return control
		}
	/*^^^^^^^^^^^^ Control Building Start ^^^^^^^^^^^*/

	/*vvvvv Editor Addon Helper Methods Start vvvvv*/
	addonRedirs()
		{ return .Editor.CollectFromAddons(#Addon_RedirMethods).Flatten().UniqueValues() }

	Default(@args)
		{
		if .redirMethods.Has?(args[0])
			.send(@args)
		else
			SuneidoLog('ERROR: (CAUGHT) CodeView, undefined method: ' $ args[0],
				params: args[1 ..], caughtMsg: 'development error')
		}

	Msg(args)
		{
		if .redirMethods.Has?(args[0])
			.send(@args)
		return super.Msg(args)
		}

	send(@args)
		{
		if .Editor isnt false
			.Editor.ConditionalSendToAddons({ .addonReady?(it) }, args)
		return 0
		}

	addonReady?(addon)
		{ return addon.Method?(#AddonReady?) ? addon.AddonReady?() : true }
	/*^^^^^ Editor Addon Helper Methods End ^^^^^*/

	Startup()
		{
		.InitialSet(.data)
		.SetReadOnly(.readonly)
		.Editor.SendToAddons(#InjectAddonControls)
		.redirMethods = .addonRedirs()
		}

	InitialSet(data)
		{
		.Table = data.GetDefault(#table, '')
		.RecName = data.GetDefault(#name, '')
		.Set(data)
		.send(#InitialSet)
		}

	Set(data)
		{
		.Table = data.GetDefault(#table, .Table)
		.RecName = data.GetDefault(#name, .RecName)

		// Only carry out the Set if the text has changed.
		// This will prevent the undo/redo queue from being unnecessarily cleared
		if data.Member?(#text) and .Get().text isnt data.text
			super.Set([text: data.text])
		else
			.EN_CHANGE()
		}

	GetSplit(split = #subSplitA)
		{ return .FindControl(this[.Name $ #Control_ $ split]) }

	Getter_Editor()
		{ return .Editor = .FindControl(#Editor) }

	GetChild()
		{ return .Editor }

	CurrentPosition()
		{ return .Editor.GetCurrentPos() }

	CurrentLine()
		{ return .Editor.LineFromPosition(.CurrentPosition()) }

	GetFirstVisibleLine()
		{ return .Editor.GetFirstVisibleLine() }

	SetFirstVisibleLine(pos)
		{ .Editor.SetFirstVisibleLine(pos) }

	CurrentTable()
		{ return .Table }

	CurrentName()
		{ return .RecName }

	Valid?()
		{ return .Editor.Valid?() }

	SetReadOnly(.readonly)
		{ .Editor.SetReadOnly(.readonly) }

	EN_CHANGE()
		{
		if not .Editor.GetReadOnly()
			.Delay(750 /*= ms */, .editorModified, uniqueID: #editor_modified)
		}

	editorModified()
		{ .Send(#SaveCode_AfterChange) }

	AfterSave()
		{ .send(#AfterSave) }

	Invalidate()
		{ .send(#Invalidate) }
	}
