all: fcbh.pdf

clean:
	rm -f *.pdf fcbh fcbh.dot

%.pdf : %.dot
	dot -Tpdf $+ >$@

fcbh : fcbh.go
	go build $+

fcbh.dot : fcbh fcbh.csv 
	./fcbh > $@

fcbh.pdf : fcbh.dot
	dot $+ -Tpdf > $@
