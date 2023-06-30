# subayes

Bayesian filter for mail subjects Ham/Spam discrimination.

## Context

Spammer uses a lot of differents subjects, sometime with wrong spelling.

Purpose of this project is a filter able to identify spam/Ham mail subjects.

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

## Spam detection
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

ex-pat contains lines to ignore patterns ( like Spam, or already detected accounts ).

subjects.sed is a sed script extracting subjects from log line.

```shell
logs/partage$ rg -z clamav  sftp_logs/$LOGDATE/*clamav.log* | rg -vf ex-pat |\
 sed -f subjects.sed  | utf8submimedecode | sort -u | subayes | rg Spam | \
 tee  subayes.spam | mail -E -s "[subayes detection]" postmaster
```
