Info:
Program dekoduje signal zakodovany v subore *.cod. Program je implementovany v jazyku Golang.
Sklada sa z 2 objektov Matrix a Codedfile. Matrix objekt obsahuje nacitane data zo suboru s maticou.
Tieto data tiez dokaze nacitat po zadani nazvu suboru, z ktoreho ma citat.
Codedfile obsahuje kodovane aj dekodovane data nacitane z *.cod suboru a taktiez vysledny signal.
Pomocou metody read() nacita data z *.cod, potom pomocou decode() rozkoduje data a pomocou synthetize() vytvori vysledny signal.
 

Kompilacia programu:
# export GOPATH=`pwd`
# #go get github.com/unixpickle/wav #neni potreba, kniznice su prilozene
# go build src/zre_proj1_linux.go

Spustenie:
./zre_proj1_linux cb512.txt gcb128.txt testmale.cod out.wav

Externe knihovny:
Kniznica github.com/unixpickle/wav pre pracu s wav subormi moznost ziskat pomocou:
go get github.com/unixpickle/wav 


Kompilacia:
Pomocou prekladaca go nasledovne:
# go build src/zre_proj1_linux.go
