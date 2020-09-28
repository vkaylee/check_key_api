# check_key_api
## Description
This tool is used to check your API with key and value. You can use it in your CI to check available service.
## Installation
Download binary file in releases tab. Execute it.

For Linux 64
    
    curl --request GET -sL \
        --url 'https://github.com/vleedev/check_key_api/releases/download/v0.1/check_key_api-linux-amd64'
        
## Use it
### Check equal
    ./check_key_api-linux-amd64 -url "your API url" -key "name" -value "value" -condition "equal"
### Check unequal
    ./check_key_api-linux-amd64 -url "your API url" -key "name" -value "value" -condition "unequal"
### The result
It will return exit code 0 when it's successful and other when it's failed.
