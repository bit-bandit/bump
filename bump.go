package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func get_urls(urls [][]string) []string {
	u := []string{}
	for i := range urls {
		// Put registries we know we can't properly account for here
		// (IE, anything that *isn't* `deno.land/x` or `esm.sh`)
		if strings.HasPrefix(urls[i][1], "https://cdn.jsdelivr.net") ||
			strings.HasPrefix(urls[i][1], "https://raw.githubusercontent.com") {
			urls[i] = urls[len(urls)-1]
		}

		r := regexp.MustCompile(`(.*)@.*?`)
		res := r.FindAllStringSubmatch(urls[i][1], -1)

		if len(res) != 0 {
			u = append(u, res[0][1])
		} else {
			u = append(u, urls[i][1])
		}
	}
	return u
}

func new_urls(urls []string) []string {
	u := []string{}

	for i := range urls {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, err := http.NewRequest("GET", urls[i], nil)

		if err != nil {
			panic(err)
		}

		res, err := client.Do(req)

		if err != nil {
			panic(err)
		}

		redir, err := res.Location()

		// Return original url if no
		// location header is provided.
		if err != nil {
			u = append(u, urls[i])
		} else {
			u = append(u, redir.String())
		}
	}
	return u
}

func write_urls(file string, old [][]string, new []string) string {
	f := file
	for i := range old {
		// Only target modules with explicit versions.
		r := regexp.MustCompile(`https?://.*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?)+`)
		n := r.FindAllStringSubmatch(new[i], -1)
		o := r.FindAllStringSubmatch(old[i][1], -1)
		if len(n) > 0 && len(o) > 0 {
			f = strings.Replace(f, o[0][0], n[0][0], -1)
		}
	}
	return f
}

// TODO: Split these off into seperate functions.
func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 1 {
		panic("Too many directories.")
	} else if len(argsWithoutProg) == 0 {
		fmt.Println(
			"Usage: bump [directory]",
			"\nFile any issues at: https://github.com/bit-bandit/bump",
		)
		os.Exit(0)
	}

	dir, err := os.ReadDir(argsWithoutProg[0])

	if err != nil {
		panic(err)
	}

	fmt.Println("Analyzing", argsWithoutProg[0], "...")

	// Main loop that does everything that it should be doing already.
	for _, entry := range dir {
		if entry.IsDir() {
			fmt.Println("Skipping Directory:", entry)
		} else {
			var filename string

			ext := filepath.Ext(entry.Name())

			switch ext {
			case ".ts", ".js", ".json", ".jsonc":
				fmt.Println("Found source file:", entry.Name())
			default:
				fmt.Println("Ingnoring:", entry.Name())
				continue
			}

			// Terrible.
			if strings.HasSuffix(argsWithoutProg[0], "/") {
				filename = argsWithoutProg[0] + entry.Name()
			} else {
				filename = argsWithoutProg[0] + "/" + entry.Name()
			}

			file, err := os.ReadFile(filename)

			if err != nil {
				panic(err)
			}

			fileinfo, err := os.Stat(filename)

			if err != nil {
				panic(err)
			}

			ff := string(file)
			r := regexp.MustCompile(`[from|import|":] ?"(https?:\/\/.*)"[;|,]?`)

			old_list := r.FindAllStringSubmatch(ff, -1)

			plain_list := get_urls(old_list)
			new_list := new_urls(plain_list)
			final := write_urls(ff, old_list, new_list)

			d := []byte(final)
			m := fileinfo.Mode()

			err = os.WriteFile(filename, d, m)

			if err != nil {
				panic(err)
			}

			fmt.Println("Updated modules in", filename)
		}
	}
}
