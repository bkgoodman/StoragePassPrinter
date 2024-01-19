package main

import (
	"fmt"
  "time"
  "os"
  "strings"
	"image/png"
  "github.com/fogleman/gg"

)

func exportbmp_dymo(filename string, usbDeviceFile *os.File) {
    logo, err := os.OpenFile(filename, os.O_RDWR, 0644)
    img, err := png.Decode(logo)
    if err != nil {
        panic(err)
    }
    defer logo.Close()

    fmt.Println("XY Are",img.Bounds().Max.X,img.Bounds().Max.Y)
    for x := 0; x < img.Bounds().Max.X; x++ {
      usbDeviceFile.Write([]byte{byte(0x16)})
    for y := img.Bounds().Max.Y-1;y>=0; y-=8 {
      data := 0
      for i:=0;i<8;i++ {
            r, _, _, _ := img.At(x, y+i).RGBA()
            if (r <= 0x8000) {
              //data |= (1 << (7-i))
              data |= (1 << i)
              //fmt.Println("-- PT",x+i,y,r)
            }
          }
          usbDeviceFile.Write([]byte{byte(data)})
        }
    }
}

func exportbmp(filename string, xstart int, ystart int, usbDeviceFile *os.File) {
    logo, err := os.OpenFile(filename, os.O_RDWR, 0644)
    img, err := png.Decode(logo)
    if err != nil {
        panic(err)
    }
    defer logo.Close()

    fmt.Println("XY Aare",img.Bounds().Max.X,img.Bounds().Max.Y)
    for y := 0; y < img.Bounds().Max.Y; y++ {
      usbDeviceFile.Write([]byte(fmt.Sprintf("BITMAP %d,%d,%d,1,1,",xstart,ystart+y,img.Bounds().Max.X/8)))
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
    if err := dc.LoadFontFace("Ubuntu-R.ttf", float64(fontsize)); err != nil {
      panic(err)
    }
    dc.DrawStringAnchored(str, float64(WIDTH/2), float64(HEIGHT/2), 0.5, 0.5)
    dc.SavePNG("/run/lableout.png")

    exportbmp("/run/lableout.png",6,y,usbDeviceFile)
}

// This tag number tried to badge in
func PrintBadge(id int64) {
	var found bool = false
	for _,tag := range validTags {
		if id == int64(tag.Tag) {
			found = true
		fmt.Printf("Tag %d Member %s",id,tag.Member)

		var topic string = fmt.Sprintf("ratt/status/node/%s/personality/storagepass",cfg.ClientID)
		var message string = fmt.Sprintf("{\"allowed\":true,\"member\":\"%s\"}",tag.Member)
    print_storagelabel(strings.Replace(tag.Member,"."," ",-1))
		client.Publish(topic,0,false,message)
			return
		} 

	}

	if (found == false) {
		fmt.Println("Tag disallowed",id)
		return
	}
}
func print_storagelabel(member string) {


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

    if (true) {
      /* DYMO PRINTER */
      lines :=960
      bpl := 38 // Bytes Per Line

      // We are doing this rotated
      var HEIGHT = (bpl * 8)
      var WIDTH = lines
      dc := gg.NewContext(WIDTH,HEIGHT)
      dc.SetRGB(1, 1, 1)
      dc.Clear()
      dc.SetRGB(0, 0, 0)
      if err := dc.LoadFontFace("Ubuntu-R.ttf", float64(84)); err != nil {
        panic(err)
      }

      offset := float64(0)
      im, err := gg.LoadPNG("milsm.png")
      if err == nil {
        dc.DrawImage(im, 20, 20)
        offset =float64(im.Bounds().Dx()) /float64(2 )
      }

      currentDate := time.Now()
      futureDate := currentDate.AddDate(0, 0, 3)
      futureDateString := futureDate.Format("Mon, 02-Jan-06")

      formattedDateTime := currentDate.Format("Mon, 02-Jan-2006 01:04 PM")

      dc.DrawStringAnchored(futureDateString, float64(WIDTH/2)+offset, 180, 0.5, 0.5)
      dc.LoadFontFace("Ubuntu-R.ttf", float64(48))
      dc.DrawStringAnchored("Temporary Storage Pass", float64(WIDTH/2)+offset, 50, 0.5, 0.5)
      dc.DrawStringAnchored(member, float64(WIDTH/2)+offset, 120, 0.5, 0.5)
      dc.LoadFontFace("Ubuntu-R.ttf", float64(24))
      dc.DrawStringAnchored(fmt.Sprintf("Left on: %s",formattedDateTime), float64(WIDTH/2)+offset, 240, 0.5, 0.5)
      dc.SetLineWidth(2)
      dc.DrawRectangle(10, 10, float64(WIDTH-10), float64(HEIGHT-10))
      dc.Stroke()


      if err := dc.LoadFontFace("Ubuntu-R.ttf", float64(18)); err != nil {
        panic(err)
      }
      textbody := "Items may be discarded and disposal charges may be incurred if items are left after specified date."
      dc.DrawStringWrapped(textbody,float64(WIDTH/2)+offset,float64(HEIGHT-32) , 0.5, 0.5, float64(WIDTH/2), 1.2, gg.AlignCenter)
      dc.SavePNG("lableout.png")

      fmt.Println("Lines",lines,"colbytes",bpl)
      usbDeviceFile.Write([]byte{27,0x44,byte(bpl)}) // Width (Bytes)
      usbDeviceFile.Write([]byte{27,0x4c,byte((lines >> 8)&0xff),byte(lines &0xff)}) // 16 lines on lable
      /*
      l:=0
      i:=0
      for l=0;l<lines;l++ {
      usbDeviceFile.Write([]byte{0x16}) // 16 lines on lable
        for i=0;i<bpl;i++ {
          usbDeviceFile.Write([]byte{0xff}) // 16 lines on lable
        }
      }
      */
    exportbmp_dymo("lableout.png", usbDeviceFile) 
      usbDeviceFile.Write([]byte{27,'E'}) // Form Feed

    } else {
      /* NON-DYMO BIGGER PRINTER */
    //arr := []byte("SIZE 6,4\nGAP 0.13,0\nDIRECTION 1\nCLS\nTEXT 10,10,\"0\",0,1,1,\"Hello, TSPL Printer!\"\nPRINT 1\n")
    arr := []byte("\n\nSIZE 6,4\nGAP 0.13,0\nCLS\n")

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

  drawCenteredString(member,220,60,usbDeviceFile)
    currentDate := time.Now()
    futureDate := currentDate.AddDate(0, 0, 3)
    futureDateString := futureDate.Format("Mon, 02-Jan-06")
    fmt.Println(futureDateString)

    formattedDateTime := currentDate.Format("Mon, 02-Jan-2006 01:04 PM")

  drawCenteredString("Item was left on",337,36,usbDeviceFile)
  drawCenteredString(formattedDateTime,400,42,usbDeviceFile)
  drawCenteredString("Must be removed on or before",480,36,usbDeviceFile)
  drawCenteredString(futureDateString,520,100,usbDeviceFile)

    usbDeviceFile.Write([]byte("PRINT 1,1\n\n\n"))
    //time.Sleep(5 * time.Second)
    var inbuf []byte
    test,err := usbDeviceFile.Read(inbuf)
    fmt.Println(test,err,inbuf)
  }
}
