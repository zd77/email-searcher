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
	emails := []Email{}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v", dirPath)
		panic(err)
	}
	for _, entry := range entries {
		if string(entry.Name()[0]) == "." {
			continue
		}
		if entry.IsDir() {
			emails = append(emails, processDirectory(filepath.Join(dirPath, entry.Name()))...)
		} else {
			email := processEmailFile(filepath.Join(dirPath, entry.Name()))
			emails = append(emails, email)
		}
	}
	return emails
}

func processEmailFile(filePath string) Email {
	file, _ := os.Open(filePath)
	lines := bufio.NewScanner(file)
	email := extractEmailInfo(lines)
	return email
}

func bulkedIndexData(emails []Email) {
	user := "admin"
	password := "Complexpass#123"
	data := Data{
		Index:   "enron_emails",
		Records: emails,
	}
	jsonData, _ := json.MarshalIndent(data, "", " ")
	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func main() {
	cpu, err := os.Create("go.performance")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()

	startTime := time.Now()
	var pathFlag string
	flag.StringVar(&pathFlag, "path", ".", "Directory route")
	flag.Parse()
	directory, err := filepath.Abs(pathFlag)
	if err != nil {
		fmt.Println("Error getting the absolute path")
		return
	}
	fmt.Println("->", directory)
	emails := processDirectory(directory)
	bulkedIndexData(emails)
	endTime := time.Now()
	fmt.Printf("emails[%d] - duration[%s]", len(emails), endTime.Sub(startTime))
}
