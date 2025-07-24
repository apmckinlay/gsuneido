// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
/* USAGE: to add an internal/system label, define "OpenImage_[label name]Label".
	These classes should inherit from OpenImageLabelBase and must define "Label: [text]"

NOTE: When sending new OpenImage label classes, either:
	A. ensure the library via: SystemAttachmentLabels.Ensure(library)
	or
	B. ensure the label directly via: [OpenImageLabel Class].Ensure()
*/
MemoizeSingle
	{
	Func()
		{
		labels = Object()
		.loop({ labels.AddUnique(it.Label) })
		return labels.Remove('')
		}

	loop(block, libraries = false)
		{
		libsInUse = Libraries()
		if libraries is false
			libraries = libsInUse
		for lib in libraries
			if not libsInUse.Has?(lib)
				.logError('cannot process library: ' $ lib $ ' (not in use)')
			else
				try
					.processLib(lib, block)
				catch (e)
					.logError(lib $ ' encountered error: ' $ e)
		}

	logError(msg)
		{ SuneidoLog('ERROR: (CAUGHT) ' $ msg, calls:, caughtMsg: 'internalLabels') }

	processLib(lib, block)
		{
		QueryApply(.query(lib), group: -1)
			{
			labelClass = Global(it.name)
			if labelClass.Condition()
				block(labelClass)
			}
		}

	query(lib)
		{ return lib $ ' where name.Prefix?(`OpenImage_`) and name.Suffix?(`Label`)' }

	Ensure(lib = false, t = false)
		{
		if lib is false or Libraries().Has?(lib)
			.ensure(lib, t)
		else
			.logError('cannot ensure labels for library: ' $ lib $ ' (not in use)')
		}

	ensure(lib, t)
		{
		if not TableExists?(#attachment_labels)
			AttachmentLabels.Ensure()
		libraries = lib isnt false ? Object(lib) : false
		.loop({ it.Ensure(t, skipResetCache:) }, :libraries)
		.ResetCache()
		}
	}
