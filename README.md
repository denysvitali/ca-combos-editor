# ca-combos-editor

A combo editor for NV ITEM 00028874

## Setup
1. Clone this repo
2. Execute the following (in the repo's root)
```
go get github.com/denysvitali/ca-combos-editor/cmd
go run cmd/main.go -h
```

## Usage

### Parse 00028774's content
1. Extract it:
```
zlib-flate --uncompress < 00028774 > extracted.bin
```

2. Parse it:
```
go run cmd/main.go parse extracted.bin
```

### Create a 00028774 file based on a band file
  
1. Provide a bands.txt file in the format shown in `test/resources/2019-10-17/bands.txt` (one combo per line)  
2. `go run cmd/main.go create bands.txt 00028774_uncompressed`
3. Compress it: `./compress.sh 00028774_uncompressed`
4. Write the new 00028774 file to your modem


