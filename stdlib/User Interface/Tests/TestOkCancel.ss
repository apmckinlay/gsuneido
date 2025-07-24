// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Window()
		{
		// buttons should NOT work, escape should not work

		Window(TestOkCancel.Regular)

		Window(TestOkCancel.Controller1)

		Window(TestOkCancel.Controller2)

		Window(TestOkCancel.Controller3)

		Window(TestOkCancel.Controller4)
		}

	Dialog()
		{
		// ok, cancel, and Esc should all work

		Dialog(0, TestOkCancel.Regular)
			// buttons work but Esc broken with defaults in WindowBase
			// Esc requires Control.Msg check
				// WindowBase.COMMAND => WindowBase.Send => Control.Msg

		Dialog(0, TestOkCancel.Controller1)
			// works with defaults in WindowBase if PassthruController

		Dialog(0, TestOkCancel.Controller2)

		Dialog(0, TestOkCancel.Controller3)

		Dialog(0, TestOkCancel.Controller4)
			// works with defaults in WindowBase
			// Control.Msg check NOT required
		}

	ModalWindow()
		{
		// ok, cancel, and Esc should all work

		ModalWindow(TestOkCancel.Regular)

		ModalWindow(TestOkCancel.Controller1)

		ModalWindow(TestOkCancel.Controller2)

		ModalWindow(TestOkCancel.Controller3)

		ModalWindow(TestOkCancel.Controller4)
		}

	Regular: (Vert
			Field
			(Static 'OK should return true, Cancel/Esc should return false')
			Skip OkCancel xstretch: 0)

	Controller1: PassthruController // default ok/cancel methods
		{
		Controls: (Vert
			(Static 'OK should return true, Cancel/Esc should return false')
			Skip (Record OkCancel) xstretch: 0)
		}
	Controller2: PassthruController
		{
		Controls: (Vert
			(Static 'OK should return "OK", Cancel/Esc should return false')
			Skip (Record OkCancel) xstretch: 0)
		On_OK()
			{
			.Window.Result('OK')
			}
		}
	Controller3: PassthruController
		{
		Controls: (Vert
			(Static 'OK should return true, Cancel/Esc should return "Cancel"')
			Skip (Record OkCancel) xstretch: 0)
		On_Cancel()
			{
			.Window.Result('Cancel')
			}
		}
	Controller4: Controller // custom ok/cancel methods
		{
		Controls: (Vert
			(Static 'OK should return "OK", Cancel/Esc should return "Cancel"')
			Skip (Record OkCancel) xstretch: 0)
		On_OK()
			{
			.Window.Result('OK')
			}
		On_Cancel()
			{
			.Window.Result('Cancel')
			}
		}
	}