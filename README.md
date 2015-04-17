# S3 Downloader
An easy CLI to download all contents from a given AWS S3 bucket

## Installation

```bash
go get github.com/supherman/s3_downloader
```

## Basic Usage

```bash
s3_downloader --bucket=my-bucket \
--access-key-id=xxxxxx \
--secret-access-key=xxxxx
```

## Specify Region

```bash
s3_downloader --bucket=my-bucket \
--access-key-id=xxxxxx \
--secret-access-key=xxxxx \
--region=us-west-2
```

## Configure donwload concurrency

```bash
s3_downloader --bucket=my-bucket \
--access-key-id=xxxxxx \
--secret-access-key=xxxxx \
--region=us-west-2 \
--concurrency=15
```
