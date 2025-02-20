# bulk-mkv-renamer
Renames files ripped from makemkv

# How to use
Build
```bash
go build .\cmd\bulk_rename.go
```
```bash
./bulk_rename bfr -p "/SomeFilePath/COWBOY_BEBOP" -c "Cowboy Bebop" -f "snake_case"
```
Renames files in a given directory.
This command takes three flags
path: 'File path where your .mkv files live'
content-name: 'The content of your mkv files are your directory'
file-name-case: The case for file names (default or camelcase)
This will then go though the files and ask for input to rename what file to what ep or movie

| Flag                     | Type   | Description                                                                  |
|--------------------------|--------|------------------------------------------------------------------------------|
| `-c`, `--content-name`   | string | Specify the name of the movie or TV show                                     |
| `-f`, `--file-name-case` | string | Specify the case for file names (default: "White space naming" "snake_case") |
| `-h`, `--help`           |        | Help for bfr                                                                 |
| `-p`, `--path`           | string | Specify the path to the directory                                            |
