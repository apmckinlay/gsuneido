// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// NOTE: flags only used when prompt is just a string
// This is the dialog, see OkCancelControl for the buttons
function (prompt = "", title = "", onDestroy = false)
	{
	return ModalWindow([OkCancelWrapper, prompt], :title, closeButton?: false, :onDestroy)
	}
