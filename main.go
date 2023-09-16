package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/disintegration/imaging"
	"github.com/thteam47/resize-image/models"
	"image"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Resize Image")
	w.Resize(fyne.NewSize(1100, 700))
	w.CenterOnScreen()
	//data := []models.Todo{}
	todos := binding.NewUntypedList()
	//for _, t := range data {
	//	todos.Append(t)
	//}
	newtodoDescTxtHeight := widget.NewEntry()
	newtodoDescTxtHeight.PlaceHolder = "Height"
	newtodoDescTxtWidth := widget.NewEntry()
	newtodoDescTxtWidth.PlaceHolder = "Width"
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)
	progressBar.Hide()
	w.SetContent(
		container.NewBorder(
			container.NewVBox(
				newtodoDescTxtHeight,
				newtodoDescTxtWidth,
				container.NewVBox(
					progressBar,
				),
				widget.NewButton("Open Files", func() {
					dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
						imagesError := []string{}
						checkResolutionImage(&imagesError, uri.URI().Path())
						fmt.Println(imagesError)
						for _, path := range imagesError {
							todos.Append(models.NewTodo(path))
						}
						if len(imagesError) == 0 {
							dialog.ShowError(errors.New("No image invalid"), w)
						}
					}, w)
				}),
				widget.NewButton("Open Folder", func() {
					dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
						if err != nil {
							// Xử lý lỗi ở đây nếu có
							w.SetTitle("Error: " + err.Error())
							return
						}
						if uri != nil {
							// Xử lý đường dẫn thư mục đã chọn ở đây
							imagesError := []string{}
							checkResolutionImage(&imagesError, uri.Path())
							fmt.Println(imagesError)
							for _, path := range imagesError {
								todos.Append(models.NewTodo(path))
							}
							if len(imagesError) == 0 {
								dialog.ShowError(errors.New("No image invalid"), w)
							}
						} else {
							w.SetTitle("No Folder Selected")
						}

					}, w)
				}),
			),
			container.NewGridWithColumns(3,
				widget.NewButton("Resize", func() {
					dataTodos, _ := todos.Get()
					progressBar.SetValue(0)
					progressBar.Show()
					count := 0
					countDont := 0
					newTodos := []models.Todo{}
					for _, urlData := range dataTodos {
						todo := urlData.(models.Todo)
						if todo.Done {
							newTodos = append(newTodos, todo)
							countDont++
						}
					}
					for _, todo := range newTodos {
						err := fomattedImage(todo, newtodoDescTxtHeight, newtodoDescTxtWidth)
						if err != nil {
							dialog.ShowError(err, w)
							break
						}
						count++
						progressBar.SetValue(float64(count) / float64(len(newTodos)))
					}
					dialog.ShowInformation("Message", "Resize successful", w)
				}), // Left
				// RIGHT ↓
				widget.NewButton("Resize All", func() {
					dataTodos, _ := todos.Get()
					count := 0
					progressBar.SetValue(0)
					progressBar.Show()
					for _, urlData := range dataTodos {
						todo := urlData.(models.Todo)
						err := fomattedImage(todo, newtodoDescTxtHeight, newtodoDescTxtWidth)
						if err != nil {
							dialog.ShowError(err, w)
							break
						}
						count++
						progressBar.SetValue(float64(count) / float64(len(dataTodos)))
					}
					dialog.ShowInformation("Message", "Resize successful", w)
				}),
				widget.NewButton("Reset", func() {
					todos.Set(nil)
					progressBar.Hide()
					dialog.ShowInformation("Message", "Reset successful", w)
				}),
			),
			nil, // Right
			nil,
			widget.NewListWithData(
				todos,
				func() fyne.CanvasObject {
					return container.NewBorder(
						nil, nil, nil,
						// left of the border
						widget.NewCheck("", func(b bool) {}),
						// takes the rest of the space
						widget.NewLabel(""),
					)
				},
				func(di binding.DataItem, o fyne.CanvasObject) {
					ctr, _ := o.(*fyne.Container)
					l := ctr.Objects[0].(*widget.Label)
					c := ctr.Objects[1].(*widget.Check)
					todo := models.NewTodoFromDataItem(di)
					l.SetText(todo.Url)
					c.SetChecked(todo.Done)

					c.OnChanged = func(checked bool) {
						// Cập nhật trạng thái Done của TODO dựa trên giá trị checked
						todo.Done = checked
						dataTodos, _ := todos.Get()
						todos.Set(nil)
						for _, urlData := range dataTodos {
							if urlData.(models.Todo).Url == todo.Url {
								todos.Append(todo)
							} else {
								todos.Append(urlData)
							}
						}

						// Gán danh sách mới cho biến todos
						//todos.Set(updatedTodos)
					}
				}),
		),
	)
	w.ShowAndRun()
}

