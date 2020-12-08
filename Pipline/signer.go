package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const GorNum = 6

func ExecutePipeline(freeFlowJobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, val := range freeFlowJobs {
		wg.Add(1)
		out := make(chan interface{})
		go jobWorker(val, in, out, wg)
		in = out
	}
	wg.Wait()
}

func jobWorker(freeFlowJobs job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	freeFlowJobs(in, out)
}

func SingleHash(in, out chan interface{}) {

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go singleHashWorker(i, out, wg, mu)
	}

	wg.Wait()
}

func singleHashWorker(in interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	data := strconv.Itoa(in.(int))

	mu.Lock()
	md5Data := DataSignerMd5(data)
	mu.Unlock()

	crc32DataChan := make(chan string)
	go func() {
		crc32DataChan <- DataSignerCrc32(data)
	}()

	crc32Md5Data := DataSignerCrc32(md5Data)
	crc32Data := <-crc32DataChan
	out <- crc32Data + "~" + crc32Md5Data
}



func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for i := range in {
		wg.Add(1)
		go multiHashWorker(i.(string), out, wg)
	}

	wg.Wait()
}

func multiHashWorker(in string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	mu := &sync.Mutex{}
	wg1 := &sync.WaitGroup{}

	concatArray := make([]string, GorNum)

	for i := 0; i < GorNum; i++ {
		wg1.Add(1)
		data := strconv.Itoa(i) + in
		go func(concatArray []string, data string, index int, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()

			data = DataSignerCrc32(data)

			mu.Lock()
			concatArray[index] = data
			mu.Unlock()

		}(concatArray, data, i, wg1, mu)
	}

	wg1.Wait()
	result := strings.Join(concatArray, "")

	out <- result
}

func CombineResults(in, out chan interface{}) {
	var array []string

	for i := range in {
		array = append(array, i.(string))
	}

	sort.Strings(array)
	result := strings.Join(array, "_")
	out <- result
}
