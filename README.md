## How does it works?
User send a request to /job endpoint, a gouroutine assigns to user job and will run based on user defined time interval. There is a goroutine limit to prevent app crash. You can edit it on .env 
## Run with docker:

```bash
docker compose up -d
```

```bash
endpoints are available at:
http://0.0.0.0:8080/
```
## Run bare metal:
Run the following command in the project root:
```bash
cd cmd
go run main.go
```


**Note:**
# License
[![Licence](https://img.shields.io/github/license/Ileriayo/markdown-badges?style=for-the-badge)](./LICENSE)
