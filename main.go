package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hinakaze/iniparser"
)

var (
	TotalFile   int = 0
	TotalLine   int = 0
	SpaceLine   int = 0
	CommentLine int = 0
	CodeLine    int = 0
)

type config struct {
	RootPath   string
	FileSuffix string
}

func main() {
	config, err := getConfig("./config.ini")
	if err != nil {
		panic("Parse config file failed ," + err.Error())
	}
	log.Printf("Start count line,RootPath [%s] FileSuffix [%s]", config.RootPath, config.FileSuffix)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		//		log.Printf("Find file [%s] ", path)
		if !strings.HasSuffix(info.Name(), config.FileSuffix) {
			return nil
		}

		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return err
		}
		breader := bufio.NewReader(file)

		totalLine, spaceLine, commentLine, codeLine := 0, 0, 0, 0

		linestr := ""
		for err := error(nil); err != io.EOF; linestr, err = breader.ReadString('\n') {
			totalLine++
			linestr = strings.TrimSpace(linestr)
			if linestr == "" {
				spaceLine++
			}
			if strings.HasPrefix(linestr, "//") || strings.HasPrefix(linestr, "/*") || strings.HasSuffix(linestr, "*/") {
				commentLine++
			}
			codeLine = totalLine - spaceLine - commentLine
		}
		log.Printf("File [%s] totalLine [%d] spaceLine [%d] commentLine[%d] codeLine [%d]", info.Name(), totalLine, spaceLine, commentLine, codeLine)

		//add to global variable
		TotalLine += totalLine
		SpaceLine += spaceLine
		CommentLine += commentLine
		CodeLine += codeLine
		TotalFile++
		return nil
	}
	err = filepath.Walk(config.RootPath, walkFunc)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Result : TotalFile [%d] TotalLine [%d] SpaceLine [%d] CommentLine [%d] CodeLine [%d] EffecientRate [%f]", TotalFile, TotalLine, SpaceLine, CommentLine, CodeLine, float64(CodeLine)/float64(TotalLine))
}

func getConfig(path string) (c config, err error) {

	iniparser.DefaultParse(path)
	section, ok := iniparser.GetSection("config")
	if !ok {
		return c, fmt.Errorf("ini section [config] not found")
	}
	c.RootPath, ok = section.GetValue("RootPath")
	if !ok {
		return c, fmt.Errorf("ini section [config] , key [RootPath] not found")
	}
	c.FileSuffix, ok = section.GetValue("FileSuffix")
	if !ok {
		return c, fmt.Errorf("ini section [config] , key [FileSuffix] not found")
	}
	return
}