func fomattedImage(todo models.Todo, newtodoDescTxtHeight *widget.Entry, newtodoDescTxtWidth *widget.Entry) error {
	srcImage, err := imaging.Open(todo.Url)
	if err != nil {
		panic(err)
	}
	newWidth := 1200
	newHeight := 1200
	if newtodoDescTxtHeight.Text != "" {
		newHeight, err = strconv.Atoi(newtodoDescTxtHeight.Text)
		if err != nil {
			return errors.New("Size height invalid")
		}
	}
	if newtodoDescTxtWidth.Text != "" {
		newWidth, err = strconv.Atoi(newtodoDescTxtWidth.Text)
		if err != nil {
			return errors.New("Size width invalid")
		}
	}
	dirName := filepath.Dir(todo.Url)
	dirNameFormatted := filepath.Join(dirName, "formatted")
	_, err = os.Stat(dirNameFormatted)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirNameFormatted, os.ModePerm)
		if err != nil {
			return errors.New(fmt.Sprintf("Lỗi khi tạo thư mục: %s", dirNameFormatted))
		}
	} else if err != nil {
		return errors.New(fmt.Sprintf("Lỗi khi kiểm tra thư mục: %s", dirNameFormatted))
	}
	formattedPath := filepath.Join(dirNameFormatted, filepath.Base(todo.Url))

	resizedImage := imaging.Resize(srcImage, newWidth, newHeight, imaging.Lanczos)
	contrastedImage := imaging.AdjustContrast(resizedImage, 10)
	sharpenedImage := imaging.Sharpen(contrastedImage, 1.0)

	err = imaging.Save(sharpenedImage, formattedPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Image %s format error", todo.Url))
	}
	return nil
}

func checkResolutionImage(imagesError *[]string, imageDir string) {
	fileInfo, err := os.Stat(imageDir)
	if err != nil {
		fmt.Printf("Lỗi khi kiểm tra đường dẫn: %v\n", err)
		return
	}
	if !fileInfo.IsDir() {
		fileName := filepath.Base(imageDir)

		if isImage(fileName) {
			// Đọc tệp tin ảnh
			imgFile, err := os.Open(imageDir)
			if err != nil {
				log.Println("Failed to open image:", imageDir)
				return
			}
			defer imgFile.Close()

			// Giải mã thông tin ảnh
			imgConfig, _, err := image.DecodeConfig(imgFile)
			if err != nil {
				log.Println("Failed to decode image:", imageDir)
			}

			if imgConfig.Width != imgConfig.Height {
				return
			} else if imgConfig.Width < 600 || imgConfig.Height < 600 {
				return
			}

			*imagesError = append(*imagesError, imageDir)
		}
		return
	}
	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if file.Name() == "formatted" {
				continue
			}
			filePath := filepath.Join(imageDir, file.Name())
			checkResolutionImage(imagesError, filePath)
		}
		// Kiểm tra xem tệp tin có phải là ảnh hay không
		if isImageFile(file) {
			// Đường dẫn đầy đủ đến tệp tin
			filePath := filepath.Join(imageDir, file.Name())

			// Đọc tệp tin ảnh
			imgFile, err := os.Open(filePath)
			if err != nil {
				log.Println("Failed to open image:", filePath)
				continue
			}
			defer imgFile.Close()

			// Giải mã thông tin ảnh
			imgConfig, _, err := image.DecodeConfig(imgFile)
			if err != nil {
				log.Println("Failed to decode image:", filePath)
				continue
			}

			if imgConfig.Width != imgConfig.Height {
				continue
			} else if imgConfig.Width < 600 || imgConfig.Height < 600 {
				continue
			}

			*imagesError = append(*imagesError, filePath)
		}
	}
}

// Kiểm tra xem tệp tin có phải là ảnh hay không
func isImageFile(fileInfo os.DirEntry) bool {
	extension := filepath.Ext(fileInfo.Name())
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}

func isImage(filename string) bool {
	// Chuyển đổi phần mở rộng của tệp thành chữ thường để so sánh dễ dàng hơn.
	ext := strings.ToLower(filepath.Ext(filename))

	// Danh sách các phần mở rộng thường được sử dụng cho hình ảnh.
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}

	return false
}
