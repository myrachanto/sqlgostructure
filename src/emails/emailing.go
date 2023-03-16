package emails

import (
	"fmt"
	"log"
	"os"

	"github.com/go-mail/mail"
	"github.com/joho/godotenv"
)

func ToEmail(name, customeremail, code, shopalias, bizname string, amount float64) {
	if customeremail == "" {
		log.Println("No customer email ....")
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes", err)
	}
	AppEmail := os.Getenv("AppEmail")
	Emailpass := os.Getenv("Emailpass")
	// secret := os.Getenv("Secret")
	// Customid := os.Getenv("Customid")
	// Owner := os.Getenv("Owner")
	OwnerEmail := "nillaveecakes@gmail.com"
	// OwnerEmail := os.Getenv("OwnerEmail")
	WebsiteLink := os.Getenv("WebsiteLink")
	Phone := os.Getenv("Phone")
	m := mail.NewMessage()

	m.SetHeader("From", OwnerEmail)
	m.SetHeader("To", customeremail, bizname)
	if bizname == "Nillavee" {
		m.SetHeader("To", customeremail, "nillaveecakes@gmail.com", "mutiapatrick35@gmail.com", "myrachanto@gmail.com", bizname)
		m.SetAddressHeader("Cc", OwnerEmail, "Nillavees Cakes and Patries")

		m.SetHeader("Subject", "Your Order at Nillavee Cakes and Patries was Successful!")

		m.SetBody("text/html", "<h3>Hello  "+name+"<br /> Thank you for shopping with us!</h3>Your Order of Ksh "+fmt.Sprintf("%.2f", amount)+" is being processed <br />We'll respond shortly.<br><br>Thanks!<br>"+OwnerEmail+"<br />Phone:"+Phone+"<br />Website: <a href='"+WebsiteLink+"'>"+WebsiteLink+"</a>")
	}

	m.SetAddressHeader("Cc", OwnerEmail, bizname)

	m.SetHeader("Subject", "Your Order at "+bizname+" was Successful!")

	m.SetBody("text/html", "<h3>Hello  "+name+"<br /> Thank you for shopping with us!</h3>Your Order of Ksh "+fmt.Sprintf("%.2f", amount)+" is being processed <br />We'll respond shortly.<br><br>Thanks!<br")

	// m.Attach("logo.png")

	d := mail.NewDialer("smtp.gmail.com", 587, AppEmail, Emailpass)

	// Send the email to Kate, Noah and Oliver.

	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}

}

func ToEmailPassword(pass, email string) {
	if email == "" {
		log.Println("No customer email ....")
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes", err)
	}
	AppEmail := os.Getenv("AppEmail")
	Emailpass := os.Getenv("Emailpass")
	OwnerEmail := "noreply@chantosweb.com"
	m := mail.NewMessage()

	m.SetHeader("From", OwnerEmail)

	m.SetHeader("To", email)

	m.SetAddressHeader("Cc", email, "New Password")

	m.SetHeader("Subject", "Your New Password is here!")

	m.SetBody("text/html", "<h3>Hello  "+email+"<br /> Your new password is "+pass+"!</h3><br />")

	d := mail.NewDialer("smtp.gmail.com", 587, AppEmail, Emailpass)
	if err := d.DialAndSend(m); err != nil {

		panic(err)

	}

}
