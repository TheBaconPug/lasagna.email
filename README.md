# lasagna.email

A easy to use temporary email service

Try it out here: https://lasagna.email

# Setup

## Requirements:

Mailgun Account\
Go

## Step one: Clone the repository
``git clone https://github.com/TheBaconPug/lasagna.email.git``

## Step two: Setup Mailgun
### Create a new route

Set the expression type to "Match Recipient"\
Set the recipient to ``.*@yourdomain``\
Turn on "Store and notify" and set the value to ``https://yourdomain/api/callback``\
Turn on "Stop"\
Set the priority to 0

## Step three: Edit the config

Set "port" to the port you're going to use\
Set "domain" to your domain name\
Rename "config.example.json" to "config.json"

## Step four: Start the site

Run ``go run .`` or ``go build && ./lasagnamail ``
