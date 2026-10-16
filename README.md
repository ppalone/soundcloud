# Soundcloud

☁️ Unofficial Soundcloud Go client

## Installation

```bash
go get github.com/ppalone/soundcloud
```

## Usage

### Initialize Client

```go
client, err := soundcloud.NewClient()
if err != nil {
    log.Fatal(err)
}
```

### Search Tracks

```go
client, err := soundcloud.NewClient()
if err != nil {
    log.Fatal(err)
}

tracks, err := client.SearchTracks(context.Background(), "test")
if err != nil {
    log.Fatal(err)
}

for _, track := range tracks {
    fmt.Println(track.Title)
}
```

### Search Track by ID

```go
client, err := soundcloud.NewClient()
if err != nil {
    log.Fatal(err)
}

track, err := client.GetTrackById(context.Background(), 123456789)
if err != nil {
    log.Fatal(err)
}

fmt.Println(track.Title)
```

### Get Stream by ID

```go
client, err := soundcloud.NewClient()
if err != nil {
    log.Fatal(err)
}

stream, err := client.GetStreamById(context.Background(), 123456789)
if err != nil {
    log.Fatal(err)
}

fmt.Println(stream)
```

## Author

Pranjal
