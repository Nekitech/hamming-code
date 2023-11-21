package main

import (
	"fmt"
	"math/rand"
	"reflect"
)

// 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0,
func main() {
	// Пример данных для кодирования
	data := []int{0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1}
	var controlBit []int
	// 1 0 0 1 1 0 0 0 0 1 0 0 0 0 1 0 1 1 1 0 1
	for i := 0; i < len(data); i++ {
		if isPowerOfTwo(i + 1) {
			controlBit = append(controlBit, i)
		}
	}

	fmt.Println(controlBit)

	// Кодирование данных
	encodedData, control_sum := encodeHamming(controlBit, data)

	fmt.Println("контрольная сумма - ", control_sum)

	fmt.Println("Исходные данные:", data)
	fmt.Println("Закодированные данные:", encodedData)
	fmt.Println("Длина закодированных данных", len(encodedData))

	// Имитация ошибок в переданных данных (поменяем один бит)

	//encodedData[6] ^= 1
	randomPlaceError(encodedData)
	randomPlaceError(encodedData)
	//encodedData[2] = 1
	//encodedData[3] = 0
	////encodedData[6] = 1

	fmt.Println("Данные с ошибкой:", encodedData)

	// Декодирование данных с возможностью исправления ошибок
	decodedData, err := decodeHamming(controlBit, encodedData, control_sum)
	if err != nil {
		fmt.Println("Ошибка декодирования:", err)
	} else {
		fmt.Println("Декодированные данные:", decodedData)
		if reflect.DeepEqual(data, decodedData) {
			fmt.Println("Массивы равны")
		} else {
			fmt.Println("Массивы не равны")
		}

	}
}

/*
Кодирует входные данные Хемминг-кодом с n-ым количеством символов и контрольных битов
*/
func encodeHamming(controlBit []int, data []int) ([]int, []int) {
	// Создаем закодированный массив
	encodedData := make([]int, len(data)+len(controlBit))

	count := 0
	for i := 0; i < len(encodedData); i++ {
		if implContains(controlBit, i) {
			continue
		}
		encodedData[i] = data[count]
		count++

	}

	for _, index := range controlBit {
		indexes := getCheckIndices(index+1, len(encodedData))
		encodedData[index] = calculateParityBit(encodedData, indexes)

	}

	controlsSum := make([]int, len(encodedData))

	for _, index := range controlBit {
		indexes := getCheckIndices(index+1, len(encodedData))
		controlsSum[index] = calculateParityBit(encodedData, indexes)
		fmt.Println(controlsSum[index])

	}

	return encodedData, controlsSum
}

// decodeHamming декодирует данные Хемминг-кода (7, 4) и исправляет ошибки
func decodeHamming(controlBit []int, encodedData []int, control_sum []int) ([]int, error) {
	decodedData := make([]int, len(encodedData)-len(controlBit))
	// Проверяем контрольные биты
	errorPosition := 0
	for _, index := range controlBit {
		indexes := getCheckIndices(index+1, len(encodedData))
		calculatedParityBit := calculateParityBit(encodedData, indexes)

		if calculatedParityBit != control_sum[index] {
			errorPosition += index + 1
		}
	}
	fmt.Println(errorPosition)

	// Если есть ошибка, исправляем
	if errorPosition > 0 {
		fmt.Println("Исправлен бит на позиции:", errorPosition-1, "::|::", encodedData[errorPosition-1], " -> ", encodedData[errorPosition-1]^1)
		encodedData[errorPosition-1] ^= 1
	}

	count := 0
	//fmt.Println(parityBitIndices, len(encodedData))
	for i := 0; i < len(encodedData); i++ {
		if implContains(controlBit, i) {
			continue
		}
		decodedData[count] = encodedData[i]
		count++

	}

	return decodedData, nil
}

/*
Рассчитывает значение контрольного бита (просто вычисляет четность суммы и в зависимости от нее выдает либо 0, либо 1)
*/
func calculateParityBit(data []int, indices []int) int {
	parityBit := 0
	for _, index := range indices {
		if data[index] == 1 {
			parityBit += data[index]
		}
	}

	if parityBit%2 == 0 {
		return 0
	} else {
		return 1
	}
}

/*
Возвращает индексы битов, используемых для проверки контрольного бита (индексы, где должны стоять те биты, которые потом будем суммировать, и смотреть на четность/нечетность суммы 1-ек и делать выводы,
была ошибка или нет
*/
func getCheckIndices(parityBitIndex int, countData int) []int {
	indices := make([]int, 0)

	// устанавливаем контрольный индекс, который будет соотв. шагу
	step := parityBitIndex
	// переменная, которая отвечает за подсчет интервала (например, при шаге в 2 в интервале должно быть пройдено 2 значения)
	i := parityBitIndex
	for step <= countData {
		if i != 0 {
			indices = append(indices, step-1)
			// нашли индекс, уменьшили переменную, двигаемся дальше
			i -= 1
			step += 1
		} else {
			// текущий интервал пройден - увеличиваем шаг на индекс контрольного бита
			step += parityBitIndex
			// обновляем значение интервала
			i = parityBitIndex
		}

	}
	//fmt.Println(indices, parityBitIndex)
	return indices
}

/*
Функция, проверяющая на вхождение в массив значения
*/
func implContains(sl []int, number int) bool {
	for _, value := range sl {
		if value == number {
			return true
		}
	}
	return false
}

/*
Функция возведения в n-ую степень
*/
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

func randomPlaceError(arr []int) {
	index := getRandomIndex(arr)
	fmt.Println("Индекс с ошибкой -> ", index)
	arr[index] ^= 1
}

func getRandomIndex(arr []int) int {
	if len(arr) == 0 {
		return -1 // возврат -1, если массив пуст
	}

	return rand.Intn(len(arr))
}
