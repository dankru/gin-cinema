package main

import (
	"encoding/json"
	"os"
)

var cfg Setting

func init() {
	// Открытие файла конфигурации
	file, e := os.Open("settings.cfg")
	if e != nil {
		panic(e.Error())
	}
	defer file.Close()

	stat, e := file.Stat()
	if e != nil {
		panic(e.Error())
	}
	
	readByte := make([]byte, stat.Size())
	
	// file.Read считывает содержимое файлы в массив readByte
	_, e = file.Read(readByte)
	if e != nil{
		panic(e.Error())
	}

	// Unmarshal делает из объекта json объект нашей структуры. Первым аргументом мы передаём массив, а вторым интерфейс, которым будет являться 
	//созданная ранее переменная cfg с типом Setting
	e = json.Unmarshal(readByte, &cfg)
	if e != nil{
		panic(e.Error())
	}
	//Теперь у нас есть строка с настройками в виде объекта setting
	//fmt.Println(cfg)
}