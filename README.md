# shissue

![travis build](https://api.travis-ci.org/arthurmco/shissue.svg?branch=master)

shissue (`shell issue`) is a software that allows you to manage git project issues 
from the command line, without the need to open up a browser and go to your issue.

> Browsers are huge beasts that can be slow sometimes.
> 
> Command line is small and fast.

Don't having a command line tool for issues made me write less issues for
my projects, both public and private ones. 

We have a command line tool for managing git source history, why not one for
managing git issue history?

## Build & Install

 - Install the go compiler (`pacman -S go`, `apt install golang`, `yum install golang`, one of the three will do)
 - Set the environment variable `GOPATH` to where you want to download
   the source
 - Type `go get github.com/arthurmco/shissue` and wait ~and make a
   coffee~. A folder named `github.com` will appear at
   `GOPATH`. Browse it until you find `shissue`folder.
 - Open the shissue folder
 - Run the build.sh script with the folder you want to install shissue as the argument (like `sh build.sh /usr/local/bin`). **Note that, for some directories, you need to run the script as root!**
 
## Usage

```
> $ shissue help

 shissue - view github/gitlab issues in command line

 Usage: shissue [options] command [commandargs...]

 Commands: 
	help                 Print this help text
	issues               List repository issues

 Options: 
 [-U|--username] <<username>>
	specify the username used in your repo account
 [-P|--password] <<password>>
	specify the password used in your repo account
 --allow-untrusted-certs
	Allow connecting to certificates not trusted by the system

```

* **issues** will list the issues from the current repository, if it does
 have a compatible remote.  
 For now, it only  supports Github public repos, but more will be added 
 over time (I *do* have  projects in other sites, too!). 
 
 * You can specify only 'username'. If you do that, the software will ask for the 
   password.

 * You can also store the username inside git configuration (using `git config`). Use `git config shissue.username <<username>>` for storing the github username, and you won't have to type it.
 
 * In Gitlab, you have the [personal access token](https://docs.gitlab.com/ce/user/profile/personal_access_tokens.html) for accessing repos without 
   needing a password. Use `git config shissue.token <<token>>` to set it
   inside shissue.
   

To see a video of shissue in action, check the video below:

[![asciicast](https://asciinema.org/a/qDxWdqzvO5VLnBlpOTdnNz1Im.png)](https://asciinema.org/a/qDxWdqzvO5VLnBlpOTdnNz1Im)

## What do we have?

What it already supports is bold, what it doesn't is not

Support for *everything* in this list is planned, so don't worry! :smile:

 - Support reading issues from
   - **Github public repos**
   - **Github private repos** (Maybe? Need to check. I don't have private repos)
   - **Gitlab public & private repos**
   - Bitbucket public & private repos
   
 - Support for creating issues
   - on Github
   - on Gitlab
   - on Bitbucket
   
 - Support for editing issues
 - Support for viewing pull requests
 - Support for viewing issues' and PRs comments
 - Support for commenting on issues & PRs
 
( I might add support for that reaction thing in github issue system)

When everything above is implemented, I'll launch an 1.0 version.

## Contributing

I do want your contribution. Don't be afraid :smile:

Issues and pull requests can be done in Portuguese or English. 
I don't have preference. 

What I need the most is someone to help building a test suite. I don't know
how to test this without importunating the github API.

## Why don't `gissue` ?

https://github.com/search?utf8=%E2%9C%93&q=gissue&type=

https://github.com/search?utf8=%E2%9C%93&q=shissue&type=

Compare. :wink:

# Images

Everybody likes nice images, so...

![shissue @ node.js](https://i.imgur.com/Ui5uYmZ.png "shissue listing node.js open issues")
![shissue @ ourselves](https://i.imgur.com/0K5udPt.png "shissue listing our own issues")

## Licensing

Everything inside here is under the MIT license

