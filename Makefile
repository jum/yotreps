TARG=yotreps
GOFILES=yotreps.go mbox.go wpt.go

$(TARG): $(GOFILES)
	go build -o $(TARG) $(GOFILES)

clean:
	rm -f $(TARG) mon.out
