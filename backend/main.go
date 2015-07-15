package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

const EMAIL_PERIOD = time.Second * 1800
const EMAIL_WAIT = time.Second * 30

var emailList = []string{"antoine.pourchet@gmail.com"}
var lastCount = -1
var lastImage = []byte{}
var fsm *FSM
var emailThreshold = time.Now()

func removeEmail(toremove string) []string {
	newEmailList := []string{}
	for _, email := range emailList {
		if email != toremove {
			newEmailList = append(newEmailList, email)
		}
	}
	return newEmailList
}

func cleanBody(body []byte) string {
	return strings.Replace(string(body), "\n", "", -1)
}

func sendEmail(count int) {
	fmt.Println("Sending the email")

	subjectString := "ALERT: INVENTORY LOW"
	bodyString := ""
	if count == 0 {
		bodyString = "Out of items, please refill"
	} else if count == 1 {
		bodyString = "You only have 1 item left, order some moar"
	} else if count < 3 {
		bodyString = "Running low, be aware"
	}

	emailUser := &EmailUser{"monitor.inventory", "squaresquaresquare1!", "smtp.gmail.com", 587}
	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer)

	emailBody := fmt.Sprintf("From: AlwaysBeer\nTo: Dear Customer\nSubject: %s\n\n%s\nBottles Left: %d\n", subjectString, bodyString, count)
	err := smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), auth,
		emailUser.Username,
		emailList,
		[]byte(emailBody))
	if err != nil {
		fmt.Println("ERROR: attempting to send a mail ", err)
	}
}

func emailSender() {
	ticker := time.NewTicker(EMAIL_PERIOD)
	for _ = range ticker.C {
		if lastCount < 3 && lastCount > 0 {
			go sendEmail(lastCount)
		}
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: statusHandler")
	fmt.Fprintf(w, "ok")
}

func addEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: addEmailHandler")
	email := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		email = cleanBody(body)
	} else {
		email = r.URL.Query().Get("email")
	}
	fmt.Printf("Email: '%s'\n", email)
	if len(email) != 0 {
		emailList = append(emailList, email)
	}
	fmt.Printf("Email List: %v\n\n", emailList)
	fmt.Fprintf(w, "Email added: %s", email)
}

func lastCountHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: lastCountHandler")
	if lastCount < 0 {
		fmt.Fprintf(w, "We do not know how many bottles are left!")
	} else {
		fmt.Fprintf(w, "Last Count: %d\n", lastCount)
	}
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: countHandler")
	count := 0
	if len(r.URL.RawQuery) == 0 {
		bodyArr, _ := ioutil.ReadAll(r.Body)
		body := cleanBody(bodyArr)
		count, _ = strconv.Atoi(body)
	} else {
		countStr := r.URL.Query().Get("count")
		count, _ = strconv.Atoi(countStr)
	}

	fmt.Printf("Got a bottle count: %d\n", count)
	fmt.Fprintf(w, "Bottle count: '%d'\n", count)
	lastCount = count
}

func removeEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: removeEmailHandler")
	email := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		email = cleanBody(body)
	} else {
		email = r.URL.Query().Get("email")
	}
	fmt.Printf("Email: '%s'\n", email)
	if len(email) > 0 {
		emailList = removeEmail(email)
	}
	fmt.Printf("Email List: %v\n\n", emailList)
	fmt.Fprintf(w, "Email removed: %s", email)
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: sendEmailHandler")
	sendEmail(lastCount)
}

func setImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: setImageHandler")
	lastImage, _ = ioutil.ReadAll(r.Body)
	fmt.Println("Set the last image.")
	go handleInput()
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: getImageHandler")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(lastImage)))
	w.Write(lastImage)
}

func handleHandlers() {
	http.HandleFunc("/_status", statusHandler)
	http.HandleFunc("/setcount", countHandler)
	http.HandleFunc("/lastcount", lastCountHandler)
	http.HandleFunc("/addemail", addEmailHandler)
	http.HandleFunc("/removeemail", removeEmailHandler)
	http.HandleFunc("/sendemail", sendEmailHandler)
	http.HandleFunc("/setimage", setImageHandler)
	http.HandleFunc("/getimage", getImageHandler)
}

func handleInput() {
	// Save lastImage into file
	ioutil.WriteFile("/tmp/lastpicture.jpeg", lastImage, 0644)
	// Shellout to classifier
	out, err := exec.Command("/root/alwaysbeer/src/test_bottle_detector", "/root/alwaysbeer/src/single/test_images.xml", "silent").Output()
	// Get output and make to integer
	if err != nil {
		fmt.Println(err)
		return
	}
	cleanOut := cleanBody(out)
	fmt.Printf("Output: %s\n", cleanOut)
	lastCount, _ = strconv.Atoi(cleanOut)
	fsm.Transition(Input(lastCount))
}

func startFSM() {
	fsm = NewFSM()

	fsm.AddState(0, func(input Input) int {
		if input == 0 {
			return 2
		}
		return 1
	})
	fsm.AddState(1, func(input Input) int {
		if input != 0 {
			return 1
		}
		return 2
	})
	fsm.AddState(2, func(input Input) int {
		if input != 0 {
			return 1
		}
		return 3
	})
	fsm.AddState(3, func(input Input) int {
		if input != 0 {
			return 1
		}
		return 4
	})
	fsm.AddState(4, func(input Input) int {
		if input != 0 {
			return 2
		}
		return 5
	})
	fsm.AddState(5, func(input Input) int {
		if input != 0 {
			return 3
		}
		sendEmail(lastCount)
		return 6
	})
	fsm.AddState(6, func(input Input) int {
		if input != 0 {
			return 1
		}
		if (time.Now().After(emailThreshold)) {
			sendEmail(lastCount)
			emailThreshold = time.Now().Add(EMAIL_WAIT)
		}
		return 6
	})
}

func main() {
	go emailSender()
	handleHandlers()
	startFSM()
	http.ListenAndServe(":8080", nil)
}
