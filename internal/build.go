package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/yuin/goldmark"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func (app *App) Build() error {
	if err := os.RemoveAll("./dist"); err != nil {
		return errors.New("unable to remove dist")
	}
	if err := BuildDirStruct("./dist"); err != nil {
		return err
	}
	if err := filepath.Walk("./static", func(path string, info os.FileInfo, err error) error{
		fileName := strings.Replace(path, "static", "dist", -1)
		if info.IsDir(){
			return os.MkdirAll(fileName, os.ModePerm)
		}
		return Copy(path, fileName)
	}); err != nil {
		panic(err)
		return errors.New("unable to copy static files")
	}


	t, err := template.ParseFiles("./templates/template.html")
	if err != nil {
		return errors.New("unable to parse template")
	}

	indexMaps := make(map[string]*Index)
	err = filepath.Walk("./content", func(path string, info os.FileInfo, err error) error{

		fileName := strings.Replace(path, "content", "dist", -1)
		fileName = strings.Replace(fileName, ".md", ".html", -1)

		parent := filepath.Dir(fileName)
		parent = filepath.Base(parent)

		base := filepath.Base(fileName)


		// Skip directories
		if info.IsDir() {
			indexMaps[base] = &Index{
				Title: strings.Title(base),
				Path: fileName,
			}
			return nil
		}

		index, ok := indexMaps[parent]
		if !ok {
			return nil
		}

		// If we are creating an index file here we don't want to overwrite it
		if base == "index.html" {
			index.Created = true
		}

		title, created := app.extractTitle(fileName)
		title = strings.Replace(title, "-", " ", -1)
		title = properTitle(title)

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.New("unable to read file")
		}

		var buf bytes.Buffer
		if err := goldmark.Convert(contents, &buf); err != nil {
			return err
		}

		data := &PageData{
			Title:   title,
			Content: template.HTML(buf.String()),
			Email: app.Email,
			Twitter: app.Twitter,
			LinkedIn: app.LinkedIn,
			GitHub: app.GitHub,
			Name: app.Name,
			Tracking: app.Tracking,
		}

		f, err := os.Create(fileName)
		if err != nil {
			return errors.New("unable to create file")
		}
		defer f.Close()


		if err := t.Execute(f, data); err != nil {
			return errors.New("unable to write file")
		}

		if index.Created != true {
			index.Entries = append(index.Entries, IndexEntry{
				Title: title,
				Path: strings.Replace(fileName, "dist", "", 1),
				Created: created,
			})
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, index := range indexMaps {
		if index.Created {
			continue
		}
		entries := index.Entries
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Created.After(entries[j].Created)
		})
		var builder strings.Builder
		builder.WriteString(`<ul class="index-list">`)
		for _, entry := range entries {
			builder.WriteString(`<li>`)
			if !entry.Created.IsZero() {
				formattedDate := entry.Created.Format("January 02, 2006")
				builder.WriteString(`<span class="date">`)
				builder.WriteString(formattedDate)
				builder.WriteString(`</span>`)
			}
			builder.WriteString(fmt.Sprintf(`<a class="title" href="%s">%s</a>`, entry.Path, entry.Title) + "\n")
		}

		data := &PageData{
			Title:   index.Title,
			Content: template.HTML(builder.String()),
			Email: app.Email,
			Twitter: app.Twitter,
			LinkedIn: app.LinkedIn,
			GitHub: app.GitHub,
			Name: app.Name,
			Tracking: app.Tracking,
		}

		f, err := os.Create(index.Path + "/index.html")
		if err != nil {
			return errors.New("unable to create file")
		}
		defer f.Close()

		if err := t.Execute(f, data); err != nil {
			return errors.New("unable to write file")
		}
	}

	return nil
}

func (app *App) extractTitle(fileName string) (title string, created time.Time) {
	title = strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	if title == "index" {
		dir := filepath.Dir(fileName)
		title = filepath.Base(dir)
		if title == "dist" {
			title = app.Name
		}
		if title == "" {
			title = "Welcome!"
		}
		return title, created
	}

	titleArray := strings.Split(title, "_")
	if len(titleArray) < 2 {
		return title, created
	}

	title = titleArray[1]
	created, err := time.Parse("2006-01-02", titleArray[0])
	if err != nil {
		fmt.Printf("unable to parse date: %s", titleArray[0])
	}

	return title, created
}

func properTitle(input string) string {
	words := strings.Fields(input)
	smallwords := " a an on the to "

	for index, word := range words {
		if index !=0 && strings.Contains(smallwords, " "+word+" "){ //Don't lowercase the first letter
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

type PageData struct{
	Title string
	Content template.HTML
	Email string
	Twitter string
	LinkedIn string
	GitHub string
	Name string
	Tracking string
}

type Index struct {
	Created bool
	Entries []IndexEntry
	Title string
	Path string
}

type IndexEntry struct {
	Title string
	Created time.Time
	Path string
}

