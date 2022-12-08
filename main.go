package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
)

// get list of directories in content/restricted folder and store in variable
var restrictedFolders, _ = ioutil.ReadDir("content/restricted")
var postsFolders, _ = ioutil.ReadDir("content/posts")

// struct of post
type Post struct {
	Id           string
	Title        string
	Content      string
	DateCreated  string
	DateModified string
	Hash         string
	Locale       string
}

type Entry struct {
	Id     string
	Title  string
	Date   string
	Locale string
	Hash   string
	Folder string
}

type Tag struct {
	Name   string
	Locale string
	Posts  []string
}

type Category struct {
	Name   string
	Locale string
	Posts  []string
}

type Frontmatter struct {
	Title      string   `yaml:"title"`
	Author     string   `yaml:"author"`
	Date       string   `yaml:"date"`
	LastMod    string   `yaml:"lastmod"`
	Tags       []string `yaml:"tags"`
	Catetories []string `yaml:"categories"`
}

// create a map of tags and categories
var tags = make(map[string]Tag)
var categories = make(map[string]Category)

// main function
func main() {
	// create public if not exist
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.Mkdir("public", 0755)
	}

	// call return_key
	key := return_key()

	// create file
	// file, _ := os.Create("public/list.json");
	id := 0

	// create a map of entries, key is locale
	entries := make(map[string][]Entry)

	// make a list of string
	var locales []string

	if _, err := os.Stat("public/restricted"); os.IsNotExist(err) {
		os.Mkdir("public/restricted", 0755)
	}

	for _, folder := range restrictedFolders {
		// create sha1 from folder name
		hash := sha1.New()
		hash.Write([]byte(folder.Name()))
		hashedFolderName := fmt.Sprintf("%x", hash.Sum(nil))

		// create folder restricted under public folder if not exist
		if _, err := os.Stat("public/restricted/" + hashedFolderName); os.IsNotExist(err) {
			os.Mkdir("public/restricted/"+hashedFolderName, 0755)
		}

		// list folder contents
		files, _ := ioutil.ReadDir("content/restricted/" + folder.Name())

		for _, file := range files {
			// if file is a markdown file
			if file.Name()[len(file.Name())-3:] == ".md" {
				// create a list entry
				var post Post
				var entry Entry
				post.Id = fmt.Sprintf("%d", id)
				entry.Id = fmt.Sprintf("%d", id)
				id += 1
				post.Hash = hashedFolderName
				entry.Hash = hashedFolderName

				// get locale of the file
				post.Locale = file.Name()[6 : len(file.Name())-3]
				entry.Locale = post.Locale

				// add locale to locales list if not present
				if !contains(locales, post.Locale) {
					locales = append(locales, post.Locale)
				}

				var matter Frontmatter

				// read file as string
				var fileContent, _ = ioutil.ReadFile("content/restricted/" + folder.Name() + "/" + file.Name())

				// []byte to string
				var fileString = string(fileContent)

				// parse frontmatter
				rest, err := frontmatter.Parse(strings.NewReader(fileString), &matter)

				if err != nil {
					fmt.Println(err)
				}

				for _, tag := range matter.Tags {
					// if tag exists
					if _, ok := tags[tag]; ok {
						// add entry to tag
						if tagentry, ok := tags[tag]; ok {
							tagentry.Posts = append(tagentry.Posts, post.Id)
							tags[tag] = tagentry
						}
					} else {
						// create tag
						var tagentry Tag
						tagentry.Name = tag
						tagentry.Locale = post.Locale
						tagentry.Posts = append(tagentry.Posts, post.Id)
						tags[tag] = tagentry
					}
				}

				for _, category := range matter.Catetories {
					// if category exists
					if _, ok := categories[category]; ok {
						// add entry to category
						if categoryentry, ok := categories[category]; ok {
							categoryentry.Posts = append(categoryentry.Posts, post.Id)
							categories[category] = categoryentry
						}
					} else {
						// create category
						var categoryentry Category
						categoryentry.Name = category
						categoryentry.Posts = append(categoryentry.Posts, post.Id)
						categories[category] = categoryentry
					}
				}

				post.Title = matter.Title
				entry.Title = post.Title

				// matter.date reformat to 2006-01-02
				post.DateCreated = matter.Date[0:10]
				post.DateModified = matter.LastMod[0:10]
				entry.Date = post.DateCreated

				post.Content = string(rest)

				// check if locale exists in map
				if entryarray, ok := entries[post.Locale]; ok {
					// append entry to list of entries
					entryarray = append(entryarray, entry)
					entries[post.Locale] = entryarray
				} else {
					// create list of entries
					var entryarray []Entry
					entryarray = append(entryarray, entry)
					// add entry to map
					entries[post.Locale] = entryarray
				}

				// post to json
				postJson, _ := json.Marshal(post)

				encrypted := CBCEncrypt(key, postJson)

				// write to file
				write_if_changes("public/restricted/"+hashedFolderName+"/"+post.Id+".json", []byte(encrypted))
				// fmt.Printf("%+v\n", post)
				_ = rest
			} else {
				// copy file to public folder
				fileContent, _ := ioutil.ReadFile("content/restricted/" + folder.Name() + "/" + file.Name())
				write_if_changes("public/restricted/"+hashedFolderName+"/"+file.Name(), fileContent)
			}
		}

		// loop through entry map
		for locale, entry := range entries {
			// create json file
			entryJson, _ := json.Marshal(entry)
			write_if_changes("public/restricted/list_"+locale+".json", entryJson)
		}

		// write tags to file
		tagsJson, _ := json.Marshal(tags)
		write_if_changes("public/restricted/tags.json", tagsJson)

		// write categories to file
		categoriesJson, _ := json.Marshal(categories)
		write_if_changes("public/restricted/categories.json", categoriesJson)
	}

	if _, err := os.Stat("public/posts"); os.IsNotExist(err) {
		os.Mkdir("public/posts", 0755)
	}

	// create a map of entries, key is locale
	entries = make(map[string][]Entry)

	// loop through postsfolders
	for _, folder := range postsFolders {
		if _, err := os.Stat("public/posts/" + folder.Name()); os.IsNotExist(err) {
			os.Mkdir("public/posts/"+folder.Name(), 0755)
		}

		// list folder contents
		files, _ := ioutil.ReadDir("content/posts/" + folder.Name())

		for _, file := range files {
			// if file is a markdown file
			if file.Name()[len(file.Name())-3:] == ".md" {
				// create a list entry
				var post Post
				var entry Entry
				post.Id = fmt.Sprintf("%d", id)
				entry.Id = fmt.Sprintf("%d", id)
				id += 1
				entry.Folder = folder.Name()

				// get locale of the file
				post.Locale = file.Name()[6 : len(file.Name())-3]
				entry.Locale = post.Locale

				if !contains(locales, post.Locale) {
					locales = append(locales, post.Locale)
				}

				var matter Frontmatter

				// read file as string
				var fileContent, _ = ioutil.ReadFile("content/posts/" + folder.Name() + "/" + file.Name())

				// []byte to string
				var fileString = string(fileContent)

				// parse frontmatter
				rest, err := frontmatter.Parse(strings.NewReader(fileString), &matter)

				if err != nil {
					fmt.Println(err)
				}

				for _, tag := range matter.Tags {
					// if tag exists
					if _, ok := tags[tag]; ok {
						// add entry to tag
						if tagentry, ok := tags[tag]; ok {
							tagentry.Posts = append(tagentry.Posts, post.Id)
							tags[tag] = tagentry
						}
					} else {
						// create tag
						var tagentry Tag
						tagentry.Name = tag
						tagentry.Locale = post.Locale
						tagentry.Posts = append(tagentry.Posts, post.Id)
						tags[tag] = tagentry
					}
				}

				for _, category := range matter.Catetories {
					// if category exists
					if _, ok := categories[category]; ok {
						// add entry to category
						if categoryentry, ok := categories[category]; ok {
							categoryentry.Posts = append(categoryentry.Posts, post.Id)
							categories[category] = categoryentry
						}
					} else {
						// create category
						var categoryentry Category
						categoryentry.Name = category
						categoryentry.Posts = append(categoryentry.Posts, post.Id)
						categories[category] = categoryentry
					}
				}

				post.Title = matter.Title
				entry.Title = post.Title

				// matter.date reformat to 2006-01-02
				post.DateCreated = matter.Date[0:10]
				post.DateModified = matter.LastMod[0:10]
				entry.Date = post.DateCreated

				post.Content = string(rest)

				// check if locale exists in map
				if entryarray, ok := entries[post.Locale]; ok {
					// append entry to list of entries
					entryarray = append(entryarray, entry)
					entries[post.Locale] = entryarray
				} else {
					// create list of entries
					var entryarray []Entry
					entryarray = append(entryarray, entry)
					// add entry to map
					entries[post.Locale] = entryarray
				}

				// post to json
				postJson, _ := json.Marshal(post)

				// write to file
				write_if_changes("public/posts/"+folder.Name()+"/"+post.Id+".json", []byte(postJson))
				// fmt.Printf("%+v\n", post)
				_ = rest
			} else {
				// copy file to public folder
				fileContent, _ := ioutil.ReadFile("content/posts/" + folder.Name() + "/" + file.Name())
				write_if_changes("public/posts/"+folder.Name()+"/"+file.Name(), fileContent)
			}
		}

		// loop through entry map
		for locale, entry := range entries {
			// create json file
			entryJson, _ := json.Marshal(entry)
			write_if_changes("public/posts/list_"+locale+".json", entryJson)
		}

		// write tags to file
		tagsJson, _ := json.Marshal(tags)
		write_if_changes("public/posts/tags.json", tagsJson)

		// write categories to file
		categoriesJson, _ := json.Marshal(categories)
		write_if_changes("public/posts/categories.json", categoriesJson)
	}

	// write locale to json
	localesJson, _ := json.Marshal(locales)
	write_if_changes("public/locales.json", localesJson)
}
