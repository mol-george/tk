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
## AND/OR NOTES
tk change _TICKET [-n "_NOTES"]
tk c      _TICKET [-n "_NOTES"]

## CHANGE STATUS (existing status --> specified status)
## AND/OR NOTES
tk c _TICKET _STATUS [-n "_NOTES"]
tk c _TICKET _STATUS [-n "_NOTES"]

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
    - update script to take in account aditional sometimesSoon status
    - if db does not exists list returns error
    - input validation
* aditional fields to the Ticket struct:
    - rAccount
    - aAccount
    - db notes to the tiket to describe the ticket main objective

### change
* I must switch to using positional arguments for notes and I must validate positional arguments:
    - if 1 argument:
        + the 1st and only argument must be an existing ticket
    - if 2 arguments:
        + 1st must be an existing ticket
        + 2nd must comment
* I must lowercase all states
* the basic idea with change is if:
    - if new status is specified it changes to the specified status
    - if new status is not specified:
        + it switches from closed to working on
        + and between any active status to closed
* if I only want to change the note I have to specify the status in order to not change the status unintentionally
* have to implement search:
    - to list a ticket if I can't locate it
    - to display tickets between dates
* there is no way to just update the ticket only without switching the state

### rm
* does not work

### work
* should automatically set a ticket to working on??