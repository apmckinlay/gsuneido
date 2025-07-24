// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
/* USAGE
CHILD CLASSES REQUIRE:
- Name: This should be unique to the Addon and is used to find AddonControl.

- Inject: This is used with both the Controls and InjectControls methods.
	There are six options: [top|bottom]Left, [top|bottom]Right, exterior, editor.
	The majority are self explanatory, the special cases are: exterior and editor.
		exterior: This injection point refers to the outer most container.
		editor: This injection point should match one of the other inject points.
			This is where the Editor can be found and is the group which has display
			priority.

- Controls OR InjectControls:
-- Controls(ob): Define this method if you want the addon controls to be built during
	the standard construction.
-- InjectControls(container): Define this method if you want the addon controls
	to be added at the Controller's discretion (Via: InjectAddonControls).
	This is useful for addons which depend on other controls being constructed first.
	IE: Addon_overview_bar depends on ScintillaAddonsControl

OPTIONAL:
- Addon_RedirMethods: Returns an object of public methods.
	These methods are automatically forwarded to the addons from CodeViewControl.

CONTROLLER REQUIRES:
- InjectPoints(): This method returns the mapping required for adding the controls
	from the addons.
	It should be formatted:
		Template: 	[<see above "Inject" for potential member names>: <control name>]
		Example:	[topLeft: 'UpperVert', editor: 'HtmlEditorVert']
	If a Controller does not provide the inject point for an addon, an error is logged
*/
ScintillaAddon
	{
	Inject: false
	AddonReady?()
		{ return .AddonControl isnt false }

	InjectAddonControls()
		{
		if not .Method?(#InjectControls)
			return

		if false is injectControl = .injectControl()
			return .logError('unable to find inject control for: ' $ .injectPoint)

		.InjectControls(injectControl)
		}

	Controls(@unused)
		{ }

	getter_injectPoint()
		{
		try
			injectPoints = .Controller.InjectPoints
		catch
			return .logError('unable to retrieve InjectPoints')
		if not injectPoints.Member?(.Inject)
			return .logError('controller does not have injection point for: ' $ .Inject)
		return .injectPoint = injectPoints[.Inject]
		}

	injectControl()
		{ return .Controller.FindControl(.injectPoint) }

	Getter_Controller()
		{ return .Parent.Controller }

	logError(e)
		{
		SuneidoLog('ERROR: (CAUGHT) CodeViewAddon ' $ e, params: [addon: .Name]
			caughtMsg: 'development error')
		return false
		}

	Name: 'NOT DEFINED' // Defined by child class
	Getter_AddonControl()
		{

		if false is injectControl = .injectControl()
			return false
		if false is ctrl = injectControl.FindControl(.Name)
			return false
		return .AddonControl = ctrl
		}

	Addon_RedirMethods()
		{ return #() }

	Getter_SubSplit()
		{
		if false is .AddonControl
			return false
		return .SubSplit = .AddonControl.Parent.Parent
		}
	}
