# traq

time tracking using text files

## Usage examples

``` bash
# start time tracking for #project
$ traq project

# start time tracking for #project with the comment 'working on the landing page'
$ traq project working on the landing page

# stop time tracking
$ traq stop

# echo the content of todays file to stdout. If the file does not exist, nothing is echoed.
$ traq

# echo the content of the file from the given date to stdout. If the file does not exist, nothing is echoed.
$ traq -d 2012-07-30

# echo the content of all files in february of the current year
$ traq -m 2

# echo the content of all files in september 2012
$ traq -m 9 -y 2012

# starts time tracking for client-a on #development
$ traq -p client-a development

# stops time tracking for client-a
$ traq -p client-a stop

# list tracked times for client-a from today
$ traq -p client-a
```

## Evaluation

To evaluate traq files pass the `-e` command line flag to the utility.
E.g.

    $ traq -p test -e
    2012-09-27
    #foo:0.1666
    #bar:0.1666
    %%

## Installation

`traq` assumes you're installing it to your home directory, into `~/.traq`. This will set you up:

``` bash
$ mkdir $HOME/.traq
$ mkdir -p $HOME/Library/traq
$ git clone git@github.com:nicolai86/traq.git ~/.traq
$ echo "export TRAQ_PATH=$HOME/.traq" >> ~/.bash_profile
$ echo "export TRAQ_DATA_DIR=$HOME/Library/traq" >> ~/.bash_profile
$ echo "export PATH=$PATH:$HOME/.traq/bin" >> ~/.bash_profile
$ ln -s $HOME/.traq/man/traq.1 /usr/local/share/man/man1/traq.1
$ . ~/.bash_profile
$ which traq
```

### Build instructions

Make sure your `$GOPATH` is set up properly. Then link traq properly:

``` bash
$ mkdir -p $HOME/go/src
$ ln -s $HOME/.traq/src/ $HOME/go/src/traq
```

``` bash
$ go build -o bin/traq app.go
```

### Go Formatting

``` bash
$ gofmt -tabs=false -tabwidth=2 -w traq.go
```
## Bash Completion

If you have `bash-completion` installed you can setup bash completion for traq as well. This example assumes you are using [HomeBrew][1] and have `bash-completion` installed.

``` bash
$ ln -s $TRAQ_PATH/traq_completion.sh $(brew --prefix)/etc/bash_completion.d/traq
```

Ubuntu users can do the following:

``` bash
$ sudo apt-get install bash-completion
$ echo ". $TRAQ_PATH/traq_completion.sh" >> ~/.bash_profile
```

## Migration to v0.5

traq 0.5 has a different, flatter directory structure. Instead of one directory per year week of the year,
we now only have one directory per year.

The following bash script helps you migrate your data:

```bash
for directory in $(find $HOME/Library/traq -maxdepth 1 -mindepth 1 -type d); do
  echo $directory
  for year in $(find $directory -maxdepth 1 -mindepth 1 -type d); do
    for week in $(find $year -maxdepth 1 -mindepth 1 -type d); do
      for file in $(find $week -maxdepth 1 -mindepth 1 -type f); do
        cleanfile="${file//timestamps-/}"
        mv "$file" "$year/${cleanfile##*/}"
      done
      rm -fr $week
    done
  done
done
```

## Hacking

All files are placed under

    $ $TRAQ_DATA_DIR/timestamps/<current year>/<year>-<month>-<day>
    # eg $TRAQ_DATA_DIR/timestamps/2012/2012-12-12

or, if `-p <project>` was given, under

    $ $TRAQ_DATA_DIR/<project>/<current year>/<year>-<month>-<day>
    # eg $TRAQ_DATA_DIR/client-a/2012/2012-12-12

Each file can contain multiple lines of the following format:


    <timestamp>;<tag>;<comment>


Here's some sample content:

    Thu Sep 27 07:05:05 +0400 2012;#foo;Worked on Foo
    Thu Sep 27 07:15:05 +0400 2012;#bar;
    Thu Sep 27 07:25:05 +0400 2012;stop;

[1]:http://mxcl.github.com/homebrew/