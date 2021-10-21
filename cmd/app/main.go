package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/ktnyt/go-moji"
)

const inFileName = "./data/bankCodeIn.csv"
const outFileName = "./out/bankCodeOut.json"

type BankCode string
type BranchCode string

func BankCodeFrom(str string) BankCode {
	num, _ := strconv.Atoi(str)
	return BankCode(fmt.Sprintf("%04d", num))
}

func BranchCodeFrom(str string) BranchCode {
	num, _ := strconv.Atoi(str)
	return BranchCode(fmt.Sprintf("%03d", num))
}

type Bank struct {
	Name          string   `json:"name"`
	Code          BankCode `json:"code"`
	Hiragana      string   `json:"hiragana"`
	HalfWidthKana string   `json:"halfWidthKana"`
	FullWidthKana string   `json:"fullWidthKana"`
	Branches      []Branch `json:"branches"`
}

type Branch struct {
	Name          string     `json:"name"`
	Code          BranchCode `json:"code"`
	Hiragana      string     `json:"hiragana"`
	HalfWidthKana string     `json:"halfWidthKana"`
	FullWidthKana string     `json:"fullWidthKana"`
}

func main() {
	inCsv, err := os.Open(inFileName)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(inCsv)
	r.FieldsPerRecord = 7

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	banks := map[BankCode]Bank{}

	for i, record := range records {
		if i == 0 {
			continue
		}

		bankCode := BankCodeFrom(record[5])

		_, ok := banks[bankCode]
		if !ok {
			halfWidthKana := record[2]
			fullWidthKana := moji.Convert(halfWidthKana, moji.HK, moji.ZK)
			banks[bankCode] = Bank{
				Code:          bankCode,
				Name:          record[0],
				Hiragana:      moji.Convert(fullWidthKana, moji.KK, moji.HG),
				FullWidthKana: fullWidthKana,
				HalfWidthKana: halfWidthKana,
			}
		}

		bank := banks[bankCode]
		halfWidthKana := record[4]
		fullWidthKana := moji.Convert(halfWidthKana, moji.HK, moji.ZK)
		bank.Branches = append(bank.Branches, Branch{
			Code:          BranchCodeFrom(record[6]),
			Name:          record[3],
			Hiragana:      moji.Convert(fullWidthKana, moji.KK, moji.HG),
			FullWidthKana: fullWidthKana,
			HalfWidthKana: halfWidthKana,
		})
		banks[bankCode] = bank
	}

	bankArray := []Bank{}
	for _, v := range banks {
		bankArray = append(bankArray, v)
	}
	sort.Slice(bankArray, func(i, j int) bool { return bankArray[i].Code < bankArray[j].Code })

	// fmt.Printf("%#v", bankArray)
	b, _ := json.Marshal(&bankArray)

	var buf bytes.Buffer
	err = json.Indent(&buf, b, "", "  ")
	if err != nil {
		panic(err)
	}
	b2 := buf.Bytes()

	err = ioutil.WriteFile(outFileName, b2, 0666)
	if err != nil {
		panic(err)
	}

}
