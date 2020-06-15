# deployer

deployer component with http endpoint for webhooks

## Content
```
.           (deployer codebase)
├── api     (deployer api)
├── appctl  (kube- and legacyctl wrapper)
├── config  (commons for configs)
├── model   (deployer model)
├── test    (integration tests)
└── util    (utilities)
```

## Usage 

Run the deployer server
```
go build && ./bsc-deployer
```

Then use the client to trigger a deployment
```
./client.sh
```

Or stop the processes
```
./client.sh stop
```