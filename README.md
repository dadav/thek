# thek

<img src="./img/thek.png" width="300">

**thek** is a small tool to schedule recordings from the german mediathek.

## install

Run `make build && sudo make install`

## usage

Adjust the `config.yaml` to your needs and run the binary.

The binary has some these flags:

```bash
$ thek -h
  -config string
        path to the config yaml file (default "config.yaml")
  -debug
        enable debug mode
  -json
        enable json output
  -version
        just print the version
```

## configuration

**Example**

```yaml
---
defaults:
  output_directory: output
  safety_duration: 5m

stations:
  3sat: "http://zdf-hls-18.akamaized.net/hls/live/2016501/dach/high/master.m3u8"
  arte: "https://artesimulcast.akamaized.net/hls/live/2030993/artelive_de/index.m3u8"
  zdf: "http://zdf-hls-15.akamaized.net/hls/live/2016498/de/high/master.m3u8"
  zdfinfo: "http://zdf-hls-17.akamaized.net/hls/live/2016500/de/high/master.m3u8"
  zdfneo: "http://zdf-hls-16.akamaized.net/hls/live/2016499/de/high/master.m3u8"
  phoenix: "http://zdf-hls-19.akamaized.net/hls/live/2016502/de/high/master.m3u8"
  ki.ka: "https://kikageohls.akamaized.net/hls/live/2022693/livetvkika_de/master.m3u8"

recording_tasks:
  - station: zdfneo
    show_keywords: psych
    output_directory: ./output/psych
    safety_duration: 5m
  - show_keywords: test
    output_directory: output/test
    safety_duration: 5m
```
