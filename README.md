# md-hugo
Converts and moves Markdown files to Hugo content folders and adds meta data.

### If you have `go` installed

Just do a git pull on this repo
Then `cd` into `md-hugo`

``` bash
git pull https://github.com/forestgiant/md-hugo.git

cd md-hugo

```
Now run `go install` and md-hugo should be installed on your computer. To run it:

``` bash
md-hugo /path/to/process/folder /path/to/hugo/contents/folder
```

If you want to automatically `push` the hugo website live you can pass in the `-push` flag

``` bash
md-hugo -push /path/to/process/folder /path/to/hugo/contents/folder
```

### If you DO NOT have `go` installed

You can just run the binary so clone the repo and cd into it:

``` bash
git pull https://github.com/forestgiant/md-hugo.git

cd md-hugo
```

Now just run `md-hugo` like this:

``` bash
./md-hugo /path/to/process/folder /path/to/hugo/contents/folder
```

If you want to automatically `push` the hugo website live you can pass in the `-push` flag

``` bash
./md-hugo -push /path/to/process/folder /path/to/hugo/contents/folder
```
