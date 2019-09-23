# tk
* create &amp; manage tickets
* works only with macOS & linux
* you need sublime installed and in your `$PATH`


## status Values:
* active = workingON + waitingFor
* workingON
* waitingFor
* closed


## usage examples
```sh
## CREATE NEW TICKET
tk new _TICKET
tk n   _TICKET

## LIST ACTIVE
tk list
tk l

## LIST TICKESTS WITH SPECIFIC STATUS
tk list _STATUS
tk l    _STATUS

## LIST ALL TICKESTS
tk list all
tk l    all

## CHANGE STATUS (active --> closed OR closed --> active (workingON))
tk change _TICKET 
tk c      _TICKET

## CHANGE STATUS (existing status --> specified status)
tk c _TICKET _STATUS
tk c _TICKET _STATUS

## CHANGE NOTES
tk change _TICKET -n "_NOTES"
tk c      _TICKET -n "_NOTES"

## REMOVE DB ENTRIES NOT THE TICKET NOTES
tk r _TICKET
tk r all

## WORK TICKET
tk w _TICKET
```

## next
* needs refactoring:
    - if db does not exists list returns error
    - input validation
* aditional fields to the Ticket struct:
    - rAccount
    - aAccount
    - db notes to the tiket to describe the ticket main objective
