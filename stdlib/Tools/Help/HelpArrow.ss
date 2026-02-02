// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
class
	{
	right: `<svg width="12px" height="12px" viewBox="0 0 6.8057642 7.3742032" ` $
		`version="1.1" id="svg1" xmlns="http://www.w3.org/2000/svg">` $
		`<g transform="translate(-29.676908,-22.386642)">` $
		`<path d="m 32.064352,26.07116 -2.387444,-3.684518 6.805764, ` $
		`3.684518 -6.805764,3.689685 z m 0.490924,0 h 3.240102 l -5.2038,-2.826692 z"` $
		` id="text1" fill="currentColor" /></g></svg>`

	left: `<svg width="12px" height="12px" viewBox="0 0 25.722656 27.871094" ` $
		`version="1.1" xmlns="http://www.w3.org/2000/svg" fill="currentColor">` $
		`<path d="M 25.722656,27.871094 0,13.925781 25.722656,0 16.699219,13.925781 ` $
		`Z M 22.265625,3.2421875 2.5976563,13.925781 H 14.84375 Z" id="text1" /></svg>`

	up: `<svg width="12px" height="12px" viewBox="0 0 7.3742034 6.8057644" ` $
		`version="1.1" xmlns="http://www.w3.org/2000/svg" fill="currentColor">` $
		`<path d="M 2.4214387e-8,6.8057644 3.6896855,0 7.3742034,6.8057644 ` $
		`3.6896855,4.4183208 Z M 6.5163774,5.8910945 3.6896855,0.68729432 ` $
		`V 3.9273963 Z" id="text1" /></svg>`

	down: `<svg width="12px" height="12px" viewBox="0 0 7.374203 6.805764" ` $
		`version="1.1" xmlns="http://www.w3.org/2000/svg" fill="currentColor">` $
		`<path d="M 7.374203,0 3.684518,6.805764 0,0 3.684518,2.387443 ` $
		`Z M 0.857826,0.91467 3.684518,6.11847 V 2.878368 Z" id="text1" /></svg>`

	CallClass(direction = false)
		{
		switch (direction)
			{
			case 'left':
				svg = .left
			case 'up':
				svg = .up
			case 'right':
				svg = .right
			case 'down':
				svg = .down
			default:
				svg = ''
			}
		return svg
		}
	}