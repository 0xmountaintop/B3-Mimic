# Stratum Mining Protocol

## Login

__SEND:__

```
{
    "method": "login",
    "params": {
        "login": "antminer_1",
        "pass": "123",
        "agent": "bmminer/2.0.0"
    },
    "id": 1
}
```

__RECVD:__

```
{
    "id":1,
    "jsonrpc":"2.0",
    "result":{
        "id":"antminer_1",
        "job":{
            "version":"0100000000000000",
            "height":"2315000000000000",
            "previous_block_hash":"ecaeb84f7787aca9ed08199169fe8dc5a01f3ba50c90c6844452f916d645d911",
            "timestamp":"239dc75a00000000",
            "transactions_merkle_root":"178f3dfaf916a5f8167100602254df50f6821c243a1f6263efedde798e9271a2",
            "transaction_status_hash":"c9c377e5192668bc0a367e4a4764f11e7c725ecced1d7b6a492974fab1b6d5bc",
            "nonce":"0c04000001000000",
            "bits":"4690890000000021",
            "job_id":"16368",
            "seed":"a2a62d7715ee2234e1d73c22d26a1707fb7bc0f4ee0c01d43a4c97b0328379c5",
            "target":"c5a70000"
        },
        "status":"OK"
    },
    "error":null
}
```

## Job Notification

__RECVD:__

```
{
    "jsonrpc":"2.0",
    "method":"job",
    "params":{
        "version":"0100000000000000",
        "height":"8c44000000000000",
        "previous_block_hash":"8f11dad3d32550bfb5a6bf9f7a344c8a436be37897887e3c74d1b380abd3587d",
        "timestamp":"4c3afa5a00000000",
        "transactions_merkle_root":"2c0abd8189aba80a8940c3bb4263c9560a97e81d89aa1389ff47e754add4ce98",
        "transaction_status_hash":"c9c377e5192668bc0a367e4a4764f11e7c725ecced1d7b6a492974fab1b6d5bc",
        "nonce":"0000000037780000",
        "bits":"1d7a04000000001d",
        "job_id":"171567",
        "seed":"6aa116ef17d69e398ba388e8f0e1805f6f98afb9dfb8fb87679ad5d3c099099d",
        "target":"ffff0000"
    }
}
```


## Submission

__SEND:__

```
{
    "method": "submit", 
    "params": {
        "id": "antminer_1", 
        "job_id": "4171", 
        "nonce": "bc000d41", 
        "result": "7f7bcc61373e63c5a97f5bfd890411ef1bd914ba586ad02acf881c771b000000"
    }, 
    "id":3
}                    
```

__RECVD:__ 

_succeed_

```
{
    "id":3,
    "jsonrpc":"2.0",
    "result":{
        "status":"OK"
    },
    "error":null
}
```

_fail_

```
{
    "id":41,
    "jsonrpc":"2.0",
    "result":null,
    "error":{
        "code":-1,
        "message":"Low difficulty share"
    }
}
```

```
{
    "id":12,
    "jsonrpc":"2.0",
    "result":null,
    "error":{
        "code":-1,
        "message":"Block expired"
    }
}
```
