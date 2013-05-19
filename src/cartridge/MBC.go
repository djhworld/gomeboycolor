package cartridge

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"types"
)

type MemoryBankController interface {
	Write(addr types.Word, value byte)
	Read(addr types.Word) byte
	SaveRam(filename string) error
	LoadRam(filename string) error
	switchROMBank(bank int)
	switchRAMBank(bank int)
}

func populateROMBanks(rom []byte, noOfBanks int) [][]byte {
	romBanks := make([][]byte, noOfBanks)

	//ROM Bank 0 and 1 are the same
	romBanks[0] = rom[0x4000:0x8000]
	var chunk int = 0x4000
	for i := 1; i < noOfBanks; i++ {
		romBanks[i] = rom[chunk : chunk+0x4000]
		chunk += 0x4000
	}

	return romBanks
}

func populateRAMBanks(noOfBanks int) [][]byte {
	ramBanks := make([][]byte, noOfBanks)

	for i := 0; i < noOfBanks; i++ {
		ramBanks[i] = make([]byte, 0x2000)
	}

	return ramBanks
}

func WriteRAMToDisk(path string, ramBanks [][]byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < len(ramBanks); i++ {
		writer.Write(ramBanks[i])
		/*
			var b bytes.Buffer
			w := zlib.NewWriter(&b)
			w.Write(ramBanks[i])
			w.Close()
			before := base64.StdEncoding.EncodeToString(b.Bytes())
			fmt.Println("Before = ", len(b.Bytes()))

			after,_ := base64.StdEncoding.DecodeString(before)
			b = *bytes.NewBuffer(after)
			r, err := zlib.NewReader(&b)
			if err != nil {
				panic(fmt.Sprintln("Error" , err))
			}
			out := bytes.NewBuffer(make([]byte,0))
			io.Copy(out, r)
			fmt.Println("After = ", len(out.Bytes()))
			fmt.Println(bytes.Compare(ramBanks[i], out.Bytes()))
			r.Close()
		*/

	}
	writer.Flush()
	return nil
}

func ReadRAMFromDisk(path string, chunkSize int, expectedSize int) ([][]byte, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(fileBytes) != expectedSize {
		return nil, errors.New(fmt.Sprintf("RAM file is not %d bytes!", expectedSize))
	}

	var chunk int = 0x0000
	var noOfBanks int = expectedSize / chunkSize
	var ramBanks [][]byte = make([][]byte, noOfBanks)

	for i := 0; i < noOfBanks; i++ {
		ramBanks[i] = fileBytes[chunk : chunk+chunkSize]
		chunk += chunkSize
	}

	return ramBanks, nil
}
