# Implement the Transport Logging

In this exercise you are required to implement the transport logging for the new endpoints, exposing when the endpoint is called 

## Your would be required to:
* Implement a logging middleware endpoint to handle the wholesale requests
* Create and return an endpoint to serve as the middleware
* Log the TotalWholesalePriceEndpoint message when it enters, and exits using a defer statement
* Call the next item in the chain
* Add the logging endpoint to the wholesaleEndpoint chain when creating the http handler

## Important to pass the validation checks:
* Do not use parameter type sharing
* Use the naming conventions as per the code and examples
