/*1. ���ĸ�Э��1��2��3��4, Э��1�Ĺ��ܾ������1��Э��2�Ĺ��ܾ������2���Դ����ơ����������ĸ��ļ�A��B��C��D����ʼ��Ϊ�գ����ʵ�������ĸ��ļ��������¸�ʽ��
A��1 2 3 4 1 2...
B��2 3 4 1 2 3...
C��3 4 1 2 3 4...
D��4 1 2 3 4 1...*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var ch1, ch2, ch3, ch4 chan int

func main() {
	ch1 = make(chan int)
	ch2 = make(chan int)
	ch3 = make(chan int)
	ch4 = make(chan int)
	go Ch1()
	go Ch2()
	go Ch3()
	go Ch4()
	magicNumber('A', 10)
	magicNumber('B', 10)
	magicNumber('C', 10)
	magicNumber('D', 10)
}

func magicNumber(fliename byte, round int) {
	var temp int = int(fliename - 65)
	for i := temp; i < temp+round; i++ {
		num := i%4 + 1
		switch {
		case num == 1:
			fliewrite(string(fliename), <-ch1)
		case num == 2:
			fliewrite(string(fliename), <-ch2)
		case num == 3:
			fliewrite(string(fliename), <-ch3)
		case num == 4:
			fliewrite(string(fliename), <-ch4)
		}
	}
}
func Ch1() {
	for {
		ch1 <- 1
	}
}

func Ch2() {
	for {
		ch2 <- 2
	}
}

func Ch3() {
	for {
		ch3 <- 3
	}
}

func Ch4() {
	for {
		ch4 <- 4
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func fliewrite(fn string, num int) {
	var filename = "./" + fn + ".txt"
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //����ļ�����
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //���ļ�
		fmt.Println("�ļ�����")
	} else {
		f, err1 = os.Create(filename) //�����ļ�
		fmt.Println("�ļ�������")
	}
	check(err1)
	w := bufio.NewWriter(f) //�����µ� Writer ����
	w.WriteString(strconv.Itoa(num))
	//fmt.Printf("д�� %d ���ֽ�n", n4)
	w.Flush()
	f.Close()
}
