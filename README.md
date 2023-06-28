# Subayes

Bayesian filter for mail subjects discrimination.

## Context

Spammer uses a lot of differents subjects, sometime with wrong spelling.
Purpose of this filter is a learning filter able to sort spam subject lines.

## Learning

```shell

$ rm db/Spam db/Ham

$ ./subayes.exe -learnHam -d testdata/Ham -v
INFO classifier corpus :  [ Ham -> 0 items ] [ Spam -> 0 items ] 
INFO classifier corpus :  [ Ham -> 4623 items ] [ Spam -> 0 items ]

$ ./subayes.exe -learnSpam -d testdata/esteban.txt -v
INFO classifier corpus :  [ Ham -> 4623 items ] [ Spam -> 0 items ]
INFO classifier corpus :  [ Ham -> 5719 items ] [ Spam -> 0 items ]

```

## Extraction des Spam

```shell

$ ./subayes.exe < testdata/2023-05  | grep Spam | wc -l
58711

$ wc -l 2023-05
241662 2023-05 ( meaning 24% Spam, WTF! )

```

## Réapprentissage

```shell

$ ./subayes.exe -learnHam -d testdata/Ham-rajout-1.txt -v
INFO classifier corpus :  [ Ham -> 4623 items ] [ Spam -> 0 items ]
INFO classifier corpus :  [ Ham -> 4718 items ] [ Spam -> 0 items ]

$ ./subayes.exe < 2023-05  | grep Spam | wc -l
58149

```
