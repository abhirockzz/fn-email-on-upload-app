# Function which triggers an email when a new object is uploaded to OCI Object Store

This Function will send an email notification (to a configurable address) in the event of an object uploaded to a bucket in [Oracle Cloud Infrastructure Object Storage](https://docs.cloud.oracle.com/iaas/Content/Object/Concepts/objectstorageoverview.htm). Email notification is powered by [Oracle Cloud Infrastructure Email Delivery](https://docs.cloud.oracle.com/iaas/Content/Email/Concepts/overview.htm)

It's a Go function just uses the plain old [SMTP package in Go](https://golang.org/pkg/net/smtp/) to send emails

## Pre-requisites

### Configure OCI Email Delivery

- [Generate SMTP Credentials for a User](https://docs.cloud.oracle.com/iaas/Content/Email/Tasks/generatesmtpcredentials.htm) - you'll have to configure these user credentials in the app (`OCI_EMAIL_DELIVERY_USER_OCID` and `OCI_EMAIL_DELIVERY_USER_PASSWORD` variables)
- [Add approved sender](https://docs.cloud.oracle.com/iaas/Content/Email/Tasks/managingapprovedsenders.htm) - use `OCI_EMAIL_DELIVERY_APPROVED_SENDER` parameter to configure this in the app
- [Note down value for the SMTP server](https://docs.cloud.oracle.com/iaas/Content/Email/Tasks/configuresmtpconnection.htm) - it'll be used in the `OCI_EMAIL_DELIVERY_SMTP_SERVER` configuration attribute

Clone this repo

### Switch to correct context

- `fn use context <your context name>`
- Check using `fn ls apps`

### Create app

`fn create app --annotation oracle.com/oci/subnetIds=<SUBNETS> --config OCI_EMAIL_DELIVERY_USER_OCID=<OCI_EMAIL_DELIVERY_USER_OCID> --config OCI_EMAIL_DELIVERY_USER_PASSWORD=<OCI_EMAIL_DELIVERY_USER_PASSWORD> --config REGION=<REGION> --config OCI_EMAIL_DELIVERY_SMTP_SERVER=<OCI_EMAIL_DELIVERY_SMTP_SERVER> --config OCI_EMAIL_DELIVERY_APPROVED_SENDER=<OCI_EMAIL_DELIVERY_APPROVED_SENDER> --config EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS=<EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS> fn-email-on-upload-app`

> please provide a valid email address for `EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS`

e.g.

`fn create app --annotation oracle.com/oci/subnetIds='["ocid1.subnet.oc1.phx.aaaaaaaaghmsma7mpqhqdhbgnby25u2zo4wqlrrcskvu7jg56dryxt3hgvkz"]' --config OCI_EMAIL_DELIVERY_USER_OCID=ocid1.user.oc1..aaaaaaaa4seqx6jeyma46ldy4cbuv42q4l26scz5p4rkz3rauuoioo42qwmq@ocid1.tenancy.oc1..aaaaaaaaydrjm77otncda2xn7qtv7l3hqnd3zxn2u6siwdhniibwfv4wwhta.3n.com --config OCI_EMAIL_DELIVERY_USER_PASSWORD='s3cr3t:-)' --config OCI_EMAIL_DELIVERY_SMTP_SERVER=smtp.us-phoenix-1.oraclecloud.com --config OCI_EMAIL_DELIVERY_APPROVED_SENDER=test@test.com --config EMAIL_NOTIFICAITON_RECEPIENT_ADDRESS=test@gmail.com fn-email-on-upload-app`

**Check**

`fn inspect app fn-email-on-upload-app`

## Moving on...

Deploy the app...

`cd fn-email-on-upload-app` and `fn -v deploy --app fn-email-on-upload-app`

**Test**

To test without end-to-end Events integration, just simulate the `create` (or upload) action by using `dummy-events-payload.json`

`cat dummy-events-payload.json | fn invoke fn-email-on-upload-app notifyonupload` - if all goes well

- the configured recepient will receive an email with the subject `File 'testfile.txt' uploaded to bucket 'test-bucket'`, and,
- the function will respond with `Notification email sent successfully!` message

After setting up Events integration, create a storage bucket in your tenancy, upload a file and it'll trigger the rest of the flow (as mentioned above)