## Simple Go Blockchain

Following the tutorial series from [Coral Health](https://mycoralhealth.medium.com/part-2-networking-code-your-own-blockchain-in-less-than-200-lines-of-go-17fe1dad46e1), with more organization and unit tests ;)

To test, try running:
    
    go main.go

And in a different terminal:
    
    curl localhost:8080

To see the genesis blockchain printed to that terminal output.

To add another block:

    curl localhost:8080 -X POST -d '{"BPM":58}'

To change the port, change the PORT value in the .env file.