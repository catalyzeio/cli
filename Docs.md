# Automatic Updates

Once downloaded, the CLI will attempt to automatically update itself when a new version becomes available. This ensures you are always running a compatible version of the Catalyze CLI. However you can always check out the latest releases on the [releases page](https://github.com/catalyzeio/cli/releases).

To ensure your CLI can automatically update itself, be sure to put the binary in a location where you have write access without the need for sudo or escalated privileges.

# Supported Platforms

Since version 2.0.0, the following platforms and architectures are supported by the Catalyze CLI.

| OS | Architecture |
|----|--------------|
| Darwin (Mac OS X) | 64-bit, 32-bit |
| Linux | 64-bit, 32-bit |
| Windows | 64-bit, 32-bit |

# Global Scope

The CLI now supports the concept of scope. Previous to version 2.0.0, all commands had to be run within an associated local git repo. Now, the only time you need to be in a local git repo is when you associate to a new environment. After the initial association, CLI commands can be run from any directory. If you have more than one environment, the CLI uses this concept of scope to decide which environment you are using for the command.

Let's say you have an environment that you associated in the directory `~/mysandbox-code` and another you associated in the directory `~/myprod-code`. These environments are named `mysandbox` and `myprod` respectively. When you are within either of those directories, the CLI knows that any command you run will be in the context of that given environment. Commands run in the `~/myprod-code` directory will be run against the `myprod` environment. Similarly for `~/mysandbox-code` and the `mysandbox` environment. What if you are outside those directories? You have three options.

First, you can tell the CLI which environment you want to use with the global option `-E` or `--env` (see [Global Options](https://resources.catalyze.io/paas/cli/sections/global-options/)). Your command might start like this

```
catalyze -E myprod ...
```

This global option will even override the environment found in a local git repo. If you don't set the `-E` flag, and the CLI can't find an environment in your local git repo, the CLI then checks for a default environment. A default environment is used whenever you are outside of a git repo and an environment is not specified. A default environment can be specified using the [default](https://resources.catalyze.io/paas/cli/sections/default/) command. You can find out which environment is the default by running the [associated](https://resources.catalyze.io/paas/cli/sections/associated/) command.

Lastly, if no environment is specified, you're outside of a git repo, and no default environment is set, then the CLI simply takes the first environment you associated and prompts you to continue with this environment. This concept of scope will make it easier for Catalyze customers with multiple environments to use the CLI!

# Environment Aliases


When you associate an environment from within a local git repo, you typically run the following command:

```
catalyze associate "My Health Tech Company Production" app01
```

Where `My Health Tech Company Production` is the name of your environment. However with the concept of [scope](https://resources.catalyze.io/paas/cli/sections/global-scope/) and being able to specify which environment to use on a command by command basis with the `-E` global option, that is a lot to type! This is where environment aliases come in handy.

When you associate an environment and you want to pick a shorter name to reference the environment by, simply add a `-a` flag to the command. Let's try the command again calling it `prod` this time:

```
catalyze associate "My Health Tech Company Production" app01 -a prod
```

Now when you run the [associated](https://resources.catalyze.io/paas/cli/sections/associated/) command, you will see the alias as well as the actual environment name.

When using aliases, there are a couple things to keep in mind. Aliases are only local and never leave your local machine. If you alias this environment `prod`, a coworker can alias the environment `healthtech-prod` with no ramifications. Second, after setting an alias you will never reference the environment by its actual name with the CLI. You will always use the alias for flags, arguments, options, etc.

To change or remove an alias, you must [disassociate](https://resources.catalyze.io/paas/cli/sections/disassociate/) and then [reassociate](https://resources.catalyze.io/paas/cli/sections/associate/) with a new alias.

# Bash Autocompletion

One feature we've found helpful on \*Nix systems is autocompletion in bash. To enable this feature, head over to the github repo and download the `catalyze_autocomplete` file. If you use a Mac, you will need to install bash-completion with `brew install bash-completion` or `source` the `catalyze_autocomplete` file each time you start up terminal. Store this file locally in `/etc/bash_completion.d/` or (`/usr/local/etc/bash_completion.d/` on Mac). Completion will be available when you restart terminal. Now simply type `catalyze ` and hit tab twice to see the list of available commands. **Please note** that autocompletion only works one level deep. The CLI will not autocomplete or suggest completions when you type `catalyze db ` and then hit tab twice. It currently only works when you have just `catalyze ` typed into your terminal. This is a feature we are looking into expanding in the future.

Note: you may have to add `source /etc/bash_completion.d/catalyze_autocomplete` (`/usr/local/etc/bash_completion.d/catalyze_autocomplete`) in your `~/.bashrc` (`~/.bash_profile`) file.

# Global Options

The following table outlines all global options available in the CLI. Global options are always set after the word `catalyze` and before any commands. Rather than setting these each time, you may also set an environment variable with the appropriate value which will automatically be used.

| Short Name | Long Name | Description | Environment Variable |
|------------|-----------|-------------|----------------------|
| -U | --username | Your catalyze username that you login to the Dashboard with | CATALYZE_USERNAME |
| -P | --password | Your catalyze password that you login to the Dashboard with | CATALYZE_PASSWORD |
| -E | --env | The local alias of the environment in which this command will be run. Read more about [environment aliases](https://resources.catalyze.io/paas/cli/sections/environment-aliases/) | CATALYZE_ENV |
| -v | --version | Prints out the CLI version | |

# Overview

Usage: catalyze [OPTIONS] COMMAND [arg...]

Catalyze CLI. Version 2.3.0

Options:

```
  -U, --username        Catalyze Username ($CATALYZE_USERNAME)
  -P, --password        Catalyze Password ($CATALYZE_PASSWORD)
  -E, --env             The local alias of the environment in which this command will be run ($CATALYZE_ENV)
  -v, --version=false   Show the version and exit
```

Commands:

```
  associate      Associates an environment
  associated     Lists all associated environments
  certs          Manage your SSL certificates and domains
  console        Open a secure console to a service
  dashboard      Open the Catalyze Dashboard in your default browser
  db             Tasks for databases
  default        Set the default associated environment
  disassociate   Remove the association with an environment
  environments   List all environments you have access to
  files          Tasks for managing service files
  invites        Manage invitations for your organizations
  keys           Tasks for SSH keys
  logs           Show the logs in your terminal streamed from your logging dashboard
  logout         Clear the stored user information from your local machine
  metrics        Print service and environment metrics in your local time zone
  rake           Execute a rake task
  redeploy       Redeploy a service without having to do a git push
  services       List all services for your environment
  sites          Tasks for updating sites, including hostnames, SSL certificates, and private keys
  ssl            Perform operations on local certificates to verify their validity
  status         Get quick readout of the current status of your associated environment and all of its services
  supportids     Print out various IDs related to your associated environment to be used when contacting Catalyze support
  update         Checks for available updates and updates the CLI if a new update is available
  users          Manage users who have access to the given organization
  vars           Interaction with environment variables for the associated environment
  whoami         Retrieve your user ID
  worker         Start a background worker
  version        Output the version and quit
```

Run 'catalyze COMMAND --help' for more information on a command.

# Associate

```
Usage: catalyze associate ENV_NAME SERVICE_NAME [-a] [-r] [-d]

Associates an environment

Arguments:
  ENV_NAME=""       The name of your environment
  SERVICE_NAME=""   The name of the primary code service to associate with this environment (i.e. 'app01')

Options:
  -a, --alias=""            A shorter name to reference your environment by for local commands
  -r, --remote="catalyze"   The name of the remote
  -d, --default=false       Specifies whether or not the associated environment will be the default
```

`associate` is the entry point of the cli. You need to associate an environment before you can run most other commands. Check out [scope](https://resources.catalyze.io/paas/cli/sections/global-scope/) and [aliases](https://resources.catalyze.io/paas/cli/sections/environment-aliases/) for more info on the value of the alias and default options. Here is a sample command

```
catalyze associate My-Production-Environment app01 -a prod -d
```

# Associated

```
Usage: catalyze associated

Lists all associated environments
```

`associated` outputs information about all previously associated environments on your local machine. The information that is printed out includes the alias, environment ID, actual environment name, service ID, the git repo directory, and whether or not it is the default environment. Here is a sample command

```
catalyze associated
```

# Certs

The `certs` command gives access to certificate and private key management for public facing services. The certs command can not be run directly but has sub commands.

# Certs Create

```
Usage: catalyze certs create HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]

Create a new domain with an SSL certificate and private key

Arguments:
  HOSTNAME=""           The hostname of this domain and SSL certificate plus private key pair
  PUBLIC_KEY_PATH=""    The path to a public key file in PEM format
  PRIVATE_KEY_PATH=""   The path to an unencrypted private key file in PEM format

Options:
  -s, --self-signed=false   Whether or not the given SSL certificate and private key are self signed
  -r, --resolve=true        Whether or not to attempt to automatically resolve incomplete SSL certificate issues
```

`certs create` allows you to upload an SSL certificate and private key which can be used to secure your public facing code service. Cert creation can be done at any time, even after environment provisioning, but must be done before [creating a site](https://resources.catalyze.io/paas/cli/sections/sites-create/). When creating a cert, the CLI will check to ensure the certificate and private key match and the given hostname is valid for the given certificate. If you are using a self signed cert, pass in the `-s` flag and the hostname check will be skipped. Catalyze requires that your certificate include your own certificate, intermediate certificates, and the root certificate in that order. If you only include your certificate, the CLI will attempt to resolve this and fetch intermediate and root certificates for you. It is advised that you create a full chain before running this command as the `-r` flag is accomplished on a "best effort" basis. Here is a sample command

```
catalyze certs create mywebsite.com ~/path/to/cert.pem ~/path/to/priv.key
```

# Certs List

```
Usage: catalyze certs list

List all existing domains that have SSL certificate and private key pairs
```

`certs list` lists all of the available certs you have created on your environment. The displayed names are the names that should be used as the `DOMAIN` parameter in the [sites create](https://resources.catalyze.io/paas/cli/sections/sites-create/) command. Here is a sample command

```
catalyze certs list
```

# Certs RM

```
Usage: catalyze certs rm HOSTNAME

Remove an existing domain and its associated SSL certificate and private key pair

Arguments:
  HOSTNAME=""   The hostname of the domain and SSL certificate and private key pair
```

`certs rm` allows you to delete old certificate and private key pairs. Only certs that are not in use by a site can be deleted. Here is a sample command

```
catalyze certs rm mywebsite.com
```

# Certs Update

```
Usage: catalyze certs update HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]

Update the SSL certificate and private key pair for an existing domain

Arguments:
  HOSTNAME=""           The hostname of this domain and SSL certificate and private key pair
  PUBLIC_KEY_PATH=""    The path to a public key file in PEM format
  PRIVATE_KEY_PATH=""   The path to an unencrypted private key file in PEM format

Options:
  -s, --self-signed=false   Whether or not the given SSL certificate and private key are self signed
  -r, --resolve=true        Whether or not to attempt to automatically resolve incomplete SSL certificate issues
```

`certs update` works nearly identical to the [certs create](https://resources.catalyze.io/paas/cli/sections/certs-create/) command. All rules regarding self signed certs and certificate resolution from the `certs create` command apply to the `certs update` command. This is useful for when your certificates have expired and you need to upload new ones. Simply update your certs, then redeploy your services. Here is a sample command

```
catalyze certs update mywebsite.com ~/path/to/new/cert.pem ~/path/to/new/priv.key
```

# Console

```
Usage: catalyze console SERVICE_NAME [COMMAND]

Open a secure console to a service

Arguments:
  SERVICE_NAME=""   The name of the service to open up a console for
  COMMAND=""        An optional command to run when the console becomes available
```

`console` gives you direct access to your database service or application shell. For example, if you open up a console to a postgres database, you will be given access to a psql prompt. You can also open up a mysql prompt, mongo cli prompt, rails console, django shell, and much more. When accessing a database service, the `COMMAND` argument is not needed because the appropriate prompt will be given to you. If you are connecting to an application service the `COMMAND` argument is required. Here are some sample commands

```
catalyze console db01
catalyze console app01 "bundle exec rails console"
```

# Dashboard

```
Usage: catalyze dashboard

Open the Catalyze Dashboard in your default browser
```

`dashboard` simply opens up the Catalyze Dashboard homepage in your default web browser. Here is a sample command

```
catalyze dashboard
```

# DB

The `db` command gives access to backup, import, and export services for databases. The db command can not be run directly but has sub commands.

# DB Backup

```
Usage: catalyze db backup DATABASE_NAME [-s]

Create a new backup

Arguments:
  DATABASE_NAME=""   The name of the database service to create a backup for (i.e. 'db01')

Options:
  -s, --skip-poll=false   Whether or not to wait for the backup to finish
```

`db backup` creates a new backup for the given database service. The backup is started and unless `-s` is specified, the CLI will poll every few seconds until it finishes. Regardless of a successful backup or not, the logs for the backup will be printed to the console when the backup is finished. If an error occurs and the logs are not printed, you can use the [db logs](https://resources.catalyze.io/paas/cli/sections/db-logs/) command to print out historical backup job logs. Here is a sample command

```
catalyze db backup db01
```

# DB Download

```
Usage: catalyze db download DATABASE_NAME BACKUP_ID FILEPATH [-f]

Download a previously created backup

Arguments:
  DATABASE_NAME=""   The name of the database service which was backed up (i.e. 'db01')
  BACKUP_ID=""       The ID of the backup to download (found from "catalyze backup list")
  FILEPATH=""        The location to save the downloaded backup to. This location must NOT already exist unless -f is specified

Options:
  -f, --force=false   If a file previously exists at "filepath", overwrite it and download the backup
```

`db download` downloads a previously created backup to your local hard drive. Be careful using this command is it could download PHI. Be sure that all hard drive encryption and necessary precautions have been taken before performing a download. The ID of the backup is found by first running the [db list](https://resources.catalyze.io/paas/cli/sections/db-list/) command. Here is a sample command

```
catalyze db download db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203 ./db.sql
```

This assumes you are downloading a MySQL or PostgreSQL backup which takes the `.sql` file format. If you are downloading a mongo backup, the command might look like this

```
catalyze db download db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203 ./db.tar.gz
```

# DB Export

```
Usage: catalyze db export DATABASE_NAME FILEPATH [-f]

Export data from a database

Arguments:
  DATABASE_NAME=""   The name of the database to export data from (i.e. 'db01')
  FILEPATH=""        The location to save the exported data. This location must NOT already exist unless -f is specified

Options:
  -f, --force=false   If a file previously exists at `filepath`, overwrite it and export data
```

`db export` is a simple wrapper around the `db backup` and `db download` commands. When you request an export, a backup is created that will be added to the list of backups shown when you perform the [db list](https://resources.catalyze.io/paas/cli/sections/db-list/) command. Then that backup is immediately downloaded. Regardless of a successful export or not, the logs for the backup will be printed to the console when the export is finished. If an error occurs and the logs are not printed, you can use the [db logs](https://resources.catalyze.io/paas/cli/sections/db-logs/) command to print out historical backup job logs. Here is a sample command

```
catalyze db export db01 ./dbexport.sql
```

This assumes you are exporting a MySQL or PostgreSQL database which takes the `.sql` file format. If you are exporting a mongo database, the command might look like this

```
catalyze db export db01 ./dbexport.tar.gz
```

# DB Import


```
Usage: catalyze db import DATABASE_NAME FILEPATH [-d [-c]]

Import data into a database

Arguments:
  DATABASE_NAME=""   The name of the database to import data to (i.e. 'db01')
  FILEPATH=""        The location of the file to import to the database

Options:
  -c, --mongo-collection=""   If importing into a mongo service, the name of the collection to import into
  -d, --mongo-database=""     If importing into a mongo service, the name of the database to import into
```

`db import` allows you to inject new data into your database service. For example, if you wrote a simple SQL file

```
CREATE TABLE mytable (
id TEXT PRIMARY KEY,
val TEXT
);

INSERT INTO mytable (id, val) values ('1', 'test');
```

and stored it at `./db.sql` you could import this into your database service. When importing data into mongo, you may specify the database and collection to import into using the `-d` and `-c` flags respectively. Regardless of a successful import or not, the logs for the import will be printed to the console when the import is finished. Before an import takes place, your database is backed up automatically in case any issues arise. Here is a sample command

```
catalyze db import db01 ./db.sql
```

# DB List

```
Usage: catalyze db list DATABASE_NAME [-p] [-n]

List created backups

Arguments:
  DATABASE_NAME=""   The name of the database service to list backups for (i.e. 'db01')

Options:
  -p, --page=1         The page to view
  -n, --page-size=10   The number of items to show per page
```

`db list` lists all previously created backups. After listing backups you can copy the backup ID and use it to [download](https://resources.catalyze.io/paas/cli/sections/db-download/) that backup or [view the logs](https://resources.catalyze.io/paas/cli/sections/db-logs/) from that backup. Here is a sample command

```
catalyze db list db01
```

# DB Logs

```
Usage: catalyze db logs DATABASE_NAME BACKUP_ID

Print out the logs from a previous database backup job

Arguments:
  DATABASE_NAME=""   The name of the database service (i.e. 'db01')
  BACKUP_ID=""       The ID of the backup to download logs from (found from "catalyze backup list")
```

`db logs` allows you to view backup logs from historical backup jobs. You can find the backup ID from using the `db list` command. Here is a sample command

```
catalyze db logs db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203
```

# Default

```
Usage: catalyze default ENV_ALIAS

Set the default associated environment

Arguments:
  ENV_ALIAS=""   The alias of an already associated environment to set as the default
```

`default` sets the default environment for all commands that don't specify an environment with the `-E` flag or commands that are run outside of a git repo. See [scope](https://resources.catalyze.io/paas/cli/sections/global-scope/) for more information on scope and default environments. When setting a default environment, you must give the alias of the environment if one was set when it was associated and not the real environment name. Here is a sample command

```
catalyze default prod
```

# Deploy Keys

The `deploy-keys` command gives access to SSH deploy keys for environment services. The deploy-keys command can not be run directly but has sub commands.

# Deploy Keys Add

```
Usage: catalyze deploy-keys add NAME KEY_PATH SERVICE_NAME [-p]

Add a new deploy key

Arguments:
  NAME=""              The name for the new key, for your own purposes
  PUBLIC_KEY_PATH=""   Relative path to the public key file
  SERVICE_NAME=""      The name of the code service to add this deploy key to

Options:
  -p, --private=false   Whether or not this is a private key
```

`deploy-keys add` allows you to upload an SSH public key or SSH private key in OpenSSH format. These keys are used for pushing code to your code services but are not required. You may use personal SSH keys with the [keys](https://resources.catalyze.io/paas/cli/sections/keys/) command instead. Deploy keys are useful for Continuous Integration or Continuous Deployment scenarios and are intended to be shared among an organization. Here are some sample commands

```
catalyze deploy-keys add app01_public ~/.ssh/app01_rsa.pub app01
catalyze deploy-keys add app01_private ~/.ssh/app01_rsa app01 -p
```

# Deploy Keys List

```
Usage: catalyze deploy-keys list SERVICE_NAME [--include-keys]

List all deploy keys

Arguments:
  SERVICE_NAME=""   The name of the code service to list deploy keys

Options:
  --include-keys=false   Print out the values of the deploy keys, as well as names
```

`deploy-keys list` will list all of your previously uploaded deploy keys by name. Optionally including the key's fingerprint in SHA256 format if the `--include-keys` flag is specified as well. Here is a sample command

```
catalyze deploy-keys list app01 --include-keys
```

# Deploy Keys Rm

```
Usage: catalyze deploy-keys rm NAME SERVICE_NAME [-p]

Remove a deploy key

Arguments:
  NAME=""           The name of the key to remove
  SERVICE_NAME=""   The name of the code service to remove this deploy key from

Options:
  -p, --private=false   Whether or not this is a private key
```

`deploy-keys rm` will remove a previously created deploy key by name. It is a good idea to rotate deploy keys on a set schedule as they are intended to be shared among an organization. Here are some sample commands

```
catalyze deploy-keys rm app01_public app01
catalyze deploy-keys rm app01_private app01 -p
```

# Disassociate

```
Usage: catalyze disassociate ENV_ALIAS

Remove the association with an environment

Arguments:
  ENV_ALIAS=""   The alias of an already associated environment to disassociate
```

`disassociate` does not have to be run from within a git repo. Disassociate removes the environment from your list of associated environments but **does not** remove the catalyze git remote on the git repo. Here is a sample command

```
catalyze disassociate myprod
```

# Environments

```
Usage: catalyze environments

List all environments you have access to
```

`environments` lists all environments that you are granted access to. These environments include those you created and those that other Catalyze customers have added you to. Here is a sample command

```
catalyze environments
```

# Files

The `files` command gives access to service files on your environment's services. Service files can include nginx configs, SSL certificates, and any other file that might be injected into your running service. The files command can not be run directly but has sub commands.

# Files Download

```
Usage: catalyze files download [SERVICE_NAME] FILE_NAME [-o] [-f]

Download a file to your localhost with the same file permissions as on the remote host or print it to stdout

Arguments:
  SERVICE_NAME="service_proxy"   The name of the service to download a file from
  FILE_NAME=""                   The name of the service file from running "catalyze files list"

Options:
  -o, --output=""     The downloaded file will be saved to the given location with the same file permissions as it has on the remote host. If those file permissions cannot be applied, a warning will be printed and default 0644 permissions applied. If no output is specified, stdout is used.
  -f, --force=false   If the specified output file already exists, automatically overwrite it
```

`files download` allows you to view the contents of a service file and save it to your local machine. Most service files are stored on your service_proxy and therefore you should not have to specify the `SERVICE_NAME` argument. Simply supply the `FILE_NAME` found from the [files list](https://resources.catalyze.io/paas/cli/sections/files-list/) command and the contents of the file, as well as the permissions string, will be printed to your console. You can always store the file locally, applying the same permissions as those on the remote server, by specifying an output file with the `-o` flag. Here is a sample command

```
catalyze files download /etc/nginx/sites-enabled/mywebsite.com
```

# Files List

```
Usage: catalyze files list [SERVICE_NAME]

List all files available for a given service

Arguments:
  SERVICE_NAME="service_proxy"   The name of the service to list files for
```

`files list` prints out a listing of all service files available for download. Nearly all service files are stored on the service_proxy and therefore you should not have to specify the `SERVICE_NAME` argument. Here is a sample command

```
catalyze files list
```

# Invites

The `invites` command gives access to organization invitations. Every environment is owned by an organization and users join organizations in order to access individual environments. You can invite new users by email and manage pending invites through the CLI. You cannot call the `invites` command directly, but must call one of its subcommands.

# Invites Accept

```
Usage: catalyze invites accept INVITE_CODE

Accept an organization invite

Arguments:
  INVITE_CODE=""   The invite code that was sent in the invite email
```

`invites accept` is an alternative form of accepting an invitation sent by email. The invitation email you receive will have instructions as well as the invite code to use with this command. Here is a sample command

```
catalyze invites accept 5a206aa8-04f4-4bc1-a017-ede7e6c7dbe2
```

# Invites List

```
Usage: catalyze invites list

List all pending organization invitations
```

`invites list` lists all pending invites for the associated environment's organization. Any invites that have already been accepted will not appear in this list. To manage users who have already accepted invitations or are already granted access to your environment, use the [users](https://resources.catalyze.io/paas/cli/sections/users/) group of commands. Here is a sample command

```
catalyze invites list
```

# Invites RM

```
Usage: catalyze invites rm INVITE_ID

Remove a pending organization invitation

Arguments:
  INVITE_ID=""   The ID of an invitation to remove
```

`invites rm` removes a pending invitation found by using the [invites list](https://resources.catalyze.io/paas/cli/sections/invites-list) command. Once an invite has already been accepted, it cannot be removed. Removing an invitation is helpful if an email was misspelled and an invitation was sent to an incorrect email address. If you want to revoke access to a user who already has been given access to your environment, use the [users rm](https://resources.catalyze.io/paas/cli/sections/users-rm) command. Here is a sample command

```
catalyze invites rm 78b5d0ed-f71c-47f7-a4c8-6c8c58c29db1
```

# Invites Send

```
Usage: catalyze invites send EMAIL [-m | -a]

Send an invite to a user by email for a given organization

Arguments:
  EMAIL=""     The email of a user to invite to the associated environment. This user does not need to have a Catalyze account prior to sending the invitation

Options:
  -m, --member=true   Whether or not the user will be invited as a basic member
  -a, --admin=false   Whether or not the user will be invited as an admin
```

`invites send` invites a new user to your environment's organization. The only piece of information required is the email address to send the invitation to. The invited user will join the organization as a basic member, unless otherwise specified with the `-a` flag. The recipient does **not** need to have a Dashboard account in order to send them an invitation. However, they will need to have a Dashboard account to accept the invitation. Here is a sample command

```
catalyze invites send coworker@catalyze.io -a
```

# Keys

The `keys` command gives access to SSH key management for your user account. SSH keys can be used for authentication and pushing code to the Catalyze platform. Any SSH keys added to your user account should not be shared but be treated as private SSH keys. Any SSH key uploaded to your user account will be able to be used with all code services and environments that you have access to. The keys command can not be run directly but has sub commands.

# Keys Add

```
Usage: catalyze keys add NAME PUBLIC_KEY_PATH

Add a public key

Arguments:
  NAME=""              The name for the new key, for your own purposes
  PUBLIC_KEY_PATH=""   Relative path to the public key file
```

`keys add` allows you to add a new SSH key to your user account. SSH keys added to your user account should be private and not shared with others. SSH keys can be used for authentication (as opposed to the traditional username and password) as well as pushing code to an environment's code services. Please note, you must specify the path to the public key file and not the private key. All SSH keys should be in either OpenSSH RSA format or PEM format. Here is a sample command

```
catalyze keys add my_prod_key ~/.ssh/prod_rsa.pub
```

# Keys List

```
Usage: catalyze keys list [OPTIONS]

List your public keys

Options:
  --include-keys=false   Print out the values of the public keys, as well as names.
```

`keys list` lists all public keys by name that have been uploaded to your user account, optionally printing out the public key itself. Here is a sample command

```
catalyze keys list --include-keys
```

# Keys RM

```
Usage: catalyze keys rm NAME

Remove a public key

Arguments:
  NAME=""      The name of the key to remove.
```

`keys rm` allows you to remove an SSH key previously uploaded to your account. The name of the key can be found by using the [keys list](https://resources.catalyze.io/paas/cli/sections/keys-list/) command. Here is a sample command

```
catalyze keys rm my_prod_key
```

# Keys Set

```
Usage: catalyze keys set PRIVATE_KEY_PATH

Set your auth key

Arguments:
  PRIVATE_KEY_PATH=""   Relative path to the private key file.
```

`keys set` allows the CLI to use an SSH key for authentication instead of the traditional username and password combination. This can be useful for automation or where a shared workstations are involved. Please note that you must pass in the path to the private key and not the public key. Here is a sample command

```
catalyze keys set ~/.ssh/my_key
```

# Logout

```
Usage: catalyze logout

Clear the stored user information from your local machine
```

When using the CLI, your username and password are **never** stored in any file on your filesystem. However, in order to not type in your username and password each and every command, a session token is stored in the CLI's configuration file and used until it expires. `logout` removes this session token from the configuration file. Here is a sample command

```
catalyze logout
```

# Logs

```
Usage: catalyze logs [QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds]

Show the logs in your terminal streamed from your logging dashboard

Arguments:
  QUERY="*"    The query to send to your logging dashboard's elastic search (regex is supported)

Options:
  -f, --follow=false   Tail/follow the logs (Equivalent to -t)
  -t, --tail=false     Tail/follow the logs (Equivalent to -f)
  --hours=0            The number of hours before now (in combination with minutes and seconds) to retrieve logs
  --minutes=1          The number of minutes before now (in combination with hours and seconds) to retrieve logs
  --seconds=0          The number of seconds before now (in combination with hours and minutes) to retrieve logs
```

`logs` prints out your application logs directly from your logging Dashboard. If you do not see your logs, try adjusting the number of hours, minutes, or seconds of logs that are retrieved with the `--hours`, `--minutes`, and `--seconds` options respectively. You can also follow the logs with the `-f` option. When using `-f` all logs will be printed to the console within the given time frame as well as any new logs that are sent to the logging Dashboard for the duration of the command. When using the `-f` option, hit ctrl-c to stop. Here is a sample command

```
catalyze logs -f --hours=6 --minutes=30
```

# Metrics

The `metrics` command gives access to environment metrics or individual service metrics through a variety of formats. This is useful for checking on the status and performance of your application or environment as a whole. The metrics command cannot be run directly but has sub commands.

# Metrics CPU

```
Usage: catalyze metrics cpu [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment CPU metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.
```

`metrics cpu` prints out CPU metrics for your environment or individual services. You can print out metrics in csv, json, plain text, or spark lines format. If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. You can only stream metrics using plain text or spark lines formats. To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands

```
catalyze metrics cpu
catalyze metrics cpu app01 --stream
catalyze metrics cpu --json
catalyze metrics cpu db01 --csv -m 60
```

# Metrics Memory

```
Usage: catalyze metrics memory [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment memory metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.
```

`metrics memory` prints out memory metrics for your environment or individual services. You can print out metrics in csv, json, plain text, or spark lines format. If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. You can only stream metrics using plain text or spark lines formats. To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands

```
catalyze metrics memory
catalyze metrics memory app01 --stream
catalyze metrics memory --json
catalyze metrics memory db01 --csv -m 60
```

# Metrics Network-In

```
Usage: catalyze metrics network-in [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment received network data metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.
```

`metrics network-in` prints out received network metrics for your environment or individual services. You can print out metrics in csv, json, plain text, or spark lines format. If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. You can only stream metrics using plain text or spark lines formats. To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands

```
catalyze metrics network-in
catalyze metrics network-in app01 --stream
catalyze metrics network-in --json
catalyze metrics network-in db01 --csv -m 60
```

# Metrics Network-Out

```
Usage: catalyze metrics network-out [SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]

Print service and environment transmitted network data metrics in your local time zone

Arguments:
  SERVICE_NAME=""   The name of the service to print metrics for

Options:
  --json=false     Output the data as json
  --csv=false      Output the data as csv
  --spark=false    Output the data using spark lines
  --stream=false   Repeat calls once per minute until this process is interrupted.
  -m, --mins=1     How many minutes worth of metrics to retrieve.
```

`metrics network-out` prints out transmitted network metrics for your environment or individual services. You can print out metrics in csv, json, plain text, or spark lines format. If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. You can only stream metrics using plain text or spark lines formats. To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands

```
catalyze metrics network-out
catalyze metrics network-out app01 --stream
catalyze metrics network-out --json
catalyze metrics network-out db01 --csv -m 60
```

# Rake

```
Usage: catalyze rake TASK_NAME

Execute a rake task

Arguments:
  TASK_NAME=""   The name of the rake task to run
```

`rake` executes a rake task by its name asynchronously. Once executed, the output of the task can be seen through your logging Dashboard. Here is a sample command

```
catalyze rake db:migrate
```

# Redeploy

```
Usage: catalyze redeploy SERVICE_NAME

Redeploy a service without having to do a git push

Arguments:
  SERVICE_NAME=""   The name of the service to redeploy (i.e. 'app01')
```

`redeploy` restarts a service without having to perform a code push. Typically when you want to update your code service you make a code change, git commit, then `git push catalyze master`. After the build finishes and a couple minutes later your code service will be redeployed. With the redeploy command, you skip the git push and the build. Here is a sample command

```
catalyze redeploy app01
```

# Services

```
Usage: catalyze services

List all services for your environment
```

`services` prints out a list of all services in your environment. The names of the services will be printed regardless of their currently running state. To see which services are currently running and which are not, use the [status](https://resources.catalyze.io/paas/cli/sections/status/) command. Here is a sample command

```
catalyze services
```

# Sites

The `sites` command gives access to hostname and SSL certificate usage for public facing services. `sites` are different from `certs` in that `sites` use an instance of a `cert` and are associated with a single service. `certs` can be used by multiple sites. The sites command can not be run directly but has sub commands.

# Sites Create

```
Usage: catalyze sites create NAME SERVICE_NAME HOSTNAME

Create a new site linking it to an existing cert instance

Arguments:
  NAME=""           The name of the site to be created. This will be used in this site's nginx configuration file
  SERVICE_NAME=""   The name of the service to add this site configuration to (i.e. 'app01')
  HOSTNAME=""       The hostname used in the creation of a certs instance with the 'certs' command
```

`sites create` allows you to create a site configuration that is tied to a single service. To create a site, you must first [create a cert](https://resources.catalyze.io/paas/cli/sections/certs-create/). A site has three pieces of information, a name, the service it's tied to, and the cert instance it will used. The name is the `server_name` that will be injected into this site's nginx configuration file. The service is a code service that will use this site configuration. Lastly, the cert instance must be specified by the `HOSTNAME` argument used in the [certs create](https://resources.catalyze.io/paas/cli/sections/certs-create/) command. Here is a sample command

```
catalyze sites create mywebsite.com app01 mywebsite.com
```

# Sites List

```
Usage: catalyze sites list

List details for all site configurations
```

`sites list` lists all sites for the given environment. The names printed out can be used in the other sites commands. Here is a sample command

```
catalyze sites list
```

# Sites RM

```
Usage: catalyze sites rm NAME

Remove a site configuration

Arguments:
  NAME=""      The name of the site configuration to delete
```

`sites rm` allows you to remove a site by name. Since sites cannot be updated, if you want to change the name of a site, you must `rm` the site and then [create](https://resources.catalyze.io/paas/cli/sections/sites-create/) it again. If you simply need to update your SSL certificates, you can use the [certs update](https://resources.catalyze.io/paas/cli/sections/certs-update/) command on the cert instance used by the site in question. Here is a sample command

```
catalyze sites rm mywebsite.com
```

# Sites Show

```
Usage: catalyze sites show NAME

Shows the details for a given site

Arguments:
  NAME=""      The name of the site configuration to show
```

`sites show` will print out detailed information for a single site. The name of the site can be found from the [sites list](https://resources.catalyze.io/paas/cli/sections/sites-list/) command. Here is a sample command

```
catalyze sites show mywebsite.com
```

# SSL

The `ssl` command offers access to subcommands that deal with SSL certificates. You cannot run the SSL command directly but must call a subcommand.

# SSL Resolve

```
Usage: catalyze ssl resolve CHAIN PRIVATE_KEY HOSTNAME [OUTPUT] [-f]

Verify that an SSL certificate is signed by a valid CA and attempt to resolve any incomplete certificate chains that are found

Arguments:
  CHAIN=""         The path to your full certificate chain in PEM format
  PRIVATE_KEY=""   The path to your private key in PEM format
  HOSTNAME=""      The hostname that should match your certificate (i.e. "*.catalyze.io")
  OUTPUT=""        The path of a file to save your properly resolved certificate chain (defaults to STDOUT)

Options:
  -f, --force=false   If an output file is specified and already exists, setting force to true will overwrite the existing output file
```

`ssl resolve` is a tool that will attempt to fix invalid SSL certificates chains. A well formatted SSL certificate will include your certificate, intermediate certificates, and root certificates. It should follow this format

```
-----BEGIN CERTIFICATE-----
<Your SSL certificate here>
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
<One or more intermediate certificates here>
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
<Root CA here>
-----END CERTIFICATE-----
```

If your certificate only includes your own certificate, such as the following format shows

```
-----BEGIN CERTIFICATE-----
<Your SSL certificate here>
-----END CERTIFICATE-----
```

then the SSL resolve command will attempt to resolve this by downloading public intermediate certificates and root certificates. A general rule of thumb is, if your certificate passes the `ssl resolve` check, it will almost always work on the Catalyze platform. You can specify where to save the updated chain or omit the `OUTPUT` argument to print it to STDOUT.

Please note you all certificates and private keys should be in PEM format. You cannot use self signed certificates with this command as they cannot be resolved as they are not signed by a valid CA. Here are some sample commands

```
catalyze ssl resolve ~/mysites_cert.pem ~/mysites_key.key *.mysite.com ~/updated_mysites_cert.pem -f
catalyze ssl resolve ~/mysites_cert.pem ~/mysites_key.key *.mysite.com
```

# SSL Verify

```
Usage: catalyze ssl verify CHAIN PRIVATE_KEY HOSTNAME [-s]

Verify whether a certificate chain is complete and if it matches the given private key

Arguments:
  CHAIN=""         The path to your full certificate chain in PEM format
  PRIVATE_KEY=""   The path to your private key in PEM format
  HOSTNAME=""      The hostname that should match your certificate (i.e. "*.catalyze.io")

Options:
  -s, --self-signed=false   Whether or not the certificate is self signed. If set, chain verification is skipped
```

`ssl verify` will tell you if your SSL certificate and private key are properly formatted for use with the Catalyze PaaS. Before uploading a certificate to Catalyze you should verify it creates a full chain and matches the given private key with this command. Both your chain and private key should be **unencrypted** and in **PEM** format. The private key is the only key in the key file. However, for the chain, you should include your SSL certificate, intermediate certificates, and root certificate in the following order and format.

```
-----BEGIN CERTIFICATE-----
<Your SSL certificate here>
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
<One or more intermediate certificates here>
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
<Root CA here>
-----END CERTIFICATE-----
```

This command also requires you to specify the hostname that you are using the SSL certificate for in order to verify that the hostname matches what is in the chain. If it is a wildcard certificate, your hostname would be in the following format: `*.catalyze.io`. This command will verify a complete chain can be made from your certificate down through the intermediate certificates all the way to a root certificate that you have given or one found in your system.

You can also use this command to verify self-signed certificates match a given private key. To do so, add the `-s` option which will skip verifying the certificate to root chain and just tell you if your certificate matches your private key. Please note that the empty quotes are required for checking self signed certificates. This is the required parameter HOSTNAME which is ignored when checking self signed certificates. Here are some sample commands

```
catalyze ssl verify ./catalyze.crt ./catalyze.key *.catalyze.io
catalyze ssl verify ~/self-signed.crt ~/self-signed.key "" -s
```

# Status

```
Usage: catalyze status

Get quick readout of the current status of your associated environment and all of its services
```

`status` will give a quick readout of your environment's health. This includes your environment name, environment ID, and for each service the name, size, build status, deploy status, and service ID. Here is a sample command

```
catalyze status
```

# Support-Ids

```
Usage: catalyze support-ids

Print out various IDs related to your associated environment to be used when contacting Catalyze support
```

`support-ids` is helpful when contacting Catalyze support by sending an email to support@catalyze.io. If you are having an issue with a CLI command or anything with your environment, it is helpful to run this command and copy the output into the initial correspondence with a Catalyze engineer. This will help Catalyze identify the environment faster and help come to resolution faster. Here is a sample command

```
catalyze support-ids
```

# Update

```
Usage: catalyze update

Checks for available updates and updates the CLI if a new update is available
```

`update` is a shortcut to update your CLI instantly. If a newer version of the CLI is available, it will be downloaded and installed automatically. This is used when you want to apply an update before the CLI automatically applies it on its own. Here is a sample command

```
catalyze update
```

# Users

The `users` command allows you to manage who has access to your environment through the organization that owns the environment. The users command can not be run directly but has three sub commands.

# Users List

```
Usage: catalyze users list

List all users who have access to the given organization
```

`users list` shows every user that belongs to your environment's organization. Users who belong to your environment's organization may access to your environment's services and data depending on their role in the organization. Here is a sample command

```
catalyze users list
```

# Users RM

```
Usage: catalyze users rm USER_ID

Revoke access to the given organization for the given user

Arguments:
  USER_ID=""   The Users ID to revoke access from for the given organization
```

`users rm` revokes a users access to your environment's organization. This is the opposite of the [users add](https://resources.catalyze.io/paas/cli/sections/users-add/) command. Revoking a user's access to your environment's organization will revoke their access to your environment. Here is a sample command

```
catalyze users rm 774bf982-fc4a-428b-a048-c38cffb7d0ab
```

The `vars` command allows you to manage environment variables for your code services. The vars command can not be run directly but has sub commands.

# Vars List

```
Usage: catalyze vars list

List all environment variables
```

`vars list` prints out all known environment variables for the associated code service. Here is a sample command

```
catalyze vars list
```

# Vars Set

```
Usage: catalyze vars set -v...

Set one or more new environment variables or update the values of existing ones

Options:
  -v, --variable    The env variable to set or update in the form "<key>=<value>"
```

`vars set` allows you to add new environment variables or update the value of an existing environment variable on your code service. You can set/update 1 or more environment variables at a time with this command by repeating the `-v` option multiple times. Once new environment variables are added or values updated, a [redeploy](https://resources.catalyze.io/paas/cli/sections/redeploy/) is required for your code service to have access to the new values. The environment variables must be of the form `<key>=<value>`. Here is a sample command

```
catalyze vars set -v AWS_ACCESS_KEY_ID=1234 -v AWS_SECRET_ACCESS_KEY=5678
```

# Vars Unset

```
Usage: catalyze vars unset VARIABLE

Unset (delete) an existing environment variable

Arguments:
  VARIABLE=""   The name of the environment variable to unset
```

`vars unset` removes an environment variables from your associated code service. Only the environment variable name is required to unset. Once environment variables are unset, a [redeploy](https://resources.catalyze.io/paas/cli/sections/redeploy/) is required for your code service to realize the variable was removed. Here is a sample command

```
catalyze vars unset AWS_ACCESS_KEY_ID
```

# Version

```
Usage: catalyze version

Output the version and quit
```

`version` prints out the current CLI version as well as the architecture it was built for (64-bit or 32-bit). This is useful to see if you have the latest version of the CLI and when working with Catalyze support engineers. Here is a sample command

```
catalyze version
```

# Whoami

```
Usage: catalyze whoami

Retrieve your user ID
```

`whoami` prints out the currently logged in user's users ID. This is used with the [users add](https://resources.catalyze.io/paas/cli/sections/users-add/) and [users rm](https://resources.catalyze.io/paas/cli/sections/users-rm/) commands as well as with Catalyze support engineers. Here is a sample command

```
catalyze whoami
```

# Worker

```
Usage: catalyze worker TARGET

Start a background worker

Arguments:
  TARGET=""    The name of the Procfile target to invoke as a worker
```

`worker` starts a background worker asynchronously. The `TARGET` argument must be specified in your `Procfile`. Once the worker is started, any output can be found in your logging Dashboard or using the [logs](https://resources.catalyze.io/paas/cli/sections/logs/) command. Here is a sample command

```
catalyze worker web
```

