# spamalot coordinator (wip)

![dashboard](https://i.imgur.com/3QJx8qC.png)

## description
The spamalot-coo is a SPA allowing the remote control of
spamalot slave instances which in turn control spammers.

Slaves communicate with the coordinator via a HTTPS-API, secured
by API tokens per slave instance.

The goal of the coordinator is to leverage unused power to increase
the IOTA Tangle's overall throughput by spamming the network
with transactions after a best effort approach.

## todo

This project is pretty much in its infancy, basically everything
has to be done.

* [ ] spammers
    * [ ] API to slave instances
* [ ] slave instance
    * [ ] API to spammers
    * [ ] API to coordinator
* [ ] backend
    * [ ] data models
    * [ ] API to slave instances
    * [ ] authentication control
* [ ] single page application (SPA)
    * [ ] CRUD instances and their configuration
    * [ ] CRUD spammer configurations
    * [ ] responsive design
