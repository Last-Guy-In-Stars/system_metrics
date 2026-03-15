package main

import (
	"fmt"
	// "net/smtp"
	"bytes"
	"log"
	"os"
	"os/exec"
	"time"
)

func scheduleEmail() {
	for {
		now := time.Now()

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			17, 30, 0, 0,
			now.Location(),
		)

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		time.Sleep(time.Until(next))

		file := fmt.Sprintf(
			"metrics_%s.csv",
			time.Now().Format("2006-01-02"),
		)

		sendMetricsEmail(
			"kartamonov1@yandex.ru",
			"ЦУП - метрики",
			file,
		)
	}
}

func sendMetricsEmail(to, subject, csvFile string) {
	data, err := os.ReadFile(csvFile)
	if err != nil {
		log.Println("Error reading CSV:", err)
		return
	}

	htmlTable := csvToHTML(string(data))

	err = sendMail(to, subject, htmlTable)
	if err != nil {
		log.Println("Mail error:", err)
	} else {
		log.Println("Daily metrics mail sent")
	}
}

func csvToHTML(csvData string) string {
	lines := bytes.Split([]byte(csvData), []byte("\n"))
	var buf bytes.Buffer
	buf.WriteString("<table border='1' cellpadding='5' cellspacing='0'>\n")
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		buf.WriteString("<tr>")
		cols := bytes.Split(line, []byte(","))
		for _, col := range cols {
			if i == 0 {
				buf.WriteString(fmt.Sprintf("<th>%s</th>", col))
			} else {
				buf.WriteString(fmt.Sprintf("<td>%s</td>", col))
			}
		}
		buf.WriteString("</tr>\n")
	}
	buf.WriteString("</table>")
	return buf.String()
}

func sendMail(to, subject, body string) error {

	cmd := exec.Command("mail", "-s", subject,
		"-a", "Content-Type: text/html; charset=UTF-8",
		"-r", "cup_system@n29.ru",
		to)
	cmd.Stdin = bytes.NewBufferString(body)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error send mail: %v, stderr: %s", err)
	}

	return nil
}
