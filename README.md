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
tk n _TICKET

## LIST ALL ACTIVE
tk l

## LIST SPECIFIC STATUS
tk l -s _STATUS

## CHANGE STATUS
#### active --> closed
#### closed --> workingON
tk c _TICKET

## CHANGE STATUS
#### active --> closed
#### closed --> workingON
tk c -s _STATUS

## REMOVE DB ENTRIES NOT THE TICKET NOTES
tk r _TICKET
tk r all

## WORK TICKET
tk w _TICKET
```

## next
* needs refactoring:
    - input validation
* aditional fields to the Ticket struct:
    - rAccount
    - aAccount
    - db notes to the tiket to describe the ticket main objective