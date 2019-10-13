# Breno Helfstein Moura						NUSP: 9790972
# Matheus Barcellos de Castro Cunha   		NUSP: 11208238

.PHONY: clean

help:
	echo "Options are:\n\nmake build\nmake runMaster\nmake runDebugMaster\nmake runSlave\nmake runDebugSlave\nmake setMasterIP IP={IP TO SET}\nmake clean"

build:
	go build main.go

runMaster:
	./main list.txt

runDebugMaster:
	./main list.txt -d

runSlave:
	./main 

runDebugSlave:
	./main -d

setMasterIP:
	echo $(IP) | cat > address.conf

clean:
	rm -rf tmp/
	rm *~ eventLog.txt sorted.txt
