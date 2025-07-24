// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
// Author: Victor Schappert
Container
	{
	// =========================================================================
	// THIS CONTAINER IS SIMILAR TO 'GROUP' IN THAT AT BASE IT PROVIDES A SIMPLE
	// horizontal or vertical stack of its constituent controls within the boxes
	// and stretch paradigm. However, it has certain features which give the GUI
	// programmer more finely-grained control over the layout of individual
	// controls so that the layout can be adjusted fluidly at runtime. These
	// features are:
	//     - It is aware of the location of its own bounding rectangle.
	//     - It is aware of the sizes of its constituent controls along the
	//       layout dimension (horizontal or vertical).
	//     - It supports assigning arbitrary sizes to a constituent control
	//       along the layout dimension so long as the size assigned does not
	//       force any control to shrink below its minimum size.
	//     - It supports showing/hiding/adding/removing/replacing constituent
	//       controls.
	// This container always respects the .[X|Y]Min member of the constituent
	// controls. However, one consequence of this container's layout strategy is
	// that a constituent control's .[X|Y]Stretch member is only considered when
	// the control is initially constructed. After that point, this container
	// maintains its own internal data values for the ratio of space along the
	// layout dimension which the contained control is allotted.
	// =========================================================================
	// This control does not implement baseline alignment (i.e. it does not
	// respect the Control.Top and Control.Left members) since that behaviour is
	// not compatible with this control's ability to show/hide contained
	// controls on demand. If baseline alignment is needed for a set of
	// controls, use a contained Vert or Horz.
	// =========================================================================


	// DATA

	Name:          "LinearLayout"
	Dir:           #vert /* Layout dimension: #vert or #horz */

	rc:            false /* LAZY | Rect of this in parent Hwnd's client area */
	xmin0:         false
	ymin0:         false
	set_xstretch?: false
	set_ystretch?: false

	ctrls:         false /* List of all controls, visible or not */
	visCtrls:      false /* Boolean vector of same size as .ctrls */
	ratios:        false /* Vector of ratios of control size to nfcs */


	// CONSTRUCTORS

	New(@args)
		{
		// Determine the control direction -- #vert or #horz
		dir = args.GetDefault("dir", "vert")
		switch (dir)
			{
		case "vert": Assert(.Dir is: dir)
		case "horz": .Dir = #horz
		default: throw "invalid direction"
			}
		// Get the the control specs
		containedCtrlSpecs = args.Values(list:)
		// Initialize
		.initMembers()
		.initCtrls(containedCtrlSpecs)
		.initCtrlVisibility(args.GetDefault("visible", #()))
		.initRatios()
		.updateFrameworkVars()
		}


	// OVERRIDES FOR ANCESTOR CLASS: Control

	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.rc.Set(x, y, w, h)
		sizes = .calcCtrlSizeVector(.ratios, .visibleIndices())
		.resizeCtrlsFromVector(sizes)
		}
	GetChildren()
		{ return .ctrls }


	// PUBLIC INTERFACE

	Tally(all = false)
		{
		Assert(all, isBoolean:)
		return all ? .ctrls.Size() : .visCtrls.Count(true)
		}
	GetControlSize(indexOrName)
		{
		i = .mapIndex(indexOrName)
		return .visCtrls[i]
			? .calcCtrlSizeVector(.ratios, .visibleIndices())[i]
			: 0
		}
	GetControlPos(indexOrName)
		{
		i = .mapIndex(indexOrName)
		sizes = .calcCtrlSizeVector(.ratios, .visibleIndices())
		pos = 0
		for (k = 0; k < i; ++k)
			if .visCtrls[k]
				pos += sizes[k]
		return pos
		}
	GetControlResizePoint(indexOrName, requestedSize)
		{
		i = .mapIndex(indexOrName)
		// Searches for the first size between the control's current size and
		// the requested size that is a valid size for the control.
		validSize = .GetControlSize(i)
		sizeMap = Object()
		while requestedSize isnt validSize
			{
			midpoint = ((requestedSize - validSize) / 2).Round(0)
			midpointSize = validSize + midpoint
			sizeMap[i] = midpointSize
			if .SetControlSizes?(sizeMap) // If it worked, try closer to requested size
				validSize = midpointSize
			else
				requestedSize = midpointSize - midpoint.Sign()
			}
		return validSize
		}
	SetControlSizes?(sizeMap/*map of index or name => integer size*/)
		{ false isnt .testSetCtrlSizes(.mapIndices(sizeMap), .visibleIndices(), .ratios) }
	SetControlSizes(sizeMap/*map of index or name => integer size*/)
		{
		visibleIndices = .visibleIndices()
		newSizes = .testSetCtrlSizes(.mapIndices(sizeMap), visibleIndices, .ratios)
		if false is newSizes
			throw "size map invalid for this linear layout"
		.ratios = .calcRatiosFromSizes(newSizes, visibleIndices) // Back out the ratios
		.resizeCtrlsFromVector(newSizes) // Resize the controls
		}
	SetControlVisibilities?(visMap/*map of index or name => boolean visible*/)
		{
		visibleIndices = .visibleIndicesFromMap(.mapIndices(visMap))
		return false isnt .testSetCtrlVisibilities(visibleIndices)
		}
	SetControlVisibilities(visMap/*map of index or name => boolean visible*/)
		{
		visMap = .mapIndices(visMap)
		visibleIndices = .visibleIndicesFromMap(visMap)
		newSizes = .testSetCtrlVisibilities(visibleIndices)
		if false is newSizes
			throw "visibility map invalid for this linear layout"
		.ratios = .calcRatiosFromSizes(newSizes, visibleIndices) // Back out the ratios
		.applyVisibility(visMap) // Apply visibility to controls
		.resizeCtrlsFromVector(newSizes) // Resize the controls
		}
	SetControlVisibilitiesAndSizes?(visMap, sizeMap)
		{
		visibleIndices = .visibleIndicesFromMap(.mapIndices(visMap))
		newSizes = .testSetCtrlVisibilities(visibleIndices)
		if false is newSizes
			return false
		ratios = .calcRatiosFromSizes(newSizes, visibleIndices)
		return false isnt .testSetCtrlSizes(.mapIndices(sizeMap), visibleIndices, ratios)
		}
	SetControlVisibilitiesAndSizes(visMap, sizeMap)
		{
		visMap = .mapIndices(visMap)
		visibleIndices = .visibleIndicesFromMap(visMap)
		newSizes = .testSetCtrlVisibilities(visibleIndices)
		if false is newSizes
			throw "visibility map invalid for this linear layout"
		ratios = .calcRatiosFromSizes(newSizes, visibleIndices) // Back out the ratios
		newSizes = .testSetCtrlSizes(.mapIndices(sizeMap), visibleIndices, ratios)
		if false is newSizes
			throw "size map invalid for this linear layout"
		.ratios = .calcRatiosFromSizes(newSizes, visibleIndices) // Back out the ratios
		.applyVisibility(visMap) // Apply visibility to controls
		.resizeCtrlsFromVector(newSizes) // Resize the controls
		}
	ControlVisible?(indexOrName)
		{ .visCtrls[.mapIndex(indexOrName)] }


	// INTERNALS

	initMembers()
		{
		.rc = Rect(0, 0, 0, 0)
		.xmin0 = .Xmin
		.ymin0 = .Ymin
		if false is .Xstretch // default in Control
			.set_xstretch? = true
		if false is .Ystretch // default in Control
			.set_ystretch? = true
		}
	initCtrls(ctrlSpecs)
		{
		.ctrls = Object()
		for cs in ctrlSpecs
			.ctrls.Add(.Construct(cs))
		}
	initCtrlVisibility(visMap)
		{
		.visCtrls = Object().AddMany!(true, .Tally(all:))
		.applyVisibility(.mapIndices(visMap))
		}
	initRatios()
		{ .ratios = .calcRatiosInit() }
	findCtrlIndex(name)
		{
		n = .Tally(all:)
		found = false
		for (i = 0; i < n; ++i)
			{
			if name is .ctrls[i].Name
				if false isnt found
					throw "duplicate control name: " $ Display(name)
				else
					found = i
			}
		return found
		}
	mapIndex(indexOrName)
		{
		if (Number?(indexOrName) and indexOrName.Int?() and
			0 <= indexOrName and indexOrName < .Tally(all:))
			return indexOrName
		else if String?(indexOrName)
			return .findCtrlIndex(indexOrName)
		else
			throw "no such control: " $ Display(indexOrName)
		}
	mapIndices(map)
		// pre:  Map is a #(key => value) map where every 'key' is either the
		//       unique string name of a contained control or the unique integer
		//       index of a contained control
		// post: The return value is a #(key => value) map that is the same as
		//       (pre)map except that every string 'key' from (pre)map has been
		//       replaced with the integer index of the control with that name
		{ map.MapMembers() { .mapIndex(it) } }
	updateFrameworkVars()
		{
		// Updates .Xmin, .Ymin, .MaxHeight; also updates. Xstretch if
		// .set_xstretch? and .Ystretch if .set_ystretch?. This method is
		// roughly equivalent to Group.Recalc()
		xmin = 0
		ymin = 0
		xstretch = false
		ystretch = false
		if #vert is .Dir
			{
			maxheight = 0
			for ctrl in .visibleCtrls()
				{
				xmin = Max(xmin, ctrl.Xmin)
				xstretch = Max(xstretch, ctrl.Xstretch)
				ymin += ctrl.Ymin
				ystretch += ctrl.Ystretch
				maxheight += ctrl.MaxHeight
				}
			.MaxHeight = maxheight
			}
		else // horz
			{
			for ctrl in .visibleCtrls()
				{
				xmin += ctrl.Xmin
				xstretch += ctrl.Xstretch
				ymin = Max(ymin, ctrl.Ymin)
				ystretch = Max(ystretch, ctrl.Ystretch)
				}
			}
		.Xmin = Max(xmin, .xmin0)
		.Ymin = Max(ymin, .ymin0)
		if .set_xstretch?
			.Xstretch = xstretch
		if .set_ystretch?
			.Ystretch = ystretch
		}
	ldSize()
		{ #vert is .Dir ? .rc.GetHeight() : .rc.GetWidth() }
	ldMin(ctrl)
		{ #vert is .Dir ? ctrl.Ymin : ctrl.Xmin }
	ldStretch(ctrl)
		{ #vert is .Dir ? ctrl.Ystretch : ctrl.Xstretch }
	stretchable?(ctrl)
		{ 0 < .ldStretch(ctrl) }
	visibleCtrls(visibleIndices = false)
		{
		if false is visibleIndices
			visibleIndices = .visibleIndices()
		return .ctrls.Project(visibleIndices)
		}
	visibleIndices()
		{ .visCtrls.MembersIf() { .visCtrls[it] } }
	visibleIndicesFromMap(visMap)
		{ .visCtrls.MembersIf() { visMap.GetDefault(it, .visCtrls[it]) } }
	stretchableIndices()
		{ .ctrls.MembersIf({|i| .stretchable?(.ctrls[i]) }) }
	calcFreeSpace(visibleIndices)
		{
		// The return value is an integer which is equal to the 'free space'.
		// The 'free space' is defined as the difference between the sum of the
		// all the .[X|Y]Min values of the visible controls along the layout
		// dimension minus the container size along the layout dimension.

		// NOTE: The value returned may be less than zero if the container size
		//       is too small to fit the contained controls even at their
		//       minimum sizes.
		.ldSize() - .visibleCtrls(visibleIndices).SumWith(.ldMin)
		}
	calcTotalStretch()
		{ .visibleCtrls().SumWith(.ldStretch) }
	calcRatiosInit()
		{
		// Return value is an array of non-negative real numbers where the i'th
		// number represents the ratio of the i'th control's stretch to the
		// total of all visible stretch along the layout dimension.
		totalStretch = .calcTotalStretch()
		n = .ctrls.Size()
		ratios = Object()
		for (k = 0; k < n; ++k)
			{
			ratios.Add(.visCtrls[k]
				? Min(.ldStretch(.ctrls[k]) / totalStretch, 1)
				: 0)
			}
		return ratios
		}
	calcRatiosFromSizes(sizes, visibleIndices)
		{
		// Note that the ratios of the visible stretchable controls will always
		// sum up to 100% but if you add in the ratios of the invisible
		// stretchable controls you will get a number bigger than 100%. This is
		// deliberate. When an invisible control is made visible, all the ratios
		// will be scaled down so that the new set of visible controls still
		// sums to 100%.
		fs = .calcFreeSpace(visibleIndices)
		if fs < 1
			ratios = .ratios
		else
			{
			n = .ctrls.Size()
			ratios = Object()
			for (k = 0; k < n; ++k) // invisible controls get a ratio too
				ratios.Add(Min((Max(sizes[k] - .ldMin(.ctrls[k]), 0)) / fs, 1))
			}
		return ratios
		}
	calcCtrlSizeVector(ratios, visibleIndices)
		{
		// The return value is a vector of integer control sizes along the
		// layout dimension. Each index in the vector corresponds to the index
		// of a control, so the vector contains sizes for invisible as well as
		// visible controls. The sum of the sizes of the visible controls in the
		// vector may exceed the container size along the layout dimension if
		// the container is too small to fit all of the contained controls at
		// their minimum sizes.
		fs = .calcFreeSpace(visibleIndices)
		if fs < 1
			return .ctrls.Map(.ldMin)
		else
			{
			n = .Tally(all:)
			sizes = Object().AddMany!(false, n)
			for k in visibleIndices
				sizes[k] = true
			remFs = fs
			// STEP 1: Allocate according to ratios (there might be some
			//         rounding error).
			for (k = 0; k < n; ++k)
				{
				ctrl = .ctrls[k]
				ctrlFs = Min(remFs, (fs * ratios[k]).Round(0))
				if sizes[k]
					remFs -= ctrlFs
				sizes[k] = #vert is .Dir
					? Min(ctrl.Ymin + ctrlFs, ctrl.MaxHeight)
					: ctrl.Xmin + ctrlFs
				}
			// STEP 2: Greedily allocate any free space that was missed due to
			//         rounding error when multiplying by the ratios. This space
			//         can only be allocated among the controls that are both
			//         visible and stretchable.
			for k in visibleIndices.Intersect(.stretchableIndices())
				{
				Assert(remFs, greaterThanOrEqualTo: 0)
				if 0 is remFs
					break
				addFs = #vert is .Dir
					? Min(remFs, .ctrls[k].MaxHeight - sizes[k])
					: remFs
				sizes[k] += addFs
				remFs -= addFs
				}
			return sizes
			}
		}
	resizeCtrlsFromVector(sizes)
		{
		n = .Tally(all:)
		Assert(n, is: sizes.Size())
		x = .rc.GetX()
		y = .rc.GetY()
		w = .rc.GetWidth()
		h = .rc.GetHeight()
		for (k = 0; k < n; ++k)
			{
			if not .visCtrls[k]
				continue
			ctrl = .ctrls[k]
			size = sizes[k]
			if #vert is .Dir
				{
				ww = 0 is ctrl.Xstretch ? ctrl.Xmin : w
				ctrl.Resize(x, y, ww, size)
				y += size
				}
			else // horz
				{
				hh = 0 is ctrl.Ystretch ? ctrl.Ymin : h
				ctrl.Resize(x, y, size, hh)
				x += size
				}
			}
		}
	applyVisibility(visMap)
		{
		for i in visMap.Members()
			{
			before = .visCtrls[i]
			after  = visMap[i]
			if before isnt after
				{
				.ctrls[i].SetVisible(after)
				.visCtrls[i] = after
				}
			}
		}
	testSetCtrlSize(ctrl, size)
		{ // Returns 0 if ctrl can be set to size
		Assert(size, isInt:)
		if #vert is .Dir
			{
			if size < ctrl.Ymin
				return size - ctrl.Ymin
			else if ctrl.CalcMaxHeight() < size
				return size - ctrl.CalcMaxHeight()
			}
		else if size < ctrl.Xmin
			return size - ctrl.Xmin
		return 0
		}
	testSetCtrlSizes(sizeMap, visibleIndices, ratios)
		{
		newSizes = .calcCtrlSizeVector(ratios, visibleIndices)
		affectedIndices = sizeMap.Members()

		// The first and easiest process is to check whether any proposed
		// control size directly violates the min/max size constraints on that
		// control.

		for (i in affectedIndices)
			{
			if 0 is .testSetCtrlSize(.ctrls[i], sizeMap[i])
				newSizes[i] = sizeMap[i]
			else
				return false
			}

		// The second and harder step is to check whether the proposed control
		// sizes indirectly violate the min/max size constraints of another
		// control (because that other control would be forced to shrink or
		// grow). Another way of looking at this is whether we can allocate the
		// change in free space caused by the control size change.

		// SUBSTEP1: Determine the free space delta that needs to be allocated
		fsDelta = newSizes.Project(visibleIndices).Sum() - .rc.GetHeight()

		// SUBSTEP2: Get the indices of all the visible, stretchable, unaffected
		//           controls (unaffected meaning whose size is not being set).
		vsuIndices = visibleIndices.
			Intersect(.stretchableIndices()).
			Difference(affectedIndices)

		// SUBSTEP3: Get the sum of the ratios for those controls.
		sumRatio = vsuIndices.SumWith{ ratios[it] }

		// SUBSTEP4: Try to allocate the delta among the controls using only the
		//           ratios.
		fsDeltaRemaining = fsDelta
		if 0 < sumRatio
			{
			for (i in vsuIndices)
				{
				ratio = ratios[i] / sumRatio
				ctrlDelta = (fsDelta * ratio).Round(0)
				leftover = .testSetCtrlSize(.ctrls[i], newSizes[i] - ctrlDelta)
				ctrlDelta -= leftover
				newSizes[i] -= ctrlDelta
				fsDeltaRemaining -= ctrlDelta
				}
			}

		// SUBSTEP5: For any remaining delta that couldn't be allocated using
		//           only the ratios, try to allocate it greedily.
		for (i in vsuIndices)
			{
			ctrlDelta = fsDeltaRemaining
			leftover = .testSetCtrlSize(.ctrls[i], newSizes[i] - ctrlDelta)
			ctrlDelta -= leftover
			newSizes[i] -= ctrlDelta
			fsDeltaRemaining -= ctrlDelta
			if 0 is fsDeltaRemaining
				break
			}

		if 0 isnt fsDeltaRemaining
			return false

		// If we got to this point without returning false (failure), return the
		// new size array. This can be used to back out the new ratios if we are
		// actually applying the size change. Otherwise, it is just an indicator
		// of success.
		return newSizes
		}
	testSetCtrlVisibilities(visibleIndices)
		{
		// STEP 1: Compute the visible ratio sum, which must be normalized to
		//         1.0.
		visibleRatioSum = visibleIndices.SumWith{ .ratios[it] }
		// STEP 2: Normalize all the ratios so that the visible control ratios
		//         sum to 1.0 and the invisible control ratios are on the same
		//         scale as the visible ones.
		ratios = 0 < visibleRatioSum
			? .ratios.Map() { it / visibleRatioSum }
			: .ratios
		// STEP 3: Calculate the control size vector that would result from the
		//         new visibility settings.
		newSizes = .calcCtrlSizeVector(ratios, visibleIndices)
		// STEP 4: Test setting the control sizes to the control size vector.
		return .testSetCtrlSizes(newSizes, visibleIndices, .ratios)
		}
	}