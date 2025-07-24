// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	dict: #(0: '', 1: 'un', 2: 'deux', 3: 'trois', 4: 'quatre',
		5: 'cinq', 6: 'six', 7: 'sept',
		8: 'huit', 9: 'neuf', 10: 'dix', 11: 'onze', 12: 'douze', 13: 'treize',
		14: 'quatorze', 15: 'quinze', 16: 'seize', 17: 'dix-sept', 18: 'dix-huit',
		19: 'dix-neuf', 20: 'vingt', 30:'trente', 40: 'quarante',50: 'cinquante',
		60: 'soixante', 70: 'soixante-dix', 80: 'quatre-vingt', 90: 'quatre-vingt-dix')
	ten: 10
	hundred: 100
	thousand: 1000
	million: 1000000
	billion: 1000000000
	NumberToWords(number)
		{
		if number > 2147483647 /*= max number */
			return "nombre trop grand"
		number = number.Int()
		texte = number is 0
			? "zero"
			: .greaterThan0(number)
		return .checkGrammar(texte)
		}
	seventy: 70
	greaterThan0(number)
		{
		if number <= 20 /*= twenty */
			return .dict[number]
		if number < .seventy
			return .dict[number - number % .ten] $
				(number % .ten > 0 ? " " $ (number % .ten is 1 ? "et " : "") $
				.dict[number%.ten] : "")
		if number is .seventy
			return .dict[number]
		return .greaterThan70(number)
		}
	eighty: 80
	ninety: 90
	greaterThan70(number)
		{
		if (number < .eighty)
			return .dict[number - (number % .ten + .ten)] $
				(number % .ten + .ten > .ten ? " " $
				(number % .ten + .ten is 11 ? 'et ' : '') $ /*= eleven */
				.dict[number % .ten + .ten] : "")
		if (number is .eighty)
			return .dict[.eighty] $ 's'
		if (number < .ninety)
			return .dict[number - number % .ten] $
				(number % .ten > 0 ? " " $ .dict[number % .ten] : "")
		if (number is .ninety)
			return .dict[number]
		return .greaterThan90(number)
		}
	greaterThan90(number)
		{
		if (number < .hundred)
			return .dict[number - (number % .ten + .ten)] $
				(number % .ten + .ten > .ten ? " "  $
				.dict[number%.ten + .ten] : "")
		if (number < 200) /*= two hundred */
			return  "cent" $ (number % .hundred is 0 ? "" : " " $
				(number % .hundred).EnFrancais())
		.greaterThan200(number)
		}
	greaterThan200(number)
		{
		if (number < .thousand)
			return  .dict[(number / .hundred).Int()] $ " cent" $
				((((number / .hundred).Int() > 1) and
					(number % .hundred is 0)) ? "s" : "") $
				(number % .hundred is 0? "" : " " $ (number % .hundred).EnFrancais())
		return .greaterThanThousand(number)
		}
	greaterThanThousand(number)
		{
		if (number < .million)
			return  ((number / .thousand).Int() is 1 ? "" :
				((number/.thousand).Int()).EnFrancais() $ " ") $ "mille" $
				(number % .thousand is 0 ? "" : " " $ (number%.thousand).EnFrancais())
		if (number < .billion)
			return  ((number / .million).Int()).EnFrancais() $ " million" $
				((number /.million).Int() > 1 ? "s" : "") $
				(number % .million is 0 ? "" : " " $ (number % .million).EnFrancais())
		return .greaterThanBillion(number)
		}
	greaterThanBillion(number)
		{
		return  ((number / .billion).Int()).EnFrancais() $ " milliard" $
			((number/.billion).Int() > 1 ? "s" : "") $
			(number%.billion is 0 ? "" : " " $ (number % .billion).EnFrancais())
		}
	checkGrammar(texte)
		{
		// Vérifications grammaticales
		chaine = texte.Split(' ')
		if chaine.Find('cents') isnt false
			if chaine.Find('mille') is chaine.Find('cents') + 1
				chaine[chaine.Find('cents')] = 'cent'
		if chaine.Find('quatre-vingts') isnt false
			if chaine.Find('mille') is chaine.Find('quatre-vingts') + 1
				chaine[chaine.Find('quatre-vingts')] = 'quatre-vingt'
		// fin vérifications grammaticales
		return chaine.Join(' ')
		}
	}