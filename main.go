package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	outDirName string
	numFile    int
)

func init() {
	outDirName = time.Now().Format("150405.00")
	if err := os.Mkdir(outDirName, os.ModeDir); err != nil {
		log.Fatalln(err)
	}

	flag.IntVar(&numFile, "c", 0, "Number of image.")
	flag.Parse()
}

func saveImage(img *image.Gray, fileName string, wg *sync.WaitGroup) error {
	defer wg.Done()

	annoFile, err := os.Create(outDirName + "/" + fileName)
	if err != nil {
		return err
	}

	defer annoFile.Close()

	if err := png.Encode(annoFile, img); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func convert(file *os.File, wg *sync.WaitGroup) error {
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	max := bounds.Max
	annoImg := image.NewGray(image.Rect(0, 0, max.X, max.Y))

	for x := 0; x < max.X; x++ {
		for y := 0; y < max.Y; y++ {
			r, g, _, a := img.At(x, y).RGBA()
			r >>= 8
			g >>= 8
			a >>= 8

			if a == 0 {
				annoImg.SetGray(x, y, color.Gray{0})
			}

			if r == 255 {
				annoImg.SetGray(x, y, color.Gray{2})
			}

			if g == 255 {
				annoImg.SetGray(x, y, color.Gray{1})
			}
		}
	}

	annoName := splitFileName(file.Name())
	if err := saveImage(annoImg, annoName, wg); err != nil {
		return err
	}

	return nil
}

func splitFileName(path string) string {
	idx := len(path)
	for ; idx > 0; idx-- {
		if path[idx-1] == '/' {
			break
		}
	}

	return path[idx:]
}

func main() {
	var wg sync.WaitGroup

	for idx := 0; idx < numFile; idx++ {
		file, err := os.OpenFile("E:/annotation/bottle_"+strconv.Itoa(idx)+".png", os.O_RDONLY, 0)
		if err != nil {
			log.Fatalln(err)
		}

		wg.Add(1)
		go func() {
			if err := convert(file, &wg); err != nil {
				log.Fatalln(err)
			}
		}()
	}

	wg.Wait()
}
