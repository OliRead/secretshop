# Secret Shop
## Ho ho, you found me!
Welcome to the Secret Shop, a light weight API built in Go for pulling item
purchase information from Dota 2 Demo files written by HonestAbe and powered
by Manta, written by the beautiful fellows over at Dotabuff

## Quick Start
### Docker Compose
The quickest way do get started is with Docker Compose, it's as simple as running

```sh
cd [secretshop-dir]/cmd
docker-compose up .
```

This will start a Maria DB container and a Secret Shop container exposing its API
on port 8080. You can check that Secret Shop has started, and is running by checking
the logs of the Secret Shop container, or by heading over to localhost:8080/status

### Uploading Replays
Secret Shop expects to recieve a multipart form request to the endpoint /replay/upload.
By default there is no authentication enabled on the API, meaning anybody can upload
and modify information about a replay. **If you are intending to use this in a public
environment ensure you add some authentication**

The easiest way to upload a replay to Secret Shop is by using cURL
``` sh
curl -F replay=@/path/to/replay.dem localhost:8080
```
this will upload your replay, parse it and store it in the database. Secret Shop has 
fairly verbose logging on both http requests and the std output. If any errors occour 
you should be able to see them in both the http response and the docker logs.

### API Documentation
Full API Documentation is available at [docs.honestabe.co.uk/secretshop](https://docs.honestabe.co.uk/secretshop)

### Special Thanks
Special thanks go to [Dotabuff team](https://www.dotabuff.com/) and the [Manta project](https://github.com/dotabuff/manta), without them this would have 
taken significantly longer to build.

Bonus thanks to [Conor Clafferty](https://github.com/cclafferty) for being my rubber duck when it comes to Go projects.
