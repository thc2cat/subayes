# Subayes

This is a naive bayesian classifier for mail subjects. Ham/Spam discrimination using
[golang jbrukh/bayesian lib](https://github.com/jbrukh/bayesian).

![go.yml](https://github.com/thc2cat/subayes/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/thc2cat/subayes.svg)](https://pkg.go.dev/github.com/thc2cat/subayes)

## Context

Spammer uses a lot of differents subjects, sometime with wrong spelling and garbage.

Purpose of this project is a basic classifier able to detect spam from mail subjects.

subayes read stdin line and output them on stdout with prefix "Spam: " or "Ham: ".

## Basics

```shell
## Building
$ go mod tidy && go build 

## Defaults options : 
$ subayes -h
Usage of subayes:
  -E    explain words scores
  -d string
        data filename (default "subayes.spam")
  -db string
         db path (default "db")
  -learnHam
        learn Ham subjects
  -learnSpam
        learn Spam subjects
  -m int
        word min length (default 4)
  -v    verbose


## Learning
$ rm db/Spam db/Ham ; mkdir db
$ ./subayes  -learnHam -d testdata/Ham -v
INFO classifier corpus :  [ Ham -> 0 items ]
INFO classifier corpus :  [ Ham -> 4623 items ]
$ ./subayes  -learnSpam -d testdata/esteban.txt -v
INFO classifier corpus :  [ Spam -> 0 items ]
INFO classifier corpus :  [ Spam -> 1096 items ]

## Testing 
$ echo "mensaje al grupo de trabajo please" | subayes
Ham: mensaje al grupo de trabajo please

$ echo "View sexy women in your neighborhood" | subayes
Spam: View sexy women in your neighborhood


## Evaluating words scores
$ echo "mensaje al grupo de trabajo please" | subayes -E    
[ mensaje = Spam ] : [Ham]{ 0.4000 } [Spam]{ 0.6000 } 
[ grupo = Ham ] : [Ham]{ 0.5096 } [Spam]{ 0.4904 } 
[ trabajo = Ham ] : [Ham]{ 0.6667 } [Spam]{ 0.3333 } 
[ please = Ham ] : [Ham]{ 0.6667 } [Spam]{ 0.3333 } 
Ham: mensaje al grupo de trabajo please

## Raw test from v0.1
$ ./subayes.exe < testdata/2023-05 |cut -d: -f1|sort|uniq -c
 176347 Ham
  57102 Spam

Meaning at least 24% Spam ! 

```

## Common usage

Use
[utf8submimedecode](https://github.com/thc2cat/utf8submimedecode)
filter to decode utf8 encoded subjects lines.

ex-pat contains lines to ignore patterns ( like Spam, [PUB] or already detected users ).

subjects.sed is a simple sed script extracting subjects from log line.

subayes will create two files in db/ : Spam and Ham

Each time you find a spammer, learn theirs subjects as spam, verify updated db against previous clean data to adjust false positives.

```shell

# Detection from clamav logs

logs/partage$ rg -z clamav  sftp_logs/$LOGDATE/*clamav.log* \
| rg -vf ex-pat | sed -f subjects.sed  | utf8submimedecode \
| sort -u | subayes | rg ^Spam \
| tee  subayes.spam | mail -E -s "[subayes detection]" postmaster

# If you want to know what are the words tagged with Spam in a line, 
# use "-E" explain option (printed on stderr).

$ subayes -E < subayes.spam  

# Learning more Ham words :  
 # edit subayes.spam  (when you have false positives and relearn :)

logs/partage$ subayes  -v -learnHam -d subayes.spam          
( -d is optional, subayes.spam is the default data file)

# Efficiency :

logs/partage$  subayes < /tmp/Hacked-account-Subjects \
| cut -d: -f1 | sort | uniq -c
5658 Ham
39016 Spam ( meaning 87% detection without false positives from filtered subjects)
                  
```

## Next move

Using this db for a postfix milter that would defer these subjects ?
