# Error Handling
This document goes over error handling in the `State > Sequencer > Orchestrator > DispatcherIface` chain.

## Regular Errors

 - Treated as a test failure
 - Error will be logged
 - Error will be submitted to the result processor
 - Sequence will continue to run as long as the `ContinueOnFail` function for given state is true

## Fatal Errors 

- Treated as a test failure
- Error will be logged
- Error will be submitted to the result processor
- Will NOT continue to run the sequence
- Orchestrator goes into a `FatalError` state
- Any dispatcher must send the `RecoverFromFatal` signal to the orchestrator to return to idle