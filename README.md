# traq

bash based time tracking using text files

Requires bash v4.2.39 or newer. Works on Linux and OS X

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

# echo the content of all files from the calendar week 31 to stdout. If the week does not contain files, nothing is echoed.
$ traq -w 31

# starts time tracking for client-a on #development
$ traq -p client-a development

# stops time tracking for client-a
$ traq -p client-a stop

# list tracked times for client-a from today
$ traq -p client-a
```

## Installation

`traq` assumes you're installing it to your home directory, into `~/.traq`. This will set you up:

``` bash
$ mkdir $HOME/.traq
$ mkdir -p $HOME/Library/traq
$ git clone git@github.com:nicolai86/traq.git ~/.traq
$ echo "export TRAQ_PATH=$HOME/.traq" >> ~/.bash_profile
$ echo "export TRAQ_DATA_DIR=$HOME/Library/traq" >> ~/.bash_profile
$ echo "export PATH=$PATH:$HOME/.traq" >> ~/.bash_profile
$ . ~/.bash_profile
$ which traq
```

To update your installation all you need to do is to

``` bash
$ cd $HOME/.traq/traq
$ git pull origin master
```

**Linux Note** `traqeval` requires `bc` to be available. If `which bc` returns nothing you need install `bc` via `aptitude` or whatever package manager you're using.

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

## Tests

The project has some tests using [bats](https://github.com/sstephenson/bats). Assuming you got `bats` installed, run them using the following command:

``` bash
$ cd $TRAQ_PATH
$ bats tests/
```

## Hacking

All files are placed under

    $ $TRAQ_DATA_DIR/timestamps/<current year>/kw-<week number>/timestamps-<date>
    # eg $TRAQ_DATA_DIR/timestamps/2012/kw-50/timestamps-2012-12-12

or, if `-p <project>` was given, under

    $ $TRAQ_DATA_DIR/<project>/<current year>/kw-<week number>/timestamps-<date>
    # eg $TRAQ_DATA_DIR/client-a/2012/kw-50/timestamps-2012-12-12

Each file can contain multiple lines of the following format:


    <timestamp>;<tag>


Here's some sample content:

    Thu Sep 27 07:05:05 +0400 2012;#foo
    Thu Sep 27 07:15:05 +0400 2012;#bar
    Thu Sep 27 07:25:05 +0400 2012;stop

## Helpers

To ease evaluation of traq-files `traq` comes with two helper scripts, `traqtrans` and `traqeval`.

`traqtrans` transforms the timestamp into a unix timestamp,
and `traqeval` sums up tags.

    $ traq -p test -w 39 | traqtrans
    1348715105;#foo
    1348715705;#bar
    1348716305;stop
    %%

Pipe both together and you'll get something like this:

    $ traq -p test -w 39 | traqtrans | traqeval
    2012-09-27
    #foo:0.1666
    #bar:0.1666
    %%

[1]:http://mxcl.github.com/homebrew/