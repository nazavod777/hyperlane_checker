package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"main/internal/accountsParser"
	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/util"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func inputUser(inputText string) string {
	if inputText != "" {
		fmt.Print(inputText)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("Unexpected Error: %v", r)
		fmt.Println("Press Enter to Exit..")
		_, err := fmt.Scanln()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func processAccounts(threads int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, threads)

	for _, account := range global.AccountsList {
		wg.Add(1)
		sem <- struct{}{}

		go func(acc types.AccountData) {
			defer wg.Done()
			accountsParser.ParseAccount(acc)
			global.CurrentProgress += 1
			<-sem
		}(account)
	}

	wg.Wait()
}

func main() {
	var inputUserData string
	// <-- init

	// init log
	initLog()

	wr, err := os.OpenFile(filepath.Join("log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Panicf("Error When Opening Log File: %s", err)
	}

	defer func(wr *os.File) {
		err = wr.Close()
		if err != nil {
			log.Panicf("Error When Closing Log File: %s", err)
		}
	}(wr)
	mw := io.MultiWriter(os.Stdout, wr)
	log.SetOutput(mw)

	// handle panic
	defer handlePanic()

	// init proxies
	err = util.InitProxies(filepath.Join("config", "proxies.txt"))
	if err != nil {
		log.Panicf("Error initializing proxies: %s", err)
	}
	// --> init

	if len(util.Proxies) <= 0 {
		global.Clients = append(global.Clients, util.CreateClient(""))
	} else {
		for _, proxy := range util.Proxies {
			global.Clients = append(global.Clients, util.CreateClient(proxy))
		}
	}

	accountsListString, err := util.ReadFileByRows(filepath.Join("config", "accounts.txt"))

	if err != nil {
		log.Panicf(err.Error())
	}

	global.AccountsList, err = util.GetAccounts(accountsListString, false)

	if err != nil {
		log.Panicf(err.Error())
	}

	log.Printf("Successfully Loaded %d Accounts", len(global.AccountsList))

	inputUserData = inputUser("\nThreads: ")
	threads, err := strconv.Atoi(inputUserData)

	if err != nil {
		log.Panicf("Wrong Threads Number: %s", inputUserData)
	}

	fmt.Printf("\n")
	global.TargetProgress = len(accountsListString)

	processAccounts(threads)

	log.Printf("The Work Has Been Successfully Finished")
	inputUser("\n\nPress Enter to Exit..")
}
