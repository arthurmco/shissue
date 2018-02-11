# shissue

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
 - Create a folder tree named `src/github.com/arthurmco/`, or whatever 
   namespace this project might be in. This is important, because the go build 
   tool is... well...
 - Set the environment variable `GOPATH` to the top of that tree.
 - `git clone` this source inside that tree
 - Run the build.sh script with the folder you want to install shissue as the argument (like `sh build.sh /usr/local/bin`). **Note that, for some directories, you need to run the script as root!**
 
 Obs: I'm not a Go expert. I am angry, but it's because I need help.
 
 Obs 2: I feel like I need to do a build script.
 
## Usage

```
 shissue - view github issues in command line

 Usage: shissue [options] command [commandargs...]

 Commands: 
	help                 Print this help text
	issues               List repository issues

 Options: 
 [-U|--username] <<username>>
	specify the username used in your github account
 [-P|--password] <<password>>
	specify the password used in your github account

```

* **issues** will list the issues from the current repository, if it does
 have a compatible remote.  
 For now, it only  supports Github public repos, but more will be added 
 over time (I *do* have  projects in other sites, too!). 

## What do we have?

What it already supports is bold, what it doesn't is not

Support for *everything* in this list is planned, so don't worry! :smile:

 - Support reading issues from
   - **Github public repos**
   - **Github private repos** (Maybe? Need to check. I don't have private repos)
   - Gitlab public & private repos
   - Bitbucket public & private repos
   
 - Support for creating issues
   - on Github
   - on Gitlab
   - on Bitbucket
   
 - Support for writing issues
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

## Licensing

Everything inside here is under the MIT license

