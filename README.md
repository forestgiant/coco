# Coco
Converts and moves Markdown files to Hugo content folders and adds meta data.

### If you have `go` installed

Just do a git pull on this repo
Then `cd` into `coco`

``` bash
git pull git@github.com:forestgiant/coco.git

cd coco

```
Now run `go install` and coco should be installed on your computer. To run it:

``` bash
 /path/to/process/folder /path/to/hugo/contents/folder
```

If you want to automatically `push` the hugo website live you can pass in the `-push` flag

``` bash
coco -push /path/to/process/folder /path/to/hugo/contents/folder
```

### If you DO NOT have `go` installed

You can just run the binary so clone the repo and cd into it:

``` bash
git pull https://github.com/forestgiant/coco.git

cd coco
```

Now just run `coco` like this:

``` bash
./coco /path/to/process/folder /path/to/hugo/contents/folder
```

If you want to automatically `push` the hugo website live you can pass in the `-push` flag

``` bash
./coco -push /path/to/process/folder /path/to/hugo/contents/folder
```
