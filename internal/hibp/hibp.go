package hibp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nhoya/goPwned"
)

//Options contains the options for HIBP module
type Options struct {
	Mails        []string
	JSONFlag     bool
	GetPasswords bool
}

//Report contains the report for dumps containing a specific email address
type report struct {
	Pwnd []*PwnedEntity `json:"pwnd"`
}

//PwnedEntity is the struct that contains the mail address and the breaches
type PwnedEntity struct {
	Email     string   `json:"email"`
	Breaches  []string `json:"breaches"`
	Passwords []string `json:"passwords,omitempty"`
}

//StartHIBP is the init function for the HIBP module
func (opts *Options) StartHIBP() {
	report := new(report)
	for n, mail := range opts.Mails {
		report.getBreachesForMail(mail, opts.GetPasswords)
		if n != 0 {
			//prevent antiflood block
			time.Sleep(time.Second * 2)
		}
	}
	report.printHIBPReport(opts.JSONFlag)
}

func (report *report) getBreachesForMail(mail string, pwds bool) {
	fmt.Println("[+] Dump for " + mail)
	stuff, err := gopwned.GetAllBreachesForAccount(mail, "", "true")
	if err == nil {
		pwnd := new(PwnedEntity)
		pwnd.Email = mail
		for _, data := range stuff {
			pwnd.Breaches = append(pwnd.Breaches, data.Name)
		}
		if pwds {
			pwnd.Passwords, _ = getPasswords(mail)
		}
		report.Pwnd = append(report.Pwnd, pwnd)
	}
}

func (report *report) printHIBPReport(jsonFlag bool) {
	if jsonFlag {
		jsonreport, _ := json.Marshal(&report)
		fmt.Println(string(jsonreport))
	} else {
		for _, k := range report.Pwnd {
			fmt.Println("Mail:", k.Email)
			fmt.Println("Breaches:", k.Breaches)
			if k.Passwords != nil {
				fmt.Println("Passwords:", k.Passwords)
			}
		}
	}
}
