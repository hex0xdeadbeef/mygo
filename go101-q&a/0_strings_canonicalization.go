package main

import "unique"

/*
	HOW TO CANONICALIZE STRINGS TO SAVE MEMORY?

	At run time of Go programs, sometimes, some equal strings don't share underlying bytes
	memory blocks, even if they can share a single common bytes memory block.
*/

/*
	WAY 1: Canonicalize two strings when they're found equal

	The logic implementation is simple:

		if str1 == str2 {
			str1 = str2
		}
*/

func CanonicalizeStringsSimple(ss []string) {
	type s struct {
		str   string
		index int
	}

	var temp = make([]s, len(ss))

	// Just fill all the strings into temp
	for i := range temp {
		temp[i] = s{
			str:   ss[i],
			index: i,
		}
	}

	for i := 0; i < len(temp); {
		var k = i + 1

		for j := k; j < len(temp); j++ {
			if temp[j].str == temp[i].str {
				temp[j].str = temp[i].str
				temp[k], temp[j] = temp[j], temp[k]
				k++
			}
		}
		i = k
	}

	for i := range temp {
		ss[temp[i].index] = temp[i].str
	}

}

/*
	WAY 2: Use Go 1.23 introduced unique.Handle

	Note: the unique.Make way isn't always suitable for every situation. The unique.Make
	function will allocate a backing bytes memory block for each distinct string. So if some
	unequal strings to be canonicalized share the same backing bytes memory block, the unique.Make
	function will allocate a new backing byte sequence memory block for each of the strings. Doing
	this actually allocates more memory (than using a single memory block)

	В Go строка (string) устроена как указатель на участок памяти с байтами и длина. Если вы создаёте несколько строковых значений, которые по-сути являются разными «видами» одного и того же большого слайса байт (например, берёте срезы одной и той же []byte), то все эти строки разделяют один и тот же буфер в памяти.

	Функция unique.Make (или любая аналогичная функция «интернирования») обычно для каждого уникального значения вашей строки выделяет новый буфер и копирует в него содержимое. То есть:
		1.	Вы дали ей, скажем, три подстроки из одного большого массива байт.
		2.	Эти три подстроки по содержанию разные, но исходно разделяют один общий буфер.
		3.	unique.Make для каждой из трёх подстрок выделит свой собственный отдельный массив байт и скопирует в него символы.
		4.	В итоге у вас вместо одного большого буфера будет три маленьких, в каждом из которых своя копия.

	Это бывает невыгодно, когда вы хотите просто убрать дубликаты и при этом максимально экономить память. В таких случаях лучше не копировать байты заново, а просто переназначить все дубликаты на одну и ту же строку-экземпляр.
*/

func CanonicalizeString(s string) string {
	return unique.Make(s).Value()
}

// The way is more flexible. Just apply the above CanonicalizeString function
// for every string at run time, then all equal strings will share the same underlying
// bytes memory blocks.
func CanonicalizeStringsModern(ss []string) {
	for i, s := range ss {
		ss[i] = CanonicalizeString(s)
	}
}
