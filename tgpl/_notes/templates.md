# Text and HTML templates
- `text/template`, `html/template` provide a mechanism for substituting the values of variables into a text or HTML template
- A template is a string or file containing one or more portions enclosed in double braces, `{{...}}`, called *actions*. Most of the string is printed literally, but the actions trigger other behaviors
- Each action contains an expression in the template language, a simple but powerful notation for printing values, selecting struct fields, calling functions and methods, expressing control flow such as `if-else` statements and `range` loops, and instantiating other templates

    ```go
    const templ = `{{.TotalCount}} issues:
    {{range .Items}}-------------------------
    Number: {{.Number}}
    User:   {{.User.Login}}
    Title:  {{.Title | printf "%.64s"}}
    Age:    {{.CreatedAt | daysAgo}} days
    {{end}}`

    func daysAgo(t time.Time) int {
        return int(time.Since(t).Hours() / 24)
    }
    ```

    - Within an action, there's a notion of the current value, referred to as `.`. The dot initially refers to the template's parameter
    - `{{range .Items}}` and `{{end}}` actions create a loop
    - Within an action, the `|` notation makes the result of one operation the argument of another (analogous to a Unix shell pipeline)
    - `printf` is a built-in synonym for `fmt.Sprintf` in all templates
    - The type of `CreatedAt` is `time.Time`, not string. In the same way that a type may control its string formatting by defining certain (`String`) methods, a type also may define methods to control its JSON marshaling and unmarshaling behavior. The JSON-marshaled value of a `time.Time` is a string in a standard format
- Producing output with a template is a two-step process
   1. Parse the template into a suitable internal representation

        ```go
        // chaining methods
        report, err := template.New("report").
            Funcs(template.FuncMap{"daysAgo": daysAgo}).
            Parse(templ)
        if err != nil {
            log.Fatal(err)
        }
        ```

        - `template.New` creates and returns a template; `Funcs` adds `daysAgo` to the set of functions accessible within this template, then returns that template; `Parse` is called on the result
        - Parsing need be done only once
   2. Execute it on specific inputs

        ```go
        var report = template.Must(template.New("issuelist").
            Funcs(template.FuncMap{"daysAgo": daysAgo}).
            Parse(templ))

        func main() {
            result, err := github.SearchIssues(os.Args[1:])
            if err != nil {
                log.Fatal(err)
            }
            if err := report.Execute(os.Stdout, result); err != nil {
                log.Fatal(err)
            }
        }
        ```

    - The `template.Must` helper function accepts a template and an error, checks that the error is nil (and panics otherwise), and then returns the template
- Because templates are usually fixed at compile time, failure to parse a template indicates a fatal bug in the program
- `html/template` uses the same API and expression language as `text/template` but adds features for automatic and context-appropriate escaping of strings appearing within HTML, JavaScript, CSS, or URLs
    - These features can help avoid a perennial problem of HTML generation, an injection attack, in which an adversary crafts a **string value** like the title of an issue to include malicious code, that improperly escaped by a template, gives them control over the page
    - > The `html/template` automatically ***HTML-escaped*** the titles so that `"<"`, `">"` appear literally (`"&lt;"`). Had we use the `text/template`, the four-character string `"&lt;"` would have been rendered as a less-than character `"<"`, and the string `"<link>"` would have become a `link` element
    - We can suppress this auto-escaping behavior for fields that contain trusted HTML data by using the named string type `template.HTML` instead of string. Similar named types exist for trusted JavaScript, CSS, and URLs