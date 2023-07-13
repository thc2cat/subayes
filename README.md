# Subayes

Bayesian filter for mail subjects Ham/Spam discrimination using [golang brukh/bayesian
 lib](https://github.com/brukh/bayesian).

## Context

Spammer uses a lot of differents subjects, sometime with wrong spelling and garbage.

Purpose of this project is a classifier able to learn/identify spam/Ham mail subjects.

## Basics

```shell
## Building
$ go mod tidy && go build 

## Learning
$ rm db/Spam db/Ham
$ ./subayes  -learnHam -d testdata/Ham -v
INFO classifier corpus :  [ Ham -> 0 items ]
INFO classifier corpus :  [ Ham -> 4623 items ]
$ ./subayes  -learnSpam -d testdata/esteban.txt -v
INFO classifier corpus :  [ Spam -> 0 items ]
INFO classifier corpus :  [ Spam -> 1096 items ]

## Evaluating words scores
$ echo "mensaje al grupo de trabajo please" | subayes -E    
[ mensaje = Spam ] : [Ham]{ 0.4000 } [Spam]{ 0.6000 } 
[ grupo = Ham ] : [Ham]{ 0.5096 } [Spam]{ 0.4904 } 
[ trabajo = Ham ] : [Ham]{ 0.6667 } [Spam]{ 0.3333 } 
[ please = Ham ] : [Ham]{ 0.6667 } [Spam]{ 0.3333 } 
Ham: mensaje al grupo de trabajo please

## Spam detection from stdin
$ ./subayes < testdata/2023-05  | grep -c Spam
59213
$ wc -l 2023-05
241662 2023-05 ( meaning 24% Spam, WTF! )

## Relearning
$ ./subayes -learnHam -d testdata/Ham-rajout-1.txt -v
INFO classifier corpus :  [ Ham -> 4623 items ]
INFO classifier corpus :  [ Ham -> 4718 items ]
$ ./subayes < testdata/2023-05  | grep -c Spam
58240

```

## Usage

Use
[utf8submimedecode](https://github.com/thc2cat/utf8submimedecode)
filter to decode  utf8 encoded subjects lines.

ex-pat contains lines to ignore patterns ( like Spam, or already detected users ).

subjects.sed is a sed script extracting subjects from log line.

subayes will create two files in db/ : Spam and Ham

```shell

# Detection

logs/partage$ rg -z clamav  sftp_logs/$LOGDATE/*clamav.log* | rg -vf ex-pat |\
 sed -f subjects.sed  | utf8submimedecode | sort -u | subayes | rg Spam | \
 tee  subayes.spam | mail -E -s "[subayes detection]" postmaster

# Learning more Ham words :  

logs/partage$ rg -z clamav  sftp_logs/$DATES/*clamav.log*  | rg -vf ex-pat|\
 sed -f subjects.sed  | utf8submimedecode | sort -u | subayes | rg Spam |\
 cut -c7- | tee subayes.spam 

 # edit ex-pat ( when you find new spammer address )

 # edit subayes.spam  (when you have false positives and relearn :)

logs/partage$ subayes  -v -learnHam -d subayes.spam          

# If you want to know what are the words tagged with Spam in a line, 
# use "-E explain", save, edit and  relearn.

$ subayes -E < 2023.subjects 2>&1 | awk '/^\[/ { if ($4=="Spam") print $2 }' |\
  sort -u | tee  subayes.words  

# Efficiency :

logs/partage$  subayes < /tmp/Hacked-account-Subjects | cut -d: -f1 | sort | uniq -c
5658 Ham
39016 Spam ( meaning 87% detection without false positives from filtered subjects)
                  
```

## Next move

Using this db for a postfix milter that would defer these subjects ?
