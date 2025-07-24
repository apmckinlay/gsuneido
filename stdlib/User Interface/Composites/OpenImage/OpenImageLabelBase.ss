// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
/* USAGE
Label: 		Child classes MUST define Label as it is used as the LinkButton text
Condition:	Override this method to specify when to use or output a label
	(output via SystemAttachmentLabels). IF the condition is not met after the
	label has been output, users will be able to remove it from the
	attachment_labels table as it will no longer be necessary.
	NOTE: Condition should ALWAYS return true or false

By default, SystemAttachmentLabels.Ensure is run when installing
new systems / applications. However, if there are specific events which can toggle the
usage of a label, [OpenImageLabelBase Child].Ensure() should be called directly.
This will reduce the overall overhead when enabling new features
*/
class
	{
	Tip: 		false 	// Tool tip explaining this label's use
	Command:	'' 		// Command to trigger on clicking link
	Column:		0		// Can be 0 ("email with ..."), 1 ("send with ..."), or 2 (other)
	Layout()
		{ return [#LinkButton, name: .Label, command: .Command, tip: .Tip] }

	Condition()
		{ return true }

	Ensure(t = false, skipResetCache = false)
		{
		if TestRunner.RunningTests?()
			return false
		AttachmentLabels.OutputLabels(.Label, t)
		if not skipResetCache
			SystemAttachmentLabels.ResetCache()
		return true
		}
	}
