# Cloud-compactor

Tool that automatically compact videos from Cloud storage:

```mermaidjs
sequenceDiagram
    Cloud Compactor ->> Cloud Storage: Request recursive list of files
    activate Cloud Storage
    Cloud Storage -->> Cloud Compactor: Retrieve list of files
    deactivate Cloud Storage
    Cloud Compactor ->> Cloud Compactor: Filter only videos that have not been processed
    loop For each file from the list
    Cloud Compactor ->> Cloud Storage: Request video download
    activate Cloud Storage
    Cloud Storage -->> Cloud Compactor: Download video file
    deactivate Cloud Storage
    Cloud Compactor ->> FFMPEG: Request processing
    FFMPEG ->> FFMPEG: Process video
    FFMPEG -->> Cloud Compactor: [Done]
    Cloud Compactor ->> Cloud Compactor: Delete downloaded file
    Cloud Compactor ->> Cloud Storage: Upload new re-encoded video file
    Cloud Compactor ->> Cloud Compactor: Delete re-encoded file
    Cloud Compactor ->> Cloud Storage: Delete original file
    Cloud Storage -->> Cloud Compactor: [Done]
    end
```

Supported access:
- FTPS

Supported output video formats:
- x265

## Requirements

- Docker
- Docker-compose

## Configure

Create a `config/config.yaml` file with the following content:

```yaml
# Path to files on server
path: /

# Speed of the video processing for the ffmpeg (slow, medium, fast, veryfast)
speed: veryfast

# Formats of the video files
formats:
  # Input formats accepted
  inputs:
    - "mkv"
    - "3gp"
    - "avi"
    - "m4v"
    - "mov"
    - "mp4"
    - "mpg"
    - "mts"
    - "ogm"
  # Processed suffix is the marker that the file is already processed
  processed_suffix: processed
  # Output format
  output: mp4

# FTP settings
ftp:
  # FTP server address (without ftp://)
  address:
  # FTP username
  username:
  # FTP password
  password:
```

## Usage

All you have to do is run the following command:

```bash
make run
```

It will build the image locally and run the container.
