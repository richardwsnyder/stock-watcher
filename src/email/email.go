package email

import (
	"log"
	"net/smtp"

	e "finnhub/src/env"
)

func PriceTargetMet(body string) {
	from := e.GoDotEnvVariable("EUSER")
	pass := e.GoDotEnvVariable("EPASS")
	to := e.GoDotEnvVariable("TO")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Stock Target Met\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s\n", err)
		return
	}

	log.Print("sent")
}
