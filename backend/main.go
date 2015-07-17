package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
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

var slackHooks = []string{"https://hooks.slack.com/services/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"}
var emailList = []string{"antoine.pourchet@gmail.com"}
var requestList = []string{}
var lastCount = -1
var slackToken = ""
var notificationsEnabled = true
var lastImage = []byte{}
var fsm *FSM
var emailThreshold = time.Now()

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
		if time.Now().After(emailThreshold) {
			sendEmail(lastCount)
			emailThreshold = time.Now().Add(EMAIL_WAIT)
		}
		return 6
	})
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

func sendSlackMessage(msg string) {
	if notificationsEnabled == false {
		return
	}
	message := fmt.Sprintf("{\"text\": \"%s\"}", msg)
	for _, hook := range slackHooks {
		resp, err := http.Post(hook, "text", bytes.NewReader([]byte(message)))
		if err != nil {
			fmt.Println("failed with error: ", err)
			return
		}
		if resp.StatusCode != 200 {
			fmt.Println("failed with code: ", resp.StatusCode)
		}
	}
}

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
	if notificationsEnabled == false {
		return
	}
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

	go sendSlackMessage(bodyString)
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
	} else if lastCount == 0 {
		fmt.Fprintf(w, "Last Count: %d\n", lastCount)
	} else {
		fmt.Fprintf(w, "There are at least %d bottles left!\n", lastCount)
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

func addSlackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: addSlackHandler")
	slack := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		slack = cleanBody(body)
	} else {
		slack = r.URL.Query().Get("slack")
	}
	fmt.Printf("slack: '%s'\n", slack)
	if len(slack) != 0 {
		slackHooks = append(slackHooks, slack)
	}
	fmt.Printf("slack List: %v\n\n", slackHooks)
	fmt.Fprintf(w, "slack added: %s", slack)
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		slackToken = cleanBody(body)
	} else {
		slackToken = r.URL.Query().Get("slack")
	}
	fmt.Fprintf(w, "slack token added: %s", slackToken)
}

func removeSlackHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler: removeEmailHandler")
	slack := ""
	if len(r.URL.RawQuery) == 0 {
		body, _ := ioutil.ReadAll(r.Body)
		slack = cleanBody(body)
	} else {
		slack = r.URL.Query().Get("slack")
	}
	fmt.Printf("slackhook: '%s'\n", slack)
	if len(slack) > 0 {
		slackHooks = removeSlackHook(slack)
	}
	fmt.Printf("Slack List: %v\n\n", slackHooks)
	fmt.Fprintf(w, "Slack hook removed: %s", slack)
}

func removeSlackHook(toremove string) []string {
	slackList := []string{}
	for _, slack := range slackHooks {
		if slack != toremove {
			slackList = append(slackList, slack)
		}
	}
	return slackList
}

func addSlackHook(slack string) []string {
	return append(slackHooks, slack)
}

func listEmailsHandler(w http.ResponseWriter, r *http.Request) {
	PrintList(w, emailList)
}

func enableEmailHandler(w http.ResponseWriter, r *http.Request) {
	notificationsEnabled = true
	fmt.Fprintf(w, "email enabled")
}

func disableEmailHandler(w http.ResponseWriter, r *http.Request) {
	notificationsEnabled = false
	fmt.Fprintf(w, "email disabled")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	tmplt.ExecuteTemplate(w, "test.html", Page{Title: "Square Inventory"})
}

func getApkHandler(w http.ResponseWriter, r *http.Request) {
	apk, _ := ioutil.ReadFile("../mobile/apk/app-debug-unaligned.apk")
	fmt.Fprintf(w, "%s", apk)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	html, _ := ioutil.ReadFile("../static/test.html")
	fmt.Fprintf(w, "%s", html)
}

func requestListHandler(w http.ResponseWriter, r *http.Request) {
	var itemReq string
	if strings.Contains(r.URL.Path, "add") {
		if len(r.URL.RawQuery) == 0 {
			body, _ := ioutil.ReadAll(r.Body)
			itemReq = cleanBody(body)
		} else {
			itemReq = r.URL.Query().Get("item")
		}
		requestList = append(requestList, itemReq)
	} else if strings.Contains(r.URL.Path, "view") {
		PrintList(w, requestList)
	} else if strings.Contains(r.URL.Path, "clear") {
		requestList = []string{}
	} else {
		fmt.Fprintf(w, "ok")
	}
}

func PrintList(w http.ResponseWriter, list []string) {
	for _, item := range list {
		fmt.Fprintf(w, fmt.Sprintf("%s\n", item))
	}
}

func handleHandlers() {
	http.Handle("/apk/", http.StripPrefix("/apk/", http.FileServer(http.Dir(os.Getenv("APK_PATH")))))

	http.HandleFunc("/_status", statusHandler)
	http.HandleFunc("/setcount", countHandler)
	http.HandleFunc("/setslacksecret", secretHandler)
	http.HandleFunc("/lastcount", lastCountHandler)
	http.HandleFunc("/addemail", addEmailHandler)
	http.HandleFunc("/listemails", listEmailsHandler)
	http.HandleFunc("/removeemail", removeEmailHandler)
	http.HandleFunc("/sendemail", sendEmailHandler)
	http.HandleFunc("/setimage", setImageHandler)
	http.HandleFunc("/getimage", getImageHandler)
	//http.HandleFunc("/apk", getApkHandler)
	http.HandleFunc("/addslackhook", addSlackHandler)
	http.HandleFunc("/removelackhook", removeSlackHandler)
	http.HandleFunc("/enableemail", enableEmailHandler)
	http.HandleFunc("/disableemail", disableEmailHandler)
	http.HandleFunc("/addreq", requestListHandler)
	http.HandleFunc("/viewreq", requestListHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/clearreq", requestListHandler)
	http.HandleFunc("/", defaultHandler)
}

func main() {
	fmt.Println("****STARTING THE SERVER!")
	go emailSender()
	handleHandlers()
	startFSM()
	http.ListenAndServe(":8080", nil)
}
