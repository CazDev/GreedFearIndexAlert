package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"
)

func main() {
	for {
		log.Println("Starting...")
		fgi, err := getFGI()
		if err == nil && fgi != 999 {
			// success
			msg := "Todays FGI is: " + fmt.Sprint(fgi)
			log.Println(msg)
			if fgi < 40 {
				// send mail notification
				log.Println("Sending email notification...")
				sendMail(msg, "Be fearful when others are greedy, and be greedy when others are fearful.")
			}
		} else {
			sendMail("FGI Index error", "There has been an error getting the FGI index, please see console log for more details")
		}
		// Wait 24hrs before starting again
		log.Println("Waiting 24hrs...")
		time.Sleep(24 * time.Hour)
	}
}

func sendMail(subject string, body string) error {

	// Setup a gmail account to send notifications
	// Make sure the new gmail account has "Less secure apps" turned on in settings

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "from@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "to@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "from@gmail.com", "<your-gmail-password>")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Sent!")
		return nil
	}
}

func getFGI() (int, error) {
	// API setup
	url := "https://fear-and-greed-index.p.rapidapi.com/v1/fgi"
	req, _ := http.NewRequest("GET", url, nil)

	// API key and host
	// get this from https://rapidapi.com/rpi4gx/api/fear-and-greed-index
	req.Header.Add("x-rapidapi-key", "<your-api-key>")
	req.Header.Add("x-rapidapi-host", "fear-and-greed-index.p.rapidapi.com")

	// response from API
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// API sends 200 OK
	if res.StatusCode == 200 {
		str := string(body)[23:26]
		re := regexp.MustCompile("[0-9]+")
		fgiArr := re.FindAllString(str, 1)
		fgi, err := strconv.Atoi(fgiArr[0])

		// handle regex errors
		if err != nil {
			log.Println(err)
			return 999, err
		}

		// success
		return fgi, nil
	} else {
		// bad API response
		log.Println("err: API response status was not 200 OK")
		return 999, nil
	}
}
