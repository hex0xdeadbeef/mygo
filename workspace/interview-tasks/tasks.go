package tasks

/*
	TASK 0 (Base for all others)

	Напишите функцию getGenerator, которая принимает на вход количество элементов и отдает канал out.
	В канал out нужно записать cnt рандомных чисел.

	getGenerator(cnt int) <-chan int {...}

*/

/*
	TASK 1

	Необходимо прочитать числа из канала from и передать в канал to.
	Функция Brigde должна учитывать контекст в процессе передачи данных из одного канала в другой.

	func Bridge(ctx context.Context, from <-chan int, to chan<- int) {...}

*/

/*
	TASK 2

	Функция FanOut принимает на вход канал src, откуда берет числа, а также каналы dests, куда прокидывает полученные из src числа.
	Задача: читать числа из канала src и записывать в каждый из dests каналов.
	Необходимо учитывать контекст, который указан в сигнатуре функции, для своевременной отмены процесса.

	FanOut(ctx context.Context, src <-chan int, dests ...customCh) {...}

*/

/*
	TASK 3

	Используя генератор чисел, обработать каждый сгерерированный элемент по принипу конвейра:
    1. Прочитать значение из канала генератора, возвести в квадрат и передать в канал squarer
	2. Прочитать значение из канала squarer и разделить на переданный коэффициент, передать в канал divider
	3. Прочитать значение из канала divider и передать в канал filterer только те элементы, которые не деляться на coef без остатка.
	4. Запринтить значения из канала filterer.

							Generator
								↓
	func squarer(ctx context.Context, in <-chan int) <-chan int {...}
								↓
	divider(ctx context.Context, coef int, in <-chan int) <-chan int {...}
								↓
	func filterer(ctx context.Context, coef int, in <-chan int) <-chan int {...}
								↓
							Print values
*/
