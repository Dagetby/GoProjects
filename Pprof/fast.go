package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mailru/easyjson"
	"hw3_bench/model"
	"io"
	"os"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции


//var dataPool = sync.Pool{
//	New: func() interface{}{
//		return bytes.NewBuffer(make([]byte,0,64))
//	},
//}


var dataPool = sync.Pool{
		New: func() interface{}{
			return &model.Data{}
		},
	}

func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileContents := bufio.NewScanner(file)
	seenBrowsersM := make(map[string]struct{})
	//foundUsers := ""
	i := 0


	//var dates []model.Data
	fmt.Fprintf(out, "found users:\n")

	for fileContents.Scan(){

		i++
		line := fileContents.Bytes()
		if (!bytes.Contains(line,[]byte("Android"))) && (!bytes.Contains(line,[]byte("MSIE"))){
			continue
		}

		isAndroid := false
		isMSIE := false

		d := dataPool.Get().(*model.Data)


		err := easyjson.Unmarshal(line, d)
		if err != nil {
			panic(err)
		}

		for _, browser := range d.Browsers {

			if strings.Contains(browser,"MSIE"){
				isMSIE = true

				if _, ok := seenBrowsersM[browser]; !ok {
					seenBrowsersM[browser] = struct{}{}
				}

			}else if strings.Contains(browser,"Android"){
				isAndroid = true

				if _, ok := seenBrowsersM[browser]; !ok {
					seenBrowsersM[browser] = struct{}{}
				}
			}
		}
		dataPool.Put(d)

		if !(isAndroid && isMSIE){
			continue
		}

		d.Email = strings.ReplaceAll(d.Email, "@", " [at] ")
		fmt.Fprintf(out,"[%d] %s <%s>\n", i-1, d.Name, d.Email)
	}

	fmt.Fprintf(out, "\nTotal unique browsers %d\n", len(seenBrowsersM))

}
