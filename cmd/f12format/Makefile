test:
	cd ../gosk && go build
	cd ../f12copy && go build
	go build
	cd ../../ && ./cmd/gosk/gosk testdata/03_day_harib00f_ipl.nas ipl10.bin
	cd ../../ && ./cmd/gosk/gosk testdata/03_day_harib00f_haribote.nas haribote.sys
	./f12format -f 1440 -B ../../ipl10.bin -i ../../output.img
	../f12copy/f12copy -f ../../haribote.sys -i ../../output.img
