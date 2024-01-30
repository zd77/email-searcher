package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	BodyMsg string `json:"bodyMsg"`
}

type Content struct {
	Message_ID                string `json:"messageId"`
	Date                      string `json:"date"`
	From                      string `json:"from"`
	To                        string `json:"to"`
	Subject                   string `json:"subject"`
	Cc                        string `json:"cc"`
	Mime_Version              string `json:"mime"`
	Content_Type              string `json:"contentTyoe"`
	Content_Transfer_Encoding string `json:"cte"`
	X_From                    string `json:"xf"`
	X_To                      string `json:"xt"`
	X_cc                      string `json:"xc"`
	X_bcc                     string `json:"xbcc"`
	X_Folder                  string `json:"xfol"`
	X_Origin                  string `json:"xo"`
	X_FileName                string `json:"xfn"`
	BodyMsg                   string `json:"bodyMsg"`
}

type Data struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}

const batchSize = 100
const maxConcurrentRequests = 10

func extractEmailInfo(lines *bufio.Scanner) Email {
	email := Content{}
	for lines.Scan() {
		if strings.Contains(lines.Text(), "Message-ID:") && email.Message_ID != " " {
			email.Message_ID = strings.TrimSpace(lines.Text() + " ")[11:]
		} else if strings.Contains(lines.Text(), "Date:") && email.Date != " " {
			email.Date = strings.TrimSpace(lines.Text() + " ")[5:]
		} else if strings.Contains(lines.Text(), "From:") && email.From != " " {
			modifiedTxt := lines.Text() + "                 "
			email.From = strings.TrimSpace(modifiedTxt[7:])
		} else if strings.Contains(lines.Text(), "To:") && email.To != " " {
			modifiedTxt := lines.Text() + "                 "
			email.To = strings.TrimSpace(modifiedTxt[5:])
		} else if strings.Contains(lines.Text(), "Subject:") && email.Subject != " " {
			email.Subject = strings.TrimSpace(lines.Text() + " ")[8:]
		} else if strings.Contains(lines.Text(), "Cc:") {
			email.Cc = strings.TrimSpace(lines.Text() + ".")[3:]
		} else if strings.Contains(lines.Text(), "Mime-Version:") {
			email.Mime_Version = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "Content-Type:") {
			email.Content_Type = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "Content-Transfer-Encoding:") {
			email.Content_Transfer_Encoding = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "X-From:") {
			email.X_From = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "X-To:") {
			email.X_To = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "X-cc:") {
			email.X_cc = strings.TrimSpace(lines.Text() + ".")[6:]
		} else if strings.Contains(lines.Text(), "X-bcc:") {
			email.X_bcc = strings.TrimSpace(lines.Text() + ".")[6:]
		} else if strings.Contains(lines.Text(), "X-Folder:") {
			email.X_Folder = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "X-Origin:") {
			email.X_Origin = strings.TrimSpace(lines.Text() + ".")[9:]
		} else if strings.Contains(lines.Text(), "X-FileName:") {
			email.X_FileName = strings.TrimSpace(lines.Text() + ".")[9:]
		} else {
			email.BodyMsg = email.BodyMsg + lines.Text()
		}
	}
	return Email{
		From:    email.From,
		To:      email.To,
		Subject: email.Subject,
		BodyMsg: email.BodyMsg,
	}
}

func processDirectory(dirPath string) []Email {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		emails []Email
	)
	processEntry := func(entry os.DirEntry) {
		defer wg.Done()
		if entry.IsDir() {
			subDirPath := filepath.Join(dirPath, entry.Name())
			subEmails := processDirectory(subDirPath)
			mu.Lock()
			emails = append(emails, subEmails...)
			mu.Unlock()
		} else {
			email := processEmailFile(filepath.Join(dirPath, entry.Name()))
			mu.Lock()
			emails = append(emails, email)
			mu.Unlock()
		}
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v", dirPath)
		panic(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		wg.Add(1)
		go processEntry(entry)
	}
	wg.Wait()
	return emails
}

func processEmailFile(filePath string) Email {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v", filePath, err)
		panic(err)
	}
	defer file.Close()
	lines := bufio.NewScanner(file)
	return extractEmailInfo(lines)
}

func bulkedIndexData(emails []Email, uri string) {
	user := "admin"
	password := "Complexpass#123"
	client := &http.Client{}
	var wg sync.WaitGroup

	for i := 0; i < len(emails); i += batchSize {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			if end > len(emails) {
				end = len(emails)
			}
			if start >= end {
				return
			}
			data := Data{
				Index:   "enron_emails",
				Records: emails[start:end],
			}
			var jsonData bytes.Buffer
			if err := json.NewEncoder(&jsonData).Encode(data); err != nil {
				fmt.Printf("Error encoding JSON: %v", err)
				panic(err)
			}
			req, err := http.NewRequest("POST", uri+"/api/_bulkv2", &jsonData)
			if err != nil {
				fmt.Printf("Error creating HTTP request: %v ", err)
				panic(err)
			}
			req.SetBasicAuth(user, password)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error making HTTP request: %v\n", err)
				panic(err)
			}
			defer resp.Body.Close()
		}(i, i+batchSize)
		if i%maxConcurrentRequests == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func envSelector(env string) (string, string) {
	if env != "LOCAL" {
		return "../../enron_mail_20110402/maildir2", "http://load-balancer-dev-2129612410.us-east-1.elb.amazonaws.com"
	}
	return "../../enron_mail_20110402/maildir", "http://localhost:4080"
}

func main() {
	cpu, err := os.Create("cpu.perf")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()

	startTimeAll := time.Now()
	var envFlag string
	flag.StringVar(&envFlag, "env", "LOCAL", "Directory route")
	flag.Parse()
	path, uri := envSelector(envFlag)
	fmt.Println("Environment -> ", envFlag)
	directory, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting the absolute path")
		return
	}
	startTimePD := time.Now()
	emails := processDirectory(directory)
	endTimePD := time.Now()
	fmt.Printf("emails[%d] - duration[%s]\n", len(emails), endTimePD.Sub(startTimePD))
	startTimeBID := time.Now()
	bulkedIndexData(emails, uri)
	endTimeBID := time.Now()
	fmt.Printf("Sending data to zinc [%s]\n", endTimeBID.Sub(startTimeBID))
	endTimeAll := time.Now()
	fmt.Printf("[%s]", endTimeAll.Sub(startTimeAll))
}
