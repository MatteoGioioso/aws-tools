# AWS local key rotation

A small utility/CLI to help your rotate local aws credentials

## What it does
For the current profile:

1. Create new Access keys
2. Stash and save the old ones
3. Make the old keys inactive
4. Set the new keys into your `.aws/credentials` file


## Installation

```
go get ...
```

## Usage
You can simply run it as a command

```
LKR_SAVE_OLD_KEYS=0 ./lkr
``` 

or use a cron job to regularly rotate the keys

## Settings
For the time being you can set the CLI only through environmental variables

- `LKR_DELETE_OLD_KEYS`: by default LKR will just deactivate the old keys from IAM, if you wish to delete them then set this
  variable to `yes`
- `LKR_BACKUP_OLD_KEYS`:
  by default LKR will create a new file in `~/.aws/` before rotating the credentials, this file is going to be named `credentials-inactive-<Epoch timestamp>`, 
  if sets to `no` it will not back it up.
  
  
The rest of the settings are the same as the AWS CLI: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html