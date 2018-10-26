package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/smtp"
	"time"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(objectUploadEmailNotificationHandler))
}

func objectUploadEmailNotificationHandler(ctx context.Context, in io.Reader, out io.Writer) {
	log.Println("objectUploadEmailNotificationHandler invoked on", time.Now())

	var evt OCIEvent
	json.NewDecoder(in).Decode(&evt)
	log.Println("Got OCI event", evt)

	var casperDetails Data
	json.Unmarshal([]byte(evt.Data), &casperDetails)
	log.Println("Got Casper details", casperDetails)

	fnCtx := fdk.GetContext(ctx)

	config := fnCtx.Config()
	username := config["OCI_EMAIL_DELIVERY_USER_OCID"]
	password := config["OCI_EMAIL_DELIVERY_USER_PASSWORD"]
	ociSMTPServer := config["OCI_EMAIL_DELIVERY_SMTP_SERVER"]
	approvedOCIEmailDeliverySender := config["OCI_EMAIL_DELIVERY_APPROVED_SENDER"]
	emailRecepientAddress := config["EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS"]

	log.Println("OCI_EMAIL_DELIVERY_USER_OCID", username)
	log.Println("OCI_EMAIL_DELIVERY_USER_PASSWORD", password)
	log.Println("OCI_EMAIL_DELIVERY_SMTP_SERVER", ociSMTPServer)
	log.Println("OCI_EMAIL_DELIVERY_APPROVED_SENDER", approvedOCIEmailDeliverySender)
	log.Println("EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS", emailRecepientAddress)

	response := sendEmailNotification(username, password, ociSMTPServer, approvedOCIEmailDeliverySender, emailRecepientAddress, casperDetails)
	//response := "Object named " + casperDetails.ObjectName + " uploaded to bucket " + casperDetails.BucketName
	log.Println("Response", response)
	out.Write([]byte(response))
}

func sendEmailNotification(username, password, ociSMTPServer, approvedOCIEmailDeliverySender, emailRecepientAddress string, casperDetails Data) string {
	log.Println("sending email notification")

	auth := smtp.PlainAuth("", username, password, ociSMTPServer)

	to := []string{emailRecepientAddress}
	subject := "File '" + casperDetails.ObjectName + "' uploaded to bucket '" + casperDetails.BucketName + "'"
	body := subject + " on " + casperDetails.CreationTime + " in tenancy " + casperDetails.TenantID
	msg := []byte("To: " + emailRecepientAddress + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	log.Println("Message ", string(msg))
	err := smtp.SendMail(ociSMTPServer+":25", auth, approvedOCIEmailDeliverySender, to, msg)
	if err != nil {
		log.Println("Error sending notification email", err.Error())
		return "Error sending notification email " + err.Error()
	}

	log.Println("Notification email sent successfully!")
	return "Notification email sent successfully!"
}

//OCIEvent ...
type OCIEvent struct {
	EventType          string `json:"eventType"`
	EventTypeVersion   string `json:"eventTypeVersion"`
	ContentType        string `json:"contentType"`
	CloudEventsVersion string `json:"cloudEventsVersion"`
	SchemaURL          string `json:"schemaURL"`
	Source             string `json:"source"`
	EventID            string `json:"eventID"`
	EventTime          string `json:"eventTime"`
	Extensions         `json:"extensions"`
	Data               string `json:"data"`
}

//Extensions - "extension" attribute in events JSON payload
type Extensions struct {
	CompartmentId string `json:"compartmentId"`
}

//Data - represents Casper data
type Data struct {
	TenantID     string `json:"tenantId"`
	BucketOcid   string `json:"bucketOcid"`
	BucketName   string `json:"bucketName"`
	API          string `json:"api"`
	ObjectName   string `json:"objectName"`
	ObjectEtag   string `json:"objectEtag"`
	ResourceType string `json:"resourceType"`
	Action       string `json:"action"`
	CreationTime string `json:"creationTime"`
}
