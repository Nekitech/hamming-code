package main

import (
	"fmt"
	"reflect"
)

func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}
func main() {
	// Пример данных для кодирования
	data := []int{1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0}
	var controlBit []int

	for i := 0; i < len(data); i++ {
		if isPowerOfTwo(i + 1) {
			controlBit = append(controlBit, i)
		}
	}

	fmt.Println(controlBit)

	// Кодирование данных
	encodedData, err := encodeHamming74(controlBit, data)
	if err != nil {
		fmt.Println("Ошибка кодирования:", err)
		return
	}

	fmt.Println("Исходные данные:", data)
	fmt.Println("Закодированные данные:", encodedData)

	// Имитация ошибок в переданных данных (поменяем один бит)
	encodedData[1] = 1
	fmt.Println("Данные с ошибкой:", encodedData)

	// Декодирование данных с возможностью исправления ошибок
	decodedData, err := decodeHamming74(controlBit, encodedData, len(data))
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

// encodeHamming74 кодирует входные данные Хемминг-кодом (7, 4)
func encodeHamming74(controlBit []int, data []int) ([]int, error) {
	// Создаем закодированный массив
	encodedData := make([]int, len(data)+len(controlBit))

	count := 0
	for i := 0; i < len(encodedData); i++ {
		if isPowerOfTwo(i + 1) {
			continue
		} else {
			//fmt.Println(i)
			encodedData[i] = data[count]
			count++
		}
	}
	//fmt.Println(encodedData, "!!!!")

	for i := 0; i < len(encodedData); i++ {
		if isPowerOfTwo(i + 1) {
			encodedData[i] = calculateParityBit(encodedData, getCheckIndices(controlBit, i+1, len(data)))
			//fmt.Println(getCheckIndices(controlBit, i+1))
		} else {
			continue
		}
	}

	return encodedData, nil
}

// calculateParityBit рассчитывает значение контрольного бита
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

// decodeHamming74 декодирует данные Хемминг-кода (7, 4) и исправляет ошибки
func decodeHamming74(parityBitIndices []int, encodedData []int, countData int) ([]int, error) {
	decodedData := make([]int, countData)

	fmt.Println(parityBitIndices)

	// Проверяем контрольные биты
	errorPosition := 0
	for _, index := range parityBitIndices {
		calculatedParityBit := calculateParityBit(encodedData, getCheckIndices(parityBitIndices, index+1, countData))
		if calculatedParityBit != encodedData[index] {
			errorPosition += index + 1
		}
	}

	// Если есть ошибка, исправляем
	if errorPosition > 0 {
		fmt.Println("Исправлен бит на позиции:", errorPosition-1, "::|::", encodedData[errorPosition-1], " -> ", encodedData[errorPosition-1]^1)
		encodedData[errorPosition-1] ^= 1
	}

	count := 0
	for i := 0; i < len(encodedData); i++ {
		if isPowerOfTwo(i + 1) {
			continue
		} else {
			fmt.Println(count, i+1, encodedData[i], encodedData)
			decodedData[count] = encodedData[i]
			count++
		}
	}
	return decodedData, nil
}

// getCheckIndices возвращает индексы битов, используемых для проверки контрольного бита
func getCheckIndices(parityBitIndices []int, parityBitIndex int, countData int) []int {
	indices := make([]int, 0)

	i := parityBitIndex
	step := i
	for i <= countData {
		if step != 0 {
			if !implContains(parityBitIndices, i-1) {
				indices = append(indices, i-1)
			}
			step -= 1
			i += 1
		} else {
			i += parityBitIndex
			step = parityBitIndex
		}

	}
	return indices
}

func implContains(sl []int, number int) bool {
	for _, value := range sl {
		if value == number {
			return true
		}
	}
	return false
}
