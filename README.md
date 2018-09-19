[![Build Status](https://travis-ci.org/vvval/go-metadata-scanner.svg?branch=master)](https://travis-ci.org/vvval/go-metadata-scanner)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/vvval/go-metadata-scanner/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/vvval/go-metadata-scanner/?branch=master)
[![codecov](https://codecov.io/gh/vvval/go-metadata-scanner/branch/master/graph/badge.svg)](https://codecov.io/gh/vvval/go-metadata-scanner)

# go-metadata-scanner
Simple app with ability to scan images metadata into output file or read metadata from input file and bulk update the photos

# scan example:
```go
go run main.go scan -d folder/with/images/ -o file-inside-scanned-directory -f format
```

Supported formats are `csv`, `json` and special `mscsv` format, `csv` is a default value and can be omitted

# write example:
```go
go run main.go write -f folder/with/image/caption/file.csv
```

If the csv file uses as a separator not-comma char, do not forget to specify it with a `s` flag:
```go
go run main.go write -f folder/with/image/caption/file.csv -s ;
```

By default, given file directory is used as a root for searching files, but it can be specified using `d` flag:
```go
go run main.go write -f folder/with/image/caption/file.csv -d folder/with/images/
```

By default, data found in the csv file will overwrite the existing data in the images and the images will be updated.
As an option, original images will be untouched (the originals will be renamed using `_original` postfix) or list-type metadata will be appended to existing data. In the first case, use `o` flag, meaning "save originals"; for the second case, use `a`, meaning "append metadata".

All flags can be used at the same time:

```go
go run main.go write -f folder/with/image/caption/file.csv -s ; -d folder/with/images/ -a -o
```

#### Files match
Files are searched using soft comparison of given filename in the caption file as endings for real files, 
prefix can be any letter, then `_` char and `0` digit (none or any amount of).
Also program tries to match the extension.
If program will find several candidates with no strict match, such file will be omitted. 
Strict match means that full path matches the candidate (with or without proposed extension)
Examples:
 
 * `1.png` will be searched as `0001.png` and `IMG_001.png`.
 * `2` will be searched as `2.png` and `2.jpg`.
 * Even if you specify `3.jpg` and have image named `3.jpg.png`, we will find it. 
 * If there files `002.png` and `0002.png` and the filename is specified as `2.png` - this line will be skipped
 * If there files `2.png`, `002.png` and `02.png` and the filename is specified as `2.png` - the `2.png` file will be used (as it strictly matched)

#### Files destination
Images should be placed relative the caption file (as the filename).
For example, if filename is:
 * `test.jpg`, image should be in the same folder as the caption file;
 * `subfolder/test.jpg` image should be in ` folder/with/image/caption/subfolder/`;
 * `../another/folder/test.jpg` image should be in ` folder/with/image/another/folder/`;
