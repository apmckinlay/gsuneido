// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
/*
This control is used to handle Ok/Cancel Screen in a generic mannor
It will handle Adding the Ok/Cancel buttons to the screen, and the event processing
for the buttons.
*/
Controller
	{
	New(prompt)
		{
		super(.layout(prompt))
// TEMPORARY - Issue 28840 - logging to track down no return value error
.whatami = Object?(prompt) ? prompt[0] : prompt
		.v = .FindControl('Vert')
		// alow parent to override OkCancel Placement
		// note - this also finds AllNoneOkCancelControl (e.g. ChooseManyAsObjectList)
		// need to check for OK control as well
		// (some screens put a clear option between Ok and Cancel)
		if false is .FindControl('OkCancel') and false is .FindControl('OK')
			{
			.v.Append(#Skip)
			.v.Append(#OkCancel)
			}

		.Redir('OK', .v.GetChildren()[0])
		.Redir('Cancel', .v.GetChildren()[0])
		}
	layout(prompt)
		{
		vert = Object('Vert')
		vert.Add(prompt)
		return vert
		}
	// Sends "OK"
	// to just handle true do not implement OK
	// if implemented return the value that is to be set for .Window.Result
	// a return value of false indicates the data was invalid
	// (.Window.Result will not be set in that case)
	// NOTE: if the parrent has overriden the OkCancel Placement,
	// (and they are not a PassthruController) they will need:
	/*
	On_OK()
		{
		.Send("On_OK")
		}
	*/
	allowDestroy?: false
	On_OK()
		{
// TEMPORARY - Issue 28840 - logging to track down no return value error
.title = .Send('GetTitle')
		if false is x = .v.Send('OK')
			return // invalid data
		.allowDestroy? = true
		if x is 0
			.Window.Result(true)
		else
			.Window.Result(x)
		}
	On_Cancel()
		{
		x = .v.Send('Cancel')
		.allowDestroy? = true
		if x is 0
			.Window.Result(false)
		else
			.Window.Result(x)
		}
	ConfirmDestroy(@args /*unused*/)
		{
		return .allowDestroy?
		}
	}
