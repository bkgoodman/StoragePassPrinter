package main

import (
	"fmt"
  "time"
  "os"
	"image/png"
  "github.com/fogleman/gg"

)

func exportbmp(filename string, xstart int, ystart int, usbDeviceFile *os.File) {
    logo, err := os.OpenFile(filename, os.O_RDWR, 0644)
    img, err := png.Decode(logo)
    if err != nil {
        panic(err)
    }

    fmt.Println("XY Aare",img.Bounds().Max.X,img.Bounds().Max.Y)
    for y := 0; y < img.Bounds().Max.Y; y++ {
      usbDeviceFile.Write([]byte(fmt.Sprintf("BITMAP %d,%d,%d,1,0,",xstart,ystart+y,img.Bounds().Max.X/8)))
    for x := 0; x < img.Bounds().Max.X; x+= 8{
      data := 0
      for i:=0;i<8;i++ {
            r, _, _, _ := img.At(x+i, y).RGBA()
            if (r > 0x8000) {
              data |= (1 << (7-i))
              //fmt.Println("-- PT",x+i,y,r)
            }
          }
          //fmt.Println("TEST",x,y,data)
          //data ^=0xAA
          usbDeviceFile.Write([]byte{byte(data)})
          if (data == 0xff) {
            //fmt.Print("  ")
          } else {
            //fmt.Printf("%x",data)
          }
        }
          usbDeviceFile.Write([]byte("\n"))
          //fmt.Println("")
    }
}

func drawCenteredString(str string,y int,fontsize int,usbDeviceFile *os.File) {
    // Now print from weird library
    var WIDTH int = 800
    var HEIGHT int = fontsize*2

    dc := gg.NewContext(WIDTH, HEIGHT)
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    dc.SetRGB(0, 0, 0)
    if err := dc.LoadFontFace("Ubuntu-R.ttf", 80); err != nil {
      panic(err)
    }
    dc.DrawStringAnchored(str, float64(WIDTH/2), float64(HEIGHT/2), 0.5, 0.5)
    dc.SavePNG("out.png")

    exportbmp("out.png",6,y,usbDeviceFile)
}

func main() {


      // Open the USB device file
    logo, err := os.OpenFile("makeit_logo_lable.png", os.O_RDWR, 0644)
    img, err := png.Decode(logo)
    if err != nil {
        panic(err)
    }

    usbDeviceFile, err := os.OpenFile("/dev/usb/lp0", os.O_RDWR, 0644)
    //usbDeviceFile, err := os.OpenFile("/dev/tty", os.O_RDWR, 0644)
    if err != nil {
        fmt.Println("Error opening USB device file:", err)
        return
    }
    defer usbDeviceFile.Close()
    //arr := []byte("SIZE 6,4\nGAP 0.13,0\nDIRECTION 1\nCLS\nTEXT 10,10,\"0\",0,1,1,\"Hello, TSPL Printer!\"\nPRINT 1\n")
    arr := []byte("SIZE 6,4\nGAP 0.13,0\nCLS\n")

    //arr := []byte("SIZE 6,4\nGAP 0.13,0\nCLS\nCIRCLE 250,20,100,5\nPRINT 1\n")
    //arr := []byte("SIZE 6,4\nGAP 0.13,0\nCLS\nTEXT 1,1,\"3\",0,1,1,\"Hello\"\nPRINT 1\n")
    usbDeviceFile.Write(arr)

    //usbDeviceFile.Write([]byte("BITMAP 10,10,4,1,0,55 55 FF FF\n"))

    fmt.Println("XY Aare",img.Bounds().Max.X,img.Bounds().Max.Y)
    for y := 0; y < img.Bounds().Max.Y; y++ {
      usbDeviceFile.Write([]byte(fmt.Sprintf("BITMAP 60,%d,%d,1,0,",y,img.Bounds().Max.X/8)))
    for x := 0; x < img.Bounds().Max.X; x+= 8{
      data := 0
      for i:=0;i<8;i++ {
            r, _, _, _ := img.At(x+i, y).RGBA()
            if (r > 0x8000) {
              data |= (1 << (7-i))
              //fmt.Println("-- PT",x+i,y,r)
            }
          }
          //fmt.Println("TEST",x,y,data)
          //data ^=0xAA
          usbDeviceFile.Write([]byte{byte(data)})
        }
          usbDeviceFile.Write([]byte("\n"))
    }

    /*
    // Now print from weird library
    var WIDTH int = 800
    var HEIGHT int = 160

    dc := gg.NewContext(WIDTH, HEIGHT)
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    dc.SetRGB(0, 0, 0)
    if err := dc.LoadFontFace("Ubuntu-R.ttf", 80); err != nil {
      panic(err)
    }
    dc.DrawStringAnchored("Hello, World!", float64(WIDTH/2), float64(HEIGHT/2), 0.5, 0.5)
    dc.SavePNG("out.png")

    exportbmp("out.png",6,300,usbDeviceFile)
    */

  drawCenteredString("Hello, World!",460,80,usbDeviceFile)
      currentDate := time.Now()
    futureDate := currentDate.AddDate(0, 0, 3)
    futureDateString := futureDate.Format("Mon, 02-Jan-06")
    fmt.Println(futureDateString)
  drawCenteredString(futureDateString,560,80,usbDeviceFile)

    usbDeviceFile.Write([]byte("PRINT 1\n"))
    //time.Sleep(5 * time.Second)
    var inbuf []byte
    test,err := usbDeviceFile.Read(inbuf)
    fmt.Println(test,err,inbuf)
    fmt.Println("Done")
  }
