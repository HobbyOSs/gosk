package main

import (
	"flag"
	"github.com/hangingman/go-fs"
	"github.com/hangingman/go-fs/fat"
	"io/ioutil"
	"os"
)

// https://www.gnu.org/software/mtools/manual/mtools.html#mformat
// --------------------------------------------------------------
// mformat -f 1440                    [-f size]
//         -l HARIBOTEOS              [-v volume_label]
//         -N 0xffffffff              [-N serial_number]
//         -C                         [   creates the disk image file to install the MS-DOS file system on it.]
//         -B ${03_day_harib00g_IPLB} [-B boot_sector]
//         -i ${03_day_harib00g_OS}   [-i output image name]
//
func main() {

	var (
		fatSize     int64
		volumeLabel string
		bootSector  string
		outputImage string
	)
	flag.Int64Var(&fatSize, "f", 1440, "FATのサイズを指定する")
	flag.StringVar(&volumeLabel, "l", "HARIBOTEOS", "FATのボリュームラベルを指定する")
	flag.StringVar(&bootSector, "B", "", "FATのブートセクタを指定する")
	flag.StringVar(&outputImage, "i", "", "出力先の指定をする")

	flag.Parse()

	// ブートセクタを読み取り
	bootSectorContent, err := ioutil.ReadFile(bootSector)
	if bootSector != "" && err != nil {
		panic(err)
	}

	// 読み書き可能, 新規作成, ファイル内容あっても切り詰め
	floppyF, err := os.OpenFile(outputImage, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	if floppyF.Truncate(fatSize * 1024); err != nil {
		panic(err)
	}
	defer floppyF.Close()

	// ブロックデバイスを構築する
	device, err := fs.NewFileDisk(floppyF)
	if err != nil {
		panic(err)
	}

	// ブロックデバイスをフォーマットする
	formatConfig := &fat.SuperFloppyConfig{
		FATType:           fat.FAT12,
		Label:             volumeLabel,
		OEMName:           "go-fs",
		BootSectorContent: bootSectorContent,
	}
	if fat.FormatSuperFloppy(device, formatConfig); err != nil {
		panic(err)
	}
}
