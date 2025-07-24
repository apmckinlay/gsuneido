// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	outstanding?: false
	calculating?: false
	newTask: false
	taskCounter: 0
	Name: 'ThreadTotal'
	New(layout, .calcOnserverFn, .screenForTotals)
		{
		super(Object(RecordControl
			Object(#Horz
				Object('Static', name: 'statusText')
				Object('EnhancedButton'
					image: 'triangle-warning', name: 'status', imagePadding: .15
					imageColor: CLR.darkorange,
					mouseOverImageColor: CLR.darkorange,
					tip: 'Calculating...')
				#(Skip small:)
				layout)))
		.status = .FindControl('status')
		.statusText = .FindControl('statusText')
		.run(.calcOnserverFn, .screenForTotals,
			'ThreadTotal_' $ .calcOnserverFn $ '_' $ Display(Timestamp()))
		}

	AfterChanged(saved, query, dirty?, filters)
		{
		// when not saving we need to update the status on UI
		if saved is false
			.outstanding? = dirty?
		else
			{
			.outstanding? = false
			.calculating? = true
			.newTask = Object(:query, filters: filters.DeepCopy(), id: ++.taskCounter)
			}
		.updateState()
		}

	run(calcOnserverFn, screenForTotals, name)
		{
		Thread({ .runCalc(calcOnserverFn, screenForTotals) }, :name)
		}

	runCalc(calcOnserverFn, screenForTotals)
		{
		forever
			{
			if false is task = .fetchNewTask()
				return

			try	// do the calculation
				x = .doCalc(calcOnserverFn, task.query, screenForTotals,
					task.filters)
			catch (e)
				{
				SuneidoLog('ERROR: ThreadTotal calculation failed - ' $ e, params: task)
				continue
				}

			.doneCalc(x, task.id)
			}
		}

	fetchNewTask()
		{
		forever
			{
			if .Destroyed?()
				return false

			if .newTask isnt false
				{
				task = .newTask
				.newTask = false
				return task
				}

			.sleep()
			}
		}

	// extracted for test
	sleep()
		{
		Thread.Sleep(500 /*=sleep time*/)
		}

	doCalc(calcOnserverFn, query, screenForTotals, filters)
		{
		return ThreadTotalCached(calcOnserverFn, query, screenForTotals, filters)
		}

	// extracted to override in test
	setData(x)
		{
		.Data.Set(x)
		}

	updateState(x = false)
		{
		if .Destroyed?()
			return

		if x isnt false
			.setData(x)

		if .calculating?
			.updateStatus('Calculating...', CLR.darkorange, CLR.WarnColor)
		else if .outstanding?
			.updateStatus('Value Changed', CLR.red, CLR.ErrorColor,
				'Totals do not reflect changes to current record')
		else
			{
			.setStatus(false, '')
			.setBackgroundColor(CLR.ButtonFace)
			}
		}

	updateStatus(text, imageColor, statusColor, tooltip = false)
		{
		.setStatus(true, text)
		.status.ToolTip(tooltip is false ? text : tooltip)
		.status.SetImageColor(imageColor, imageColor)
		.setBackgroundColor(statusColor)
		}

	setBackgroundColor(color)
		{
		for f in .Data.GetControlData().Members()
			if ((false isnt ctrl = .FindControl(f)) and ctrl.Method?('SetReadOnlyBrush'))
				{
				ctrl.SetReadOnlyBrush(color)
				ctrl.Repaint()
				}
		}

	setStatus(status, statusText)
		{
		.status.SetVisible(status)
		.statusText.Set(statusText)
		}

	doneCalc(x, id)
		{
		if .Destroyed?()
			return

		.Defer()
			{
			.calculating? = .taskCounter > id
			.updateState(x)
			}
		}

	SetVisible(visible?)
		{
		super.SetVisible(visible?)
		if visible? and .statusText.Get() is ''
			.status.SetVisible(false)
		}
	}
