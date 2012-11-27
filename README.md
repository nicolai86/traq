# traq

bash based time tracking using only files

## Installation

`traq` assumes you're installing it to your home directory, into `~/.traq`. You're all setup with these four commands:

```
$ mkdir ~/.traq
$ git clone git@github.com:nicolai86/traq.git ~/.traq/traq
$ echo "export PATH=$PATH:$HOME/.traq/traq" >> ~/.bash_profile
$ . ~/.bash_profile
$ which traq # => ~/.traq/traq/traq
```

## Tests

The project has some tests using [bats](https://github.com/sstephenson/bats). Assuming you got `bats` installed, run them using the following command:

```
cd $HOME/.traq/traq
bats traq_test.bats
```

## Usage

`$ traq foo`: creates an entry for #foo at todays date

`$ traq stop`: appends the special stop delimitor to todays file

`$ traq`: echo the content of todays file to stdout. If the file does not exist, nothing is echoed.

`$ traq -d 2012-07-30` echo the content of the file from the given date to stdout. If the file does not exist, nothing is echoed.

`$ traq -w 31` echo the content of all files from the calendar week 31 to stdout. If the week does not contain files, nothing is echoed.

## Hacking

All files are placed under

```
$ $HOME/.traq/timestamps/kw-<week number>/timestamps-<date>
```

or, if `-p <project>` was given, under

```
$ $HOME/.traq/<project>/kw-<week number>/timestamps-<date>
```

Each file can contain multiple lines of the following format:

```
<timestamp>;<tag>
```

Here's some sample content:

```
Thu Sep 27 07:05:05 +0400 2012;#foo
Thu Sep 27 07:15:05 +0400 2012;#bar
Thu Sep 27 07:25:05 +0400 2012;stop
```

## Helpers

To ease evaluation of traq-files `traq` comes with two helper scripts, `traqtrans` and `traqeval`.

`traqtrans` transforms the timestamp into a unix timestamp,
and `traqeval` sums up tags.

Just pipe the output and you'll get something like this:

```
$ traq -p test -w 39 | traqtrans | traqeval

2012-09-27
#foo:0.16666666666666666
#bar:0.16666666666666666
%%
```