package main

import "github.com/disintegration/imaging"

//func main() {
//	// Open a test image.
//	src, err := imgconv.Open("./ong_sao.jpg")
//	if err != nil {
//		log.Fatalf("failed to open image: %v", err)
//	}
//
//	// Resize the image to width = 200px preserving the aspect ratio.
//	mark := imgconv.Resize(src, &imgconv.ResizeOption{Width: 1200, Height: 1200})
//
//	//dst := imgconv.Watermark(src, &imgconv.WatermarkOption{Mark: mark, Opacity: 128, Random: true})
//
//	err = imgconv.Write(io.Discard, mark, &imgconv.FormatOption{Format: imgconv.PNG})
//	if err != nil {
//		log.Fatalf("failed to write image: %v", err)
//	}
//}

func main111() {
	// Đọc hình ảnh từ tệp nguồn
	srcImage, err := imaging.Open("./ong_sao.jpg")
	if err != nil {
		panic(err)
	}

	// Chuyển đổi độ phân giải của hình ảnh (ví dụ: 800x600)
	newWidth := 1200
	newHeight := 1200
	resizedImage := imaging.Resize(srcImage, newWidth, newHeight, imaging.Lanczos)
	contrastedImage := imaging.AdjustContrast(resizedImage, 10)
	sharpenedImage := imaging.Sharpen(contrastedImage, 1.0)
	err = imaging.Save(sharpenedImage, "output.jpg")
	if err != nil {
		panic(err)
	}
}
