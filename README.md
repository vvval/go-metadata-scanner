# go-metadata-scanner
Simple app with ability to scan images metadata into output file or read metadata from input file and bulk update the photos

# scan example:
```go
go run main.go scan -d folder/with/images/ -o file-inside-scanned-directory -f format
```

Supported formats are `csv`, `json` and special `mscsv` format

# write example:
```go
go run main.go write -f folder/with/image/caption/file.csv
```

#### Files match
Files are searched using soft comparison of given filename in the caption file as endings for real files, 
prefix can be any letter, then `_` char and `0` digit (none or any amount of).
Also program tries to match the extension. If program will find several candidates with no strict match, such file will be omitted.
Examples:
 
 * `1.png` will be searched as `0001.png` and `IMG_001.png`.
 * `2` will be searched as `2.png` and `2.jpg`.
 * Even if you specify `3.jpg` and have image named `3.jpg.png`, we will find it. 

#### Files destination
Images should be placed relative the caption file (as the filename).

For example, if filename is:
 * `test.jpg`, image should be in the same folder as the caption file;
 * `subfolder/test.jpg` image should be in ` folder/with/image/caption/subfolder`;
 * `../another/folder/test.jpg` image should be in ` folder/with/image/another/folder`;
