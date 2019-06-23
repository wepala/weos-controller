# README #

The Controller is meant to handle incoming requests and route to the appropriate business login

### What is this repository for? ###

* Version: 0.1.0

### How do I get set up? ###

You can run the serve command to give access via http. There are a few ways to start the server

1. Use cli parameters `weos-controller serve http-mock -a site-api.yml -c site-config.yml`
1. Use environment variables set environment variable `API_YAML` and `CONFIG_YAML` and then start the server `weos-controller serve http-mock`
1. Configure parameters in a config file `weos-controller serve http-mock -c weoscontroller.yml`
1. Place a config in the home folder of the service `weos-controller serve http-mock`

### Contribution guidelines ###

This project uses [gitflow workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)

* Clone the repo to local
* Create feature branch from dev branch (e.g. feature/WEOS-1)
* Push the feature branch to the remote
* Create PR from  the ticket branch to dev branch 
* When the item is merged to master it will be deployed

To aid with this use the git flow cli (you will be able to create feature branches e.g. git flow feature start APO-1)

#### New Features ####





### Who do I talk to? ###

* Admin - Akeem Philbert <akeem.philbert@wepala.com>