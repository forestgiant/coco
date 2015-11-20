# Coco
Converts and moves Markdown files to Hugo content folders and adds meta data.

Just do a `git clone` on this repo
Then `cd` into `coco`

``` bash
git pull git@github.com:forestgiant/coco.git

cd coco

```
Now run `go install` and coco should be installed on your computer. To run it:

``` bash
 coco /path/to/process/folder /path/to/hugo/contents/folder
```

If you want to automatically `push` the hugo website live you can pass in the `-push` flag

``` bash
coco -push /path/to/process/folder /path/to/hugo/contents/folder
```
