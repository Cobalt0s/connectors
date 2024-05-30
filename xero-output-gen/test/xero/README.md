# Running instructions

## Step 1: prepare "xero-creds.json" file

Create a file called "xero-creds.json" in the root of the project with the following contents

    e.g. {
        "CLIENT_ID": "<client id goes here>",
        "CLIENT_SECRET": "<client secret goes here>",
        "ACCESS_TOKEN": "<access token goes here>",
        "REFRESH_TOKEN": "<refresh token goes here>"
    }

or export to an environment variable XERO_CRED_FILE by following command

$> export XERO_CRED_FILE=./xero-creds.json # or the path to your xero-creds.json file


In 1password, you can find a Xero creds.json file in the "Shared" vault.
The 1password item has an attached file called "creds.json" that contains the JSON.

## Step 2: run the following command

    $> go run test/xero/read/main.go
